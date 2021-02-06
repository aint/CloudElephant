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
	"github.com/aws/aws-sdk-go/service/ec2"
)

func describeEC2Instances(ids []*string, filters []*ec2.Filter, ec2Svc *ec2.EC2) ([]*ec2.Instance, error) {
	instancesInput := &ec2.DescribeInstancesInput{
		Filters: filters,
		InstanceIds: ids,
	}

	instances := make([]*ec2.Instance, 0)
	err := ec2Svc.DescribeInstancesPages(instancesInput, func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
		for _, reservation := range page.Reservations {
			instances = append(instances, reservation.Instances...)
		}
		return lastPage
	})

	return instances, err
}
