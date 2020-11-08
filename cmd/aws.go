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
)

// ListUnattachedClassicLBs lists unattached ELBs
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
