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
	"github.com/aws/aws-sdk-go/aws"
)

// ListIdleEBSs lists idle EBS volumes
func ListIdleEBSs() ([]Result, error) {
	sess, err := newSession()
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)
	cwSvc := cloudwatch.New(sess)

	filter := &ec2.Filter{
		Name:   aws.String("status"),
		Values: aws.StringSlice([]string{"in-use"}),
	}
	volumes, err := describeVolumes(nil, []*ec2.Filter{filter}, ec2Svc)
	if err != nil {
		return nil, fmt.Errorf("error describing EBSs: %w", err)
	}
	ebsList := make([]string, 0)
	for _, volume := range volumes {
		if !isRootVolume(volume.Attachments) {
			endTime := time.Now()
			startTime := time.Now().AddDate(0, 0, -7)
			period := int64(3600)
			dimension := cloudwatch.Dimension{
				Name:  aws.String("VolumeId"),
				Value: volume.VolumeId,
			}

			metricInput := &cloudwatch.GetMetricStatisticsInput{
				MetricName: aws.String("VolumeReadOps"),
				Namespace:  aws.String("AWS/EBS"),
				Statistics: aws.StringSlice([]string{"Sum"}),
				Dimensions: []*cloudwatch.Dimension{&dimension},
				StartTime:  &startTime,
				EndTime:    &endTime,
				Period:     &period,
			}
			metricOutput, err := cwSvc.GetMetricStatistics(metricInput)
			if err != nil {
				return nil, err
			}
			count := 0.0
			for _, datapoint := range metricOutput.Datapoints {
				count = count + *datapoint.Sum
			}
			if count <= 1 {
				metricInput.MetricName = aws.String("VolumeWriteOps")
				metricOutput, err := cwSvc.GetMetricStatistics(metricInput)
				if err != nil {
					return nil, err
				}
				count := 0.0
				for _, datapoint := range metricOutput.Datapoints {
					count = count + *datapoint.Sum
				}
				if count <= 1 {
					// fmt.Println("No Write Ops: ", *attachment.Device, ", ", *attachment.VolumeId)
					ebsList = append(ebsList, *volume.VolumeId)
				}
			}
		}
	}

	return []Result{{"Idle EBS volumes:", ebsList}}, nil
}

func isRootVolume(attachments []*ec2.VolumeAttachment) bool {
	xvda := "/dev/xvda"
	sda1 := "/dev/sda1"
	for _, attachment := range attachments {
		device := attachment.Device
		if strings.HasPrefix(*device, xvda) || strings.HasPrefix(*device, sda1) {
			return true
		}
	}
	return false
}

func ListUnusedEBSs() ([]Result, error) {
	l1, err := listAvailableEBSs()
	if err != nil {
		return nil, err
	}
	l2, err := listEBSsOnStoppedEC2()
	return append(l1, l2...), err
}


// ListAvailableEBSs lists EBS volumes with available status
func listAvailableEBSs() ([]Result, error) {
	sess, err := newSession()
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	filter := &ec2.Filter{
		Name:   aws.String("status"),
		Values: aws.StringSlice([]string{"available"}),
	}
	volumes, err := describeVolumes(nil, []*ec2.Filter{filter}, ec2Svc)
	if err != nil {
		return nil, fmt.Errorf("error describing EBSs: %w", err)
	}

	ebsList := make([]string, 0)
	for _, volume := range volumes {
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

	return []Result{{"Available EBS volumes:", ebsList}}, nil
}

// ListEBSsOnStoppedEC2 lists EBS volumes attached to stopped EC2 instances
func listEBSsOnStoppedEC2() ([]Result, error) {
	sess, err := newSession()
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	volumeIDs, err := getVolumeIDsOnStoppedEC2(ec2Svc)
	if err != nil {
		return nil, err
	}

	volumes, err := describeVolumes(volumeIDs, nil, ec2Svc)
	if err != nil {
		return nil, fmt.Errorf("error describing EBS volumes: %w", err)
	}

	ebsList := make([]string, 0)
	for _, volume := range volumes {
		var name *string
		for _, tag := range volume.Tags {
			if *tag.Key == "Name" {
				name = tag.Value
			}
		}

		ebsEntry := *volume.VolumeId + ", type: " + *volume.VolumeType
		if name != nil {
			ebsEntry = *name + ", " + ebsEntry
		}
		ebsList = append(ebsList, ebsEntry)
	}

	return []Result{{"EBS volumes on stopped EC2:", ebsList}}, nil
}

func getVolumeIDsOnStoppedEC2(ec2Svc *ec2.EC2) ([]*string, error) {
	filter := &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: aws.StringSlice([]string{"stopped"}),
	}

	instances, err := describeEC2Instances(nil, []*ec2.Filter{filter}, ec2Svc)
	if err != nil {
		return nil, fmt.Errorf("error describing EC2 instances: %w", err)
	}

	volumeIDs := make([]*string, 0)
	for _, instance := range instances {
		for _, blockDev := range instance.BlockDeviceMappings {
			volumeIDs = append(volumeIDs, blockDev.Ebs.VolumeId)
		}
	}

	return volumeIDs, nil
}

func describeVolumes(ids []*string, filters []*ec2.Filter, ec2Svc *ec2.EC2) ([]*ec2.Volume, error) {
	volumes := make([]*ec2.Volume, 0)
	if len(ids) == 0 {
		return volumes, nil
	}

	volumesInput := &ec2.DescribeVolumesInput{
		Filters:   filters,
		VolumeIds: ids,
	}

	err := ec2Svc.DescribeVolumesPages(volumesInput, func(page *ec2.DescribeVolumesOutput, lastPage bool) bool {
		volumes = append(volumes, page.Volumes...)
		return lastPage
	})

	return volumes, err
}

