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

// ListAvailableEBSs lists EBS volumes with available status
func ListAvailableEBSs() ([]string, error) {
	sess, err := newSession()
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	ebsAvailableStatus := "available"
	ebsStatusFilterName := "status"
	filter := &ec2.Filter{
		Name:   &ebsStatusFilterName,
		Values: []*string{&ebsAvailableStatus},
	}
	volumeInput := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{filter},
	}

	volumeOutput, err := ec2Svc.DescribeVolumes(volumeInput)
	if err != nil {
		return nil, fmt.Errorf("Error describing EBSs: %w", err)
	}

	ebsList := make([]string, 0)
	for _, volume := range volumeOutput.Volumes {
		var ebsName *string
		for _, tag := range volume.Tags {
			if *tag.Key == "Name" {
				ebsName = tag.Value
			}
		}
		ebsEntry := *volume.VolumeId
		if ebsName != nil {
			ebsEntry = ebsEntry + ", " + *ebsName
		}
		ebsList = append(ebsList, ebsEntry)
	}

	return ebsList, nil
}

// ListEBSsOnStoppedEC2 lists EBS volumes attached to stopped EC2 instances
func ListEBSsOnStoppedEC2() ([]string, error) {
	sess, err := newSession()
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	instanceStateName := "instance-state-name"
	stoppedStatus := "stopped"
	filter := &ec2.Filter{
		Name:   &instanceStateName,
		Values: []*string{&stoppedStatus},
	}

	instancesInput := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{filter},
	}
	instancesOutput, err := ec2Svc.DescribeInstances(instancesInput)
	if err != nil {
		return nil, fmt.Errorf("Error describing EC2 instances: %w", err)
	}

	ebsList := make([]string, 0)
	for _, reservation := range instancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			for _, blockDev := range instance.BlockDeviceMappings {
				volumeInput := &ec2.DescribeVolumesInput{
					VolumeIds: []*string{blockDev.Ebs.VolumeId},
				}
				volumeOutput, err := ec2Svc.DescribeVolumes(volumeInput)
				if err != nil {
					return nil, fmt.Errorf("Error describing EBS volumes: %w", err)
				}

				var name *string
				for _, tag := range volumeOutput.Volumes[0].Tags {
					if *tag.Key == "Name" {
						name = tag.Value
					}
				}

				ebsEntry := *volumeOutput.Volumes[0].VolumeId + ", type: " + *volumeOutput.Volumes[0].VolumeType
				if name != nil {
					ebsEntry = *name + ", " + ebsEntry
				}
				ebsList = append(ebsList, ebsEntry)
			}
		}
	}

	return ebsList, nil
}
