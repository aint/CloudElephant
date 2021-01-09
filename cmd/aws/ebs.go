/*
Copyright © 2020 - 2021 Oleksandr Tyshkovets <olexandr.tyshkovets@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package aws

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// ListIdleEBSs lists idle EBS volumes
func ListIdleEBSs() ([]string, error) {
	sess, err := newSession()
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	ebsInUseStatus := "in-use"
	ebsStatusFilterName := "status"
	filter := &ec2.Filter{
		Name:   &ebsStatusFilterName,
		Values: []*string{&ebsInUseStatus},
	}

	volumeInput := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{filter},
	}

	volumeOutput, err := ec2Svc.DescribeVolumes(volumeInput)
	if err != nil {
		return nil, err
	}
	for _, volume := range volumeOutput.Volumes {
		for _, attachment := range volume.Attachments {
			// fmt.Println(*volume.VolumeId)
			if isNotRootVolume(*attachment.Device) {
				// fmt.Println(*attachment.Device, ", ", *attachment.VolumeId)

				cwSvc := cloudwatch.New(sess)

				volumeReadOpsMetric := "VolumeReadOps"
				namespace := "AWS/EBS"
				statistic := "Sum"
				endTime := time.Now()
				startTime := time.Now().AddDate(0, -1, 0)
				period := int64(3600)
				volumeIDDimension := "VolumeId"
				dimension := cloudwatch.Dimension{
					Name:  &volumeIDDimension,
					Value: volume.VolumeId,
				}

				metricInput := &cloudwatch.GetMetricStatisticsInput{
					MetricName: &volumeReadOpsMetric,
					Namespace:  &namespace,
					Statistics: []*string{&statistic},
					Dimensions: []*cloudwatch.Dimension{&dimension},
					StartTime:  &startTime,
					EndTime:    &endTime,
					Period:     &period,
				}
				metricOutput, err := cwSvc.GetMetricStatistics(metricInput)
				if err != nil {
					return nil, err
				}
				fmt.Println("len ", len(metricOutput.Datapoints))
				count := 0.0
				for _, datapoint := range metricOutput.Datapoints {
					count = count + *datapoint.Sum
				}
				if count <= 1 {
					fmt.Println(*attachment.Device, ", ", *attachment.VolumeId)
					fmt.Println("count ", count)
					for _, datapoint := range metricOutput.Datapoints {
						fmt.Println("d: ", *datapoint.Sum)
					}
				}

			}
		}
	}

	ebsList := make([]string, 0)
	return ebsList, nil
}

func isNotRootVolume(device string) bool {
	xvda := "/dev/xvda"
	sda1 := "/dev/sda1"
	return !strings.HasPrefix(device, xvda) && !strings.HasPrefix(device, sda1)
}

// func listAvailableVolumes() ([]string, error) {
// }
