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

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

// ListUnattachedELBs returns unattached Application and Network Load Balancers
func ListUnattachedELBs() ([]string, error) {
	sess, err := newSession()
	if err != nil {
		return nil, err
	}
	svc := elbv2.New(sess)

	elbList, err := describeAllELBs(svc)
	if err != nil {
		return nil, err
	}

	unattachedELBList := make([]string, 0)
	for _, elb := range elbList {
		input := &elbv2.DescribeTargetGroupsInput{
			LoadBalancerArn: elb.LoadBalancerArn,
		}
		output, err := svc.DescribeTargetGroups(input)
		if err != nil {
			return nil, fmt.Errorf("Error describing Target Groups for %v: %w", *elb.LoadBalancerName, err)
		}
		unused, err := targetGroupsNotInUse(svc, output.TargetGroups)
		if err != nil {
			return nil, err
		}
		if unused {
			region := strings.Split(*elb.LoadBalancerArn, ":")[3]
			elbWithRegion := fmt.Sprint(*elb.LoadBalancerName, ", region: ", region)
			unattachedELBList = append(unattachedELBList, elbWithRegion)
		}
	}
	return unattachedELBList, nil
}

func targetGroupsNotInUse(elbSvc *elbv2.ELBV2, targetGroups []*elbv2.TargetGroup) (bool, error) {
	for _, targetGroup := range targetGroups {
		input := &elbv2.DescribeTargetHealthInput{
			TargetGroupArn: targetGroup.TargetGroupArn,
		}
		output, err := elbSvc.DescribeTargetHealth(input)
		if err != nil {
			return false, fmt.Errorf("Error describing Target Groups Health %v: %w", *targetGroup.TargetGroupArn, err)
		}
		if len(output.TargetHealthDescriptions) > 0 {
			return false, nil
		}
	}
	return true, nil
}

func describeAllELBs(elbV2Svc *elbv2.ELBV2) ([]*elbv2.LoadBalancer, error) {
	input := &elbv2.DescribeLoadBalancersInput{}
	result, err := elbV2Svc.DescribeLoadBalancers(input)
	if err != nil {
		return nil, fmt.Errorf("Error describing ELBs: %w", err)
	}

	return result.LoadBalancers, nil
}

// ListUnattachedClassicLBs returns unattached Classic Load Balancers
func ListUnattachedClassicLBs() ([]string, error) {
	sess, err := newSession()
	if err != nil {
		return nil, err
	}
	elbSvc := elb.New(sess)

	elbList, err := describeAllClassicLBs(elbSvc)
	if err != nil {
		return nil, err
	}

	unattachedELBList := make([]string, 0)
	for _, elb := range elbList {
		if len(elb.Instances) == 0 {
			region := extractClassicLBRegion(elb)
			elbWithRegion := fmt.Sprint(*elb.LoadBalancerName, ", region: ", region)
			unattachedELBList = append(unattachedELBList, elbWithRegion)
		}
	}

	return unattachedELBList, nil
}

func describeAllClassicLBs(elbSvc *elb.ELB) ([]*elb.LoadBalancerDescription, error) {
	input := &elb.DescribeLoadBalancersInput{}
	output, err := elbSvc.DescribeLoadBalancers(input)
	if err != nil {
		return nil, fmt.Errorf("Error describing classic ELBs: %w", err)
	}

	return output.LoadBalancerDescriptions, nil
}

func extractClassicLBRegion(elb *elb.LoadBalancerDescription) string {
	var region string
	if elb.AvailabilityZones != nil && len(elb.AvailabilityZones) > 0 {
		az := *elb.AvailabilityZones[0]
		region = az[:len(az)-1]
	}
	if elb.CanonicalHostedZoneName != nil {
		region = strings.Split(*elb.CanonicalHostedZoneName, ".")[1]
	} else if elb.DNSName != nil {
		region = strings.Split(*elb.DNSName, ".")[1]
	}
	return region
}
