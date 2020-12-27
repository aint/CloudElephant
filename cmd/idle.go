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
	"fmt"

	"github.com/spf13/cobra"
)

// idleCmd represents the idle command
var idleCmd = &cobra.Command{
	Use:       "idle",
	Short:     "Find idle load balancers",
	Long:      `Scan your load balancers and find idle ones.`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"elb", "ebs", "eip"},
	Run: func(cmd *cobra.Command, args []string) {
		resultList := make([]string, 0)
		switch arg1 := args[0]; arg1 {
		case "eip":
			eipList, err := ListUnattachedElasticIPs()
			if err != nil {
				fmt.Println("Error", err)
			}
			resultList = eipList
			break
		case "elb":
			clbList, err1 := ListUnattachedClassicLBs()
			if err1 != nil {
				fmt.Println("Error", err1)
			}
			elbList, err2 := ListUnattachedELBs()
			if err2 != nil {
				fmt.Println("Error", err2)
			}
			resultList = append(clbList, elbList...)
		}

		for _, el := range resultList {
			fmt.Println(" - ", el)
		}
	},
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
