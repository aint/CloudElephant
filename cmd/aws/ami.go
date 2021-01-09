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
	"github.com/aws/aws-sdk-go/service/ec2"
)

// ListUnusedAMIs lists unused AMIs
func ListUnusedAMIs() ([]string, error) {
	sess, err := newSession()
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	self := "self"
	imageInput := &ec2.DescribeImagesInput{
		Owners: []*string{&self},
	}
	imagesOutput, err := ec2Svc.DescribeImages(imageInput)
	if err != nil {
		return nil, fmt.Errorf("Error describing images: %w", err)
	}

	amiList := make([]string, 0)
	for _, img := range imagesOutput.Images {
		imageID := img.ImageId
		imageIDFilterName := "image-id"
		filter := &ec2.Filter{
			Name:   &imageIDFilterName,
			Values: []*string{imageID},
		}
		ec2Input := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{filter},
		}
		ec2Output, err := ec2Svc.DescribeInstances(ec2Input)
		if err != nil {
			return nil, fmt.Errorf("Error describing ec2 instances: %w", err)
		}
		if len(ec2Output.Reservations) == 0 {
			amiEntry := fmt.Sprint(*img.Name, ", imageId: ", *imageID)
			amiList = append(amiList, amiEntry)
		}
	}

	return amiList, nil
}