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
package azure

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-11-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/aint/CloudElephant/cmd/aws"
)

func ListUnusedLBs() ([]aws.Result, error) {
	lbClient, err := createLBClient()
	if err != nil {
		return nil, err
	}

	lbResultPage, err := lbClient.ListAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error listing load balancers: %w", err)
	}

	unattachedLBList := make([]string, 0)
	for _, lb := range lbResultPage.Values() {
		fmt.Println("LB:", *lb.Name)
		if isBackendAddressPoolsEmpty(lb.BackendAddressPools) {
			lbLocationStr := fmt.Sprint(*lb.Name, ", location: ", *lb.Location)
			unattachedLBList = append(unattachedLBList, lbLocationStr)
		}
	}

	return []aws.Result{{"Unattached LBs:", unattachedLBList}}, nil
}

func isBackendAddressPoolsEmpty(pools *[]network.BackendAddressPool) bool {
	if pools != nil || len(*pools) == 0 {
		return true
	}

	for _, bap := range *pools {
		fmt.Println("Backend Address Pool:", *bap.Name)
		ips := bap.BackendIPConfigurations
		if ips != nil && len(*ips) > 0 {
			return false
		}
	}

	return true
}

func createLBClient() (*network.LoadBalancersClient, error) {
	subID, ok := os.LookupEnv("AZURE_SUBSCRIPTION_ID")
	if !ok {
		return nil, errors.New("AZURE_SUBSCRIPTION_ID env var is not set")
	}

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("error creating authorizer configured from env vars: %w", err)
	}

	lbClient := network.NewLoadBalancersClient(subID)
	lbClient.Authorizer = authorizer

	return &lbClient, nil
}
