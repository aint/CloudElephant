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
package cmd

import (
	"fmt"
	"time"

	"github.com/aint/CloudElephant/cmd/aws"

	"github.com/spf13/cobra"
)

// idleCmd represents the idle command
var idleCmd = &cobra.Command{
	Use:       "idle",
	Short:     "Find idle cloud resources",
	Long:      `Scan your EBS and find idle ones.`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"ebs"},
	Run: func(cmd *cobra.Command, args []string) {
		ticker := time.NewTicker(200 * time.Millisecond)
		tickerDone := make(chan bool)

		go printProgressBar(ticker, tickerDone)

		resultList, err := findIdleResources(args[0])
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}

		tickerDone <- true

		for _, result := range resultList {
			fmt.Println("\n", result.Label)
			for _, res := range result.Resources {
				fmt.Println(" - ", res)
			}
		}
	},
}

func findIdleResources(resourceType string) ([]aws.Result, error) {
	switch resourceType {
	case "ebs":
		return aws.ListIdleEBSs()
	default:
		return nil, fmt.Errorf("Unknown resource type '%s", resourceType)
	}
}

func init() {
	rootCmd.AddCommand(idleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// idleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// idleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
