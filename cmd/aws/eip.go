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
package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// ListUnattachedElasticIPs returns unattached elastic IP addresses
func ListUnattachedElasticIPs() ([]string, error) {
	sess, err := newSession()
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	describeInput := &ec2.DescribeAddressesInput{}
	output, err := ec2Svc.DescribeAddresses(describeInput)
	if err != nil {
		return nil, err
	}

	unattachedEIPList := make([]string, 0)
	for _, address := range output.Addresses {
		if address.AssociationId == nil {
			var name *string
			for _, tag := range address.Tags {
				if *tag.Key == "Name" {
					name = tag.Value
				}
			}

			if name != nil {
				eipWithRegion := fmt.Sprint(*name, ", region: ", *address.NetworkBorderGroup)
				unattachedEIPList = append(unattachedEIPList, eipWithRegion)
			} else {
				eipWithRegion := fmt.Sprint(*address.PublicIp, ", region: ", *address.NetworkBorderGroup)
				unattachedEIPList = append(unattachedEIPList, eipWithRegion)
			}
		}
	}

	return unattachedEIPList, nil
}
