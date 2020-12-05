/*
Copyright © 2020 Oleksandr Tyshkovets <olexandr.tyshkovets@gmail.com>

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
package cmd

import (
	"strings"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

// ListUnattachedELBs lists unattached Application and Network Load Balancers
func ListUnattachedELBs() {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	svc := elbv2.New(sess)

	elbList := describeAllELBs()
	for _, elb := range elbList {
		input := &elbv2.DescribeTargetGroupsInput{
			LoadBalancerArn: elb.LoadBalancerArn,
		}
		output, err1 := svc.DescribeTargetGroups(input)
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		unused, err2 := targetGroupsNotInUse(svc, output.TargetGroups)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		if unused {
			fmt.Println("ELB is unused: ", *elb.LoadBalancerName)
		}
	}
}

func targetGroupsNotInUse(elbSvc *elbv2.ELBV2, targetGroups []*elbv2.TargetGroup) (bool, error) {
	for _, targetGroup := range targetGroups {
		input := &elbv2.DescribeTargetHealthInput{
			TargetGroupArn: targetGroup.TargetGroupArn,
		}
		output, err := elbSvc.DescribeTargetHealth(input)
		if err != nil {
			return false, err
		}
		if len(output.TargetHealthDescriptions) > 0 {
			// fmt.Println("TargetGroup is in use: ", *targetGroup.TargetGroupName)
			return false, nil
		}
	}
	return true, nil
}

func describeAllELBs() []*elbv2.LoadBalancer {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	svc := elbv2.New(sess)
	input := &elbv2.DescribeLoadBalancersInput{}

	result, err := svc.DescribeLoadBalancers(input)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return result.LoadBalancers
}

// ListUnattachedClassicLBs lists unattached Classic Load Balancers
func ListUnattachedClassicLBs() {
	elbList := describeAllClassicLBs()
	for _, elb := range elbList {
		if len(elb.Instances) == 0 {
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
			fmt.Println(*elb.LoadBalancerName, ", region: ", region)

			// fmt.Println(elb)
		}
	}
}

func describeAllClassicLBs() []*elb.LoadBalancerDescription {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	svc := elb.New(sess)
	input := &elb.DescribeLoadBalancersInput{}

	result, err := svc.DescribeLoadBalancers(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elb.ErrCodeAccessPointNotFoundException:
				fmt.Println(elb.ErrCodeAccessPointNotFoundException, aerr.Error())
			case elb.ErrCodeDependencyThrottleException:
				fmt.Println(elb.ErrCodeDependencyThrottleException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result.LoadBalancerDescriptions

}
