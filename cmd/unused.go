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

// unusedCmd represents the unused command
var unusedCmd = &cobra.Command{
	Use:       "unused",
	Short:     "Find unused cloud resources",
	Long:      `Scan your ELBs, EBSs, EIPs, AMIs and find unused ones.`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"elb", "ebs", "eip", "ami"},
	Run: func(cmd *cobra.Command, args []string) {
		ticker := time.NewTicker(200 * time.Millisecond)
		tickerDone := make(chan bool)

		go printProgressBar(ticker, tickerDone)

		resultList := make([]string, 0)
		switch arg1 := args[0]; arg1 {
		case "eip":
			eipList, err := aws.ListUnattachedElasticIPs()
			if err != nil {
				fmt.Println("Error", err)
			}
			resultList = eipList
			break
		case "elb":
			clbList, err1 := aws.ListUnattachedClassicLBs()
			if err1 != nil {
				fmt.Println("Error", err1)
			}
			elbList, err2 := aws.ListUnattachedELBs()
			if err2 != nil {
				fmt.Println("Error", err2)
			}
			resultList = append(clbList, elbList...)
			break
		case "ebs":
			ebsList1, err1 := aws.ListAvailableEBSs()
			if err1 != nil {
				fmt.Println("Error", err1)
			}
			ebsList2, err2 := aws.ListEBSsOnStoppedEC2()
			if err2 != nil {
				fmt.Println("Error", err2)
			}
			resultList = append(ebsList1, ebsList2...)
			break
		case "ami":
			amiList, err := aws.ListUnusedAMIs()
			if err != nil {
				fmt.Println("Error", err)
			}
			resultList = amiList
			break
		}

		tickerDone <- true

		for _, el := range resultList {
			fmt.Println(" - ", el)
		}
	},
}

func printProgressBar(ticker *time.Ticker, done chan bool) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			fmt.Print(".")
		}
	}
}

func init() {
	rootCmd.AddCommand(unusedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unusedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unusedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
