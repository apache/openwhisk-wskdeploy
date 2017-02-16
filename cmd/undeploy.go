// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"errors"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// undeployCmd represents the undeploy command
var undeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Remove all entities currently deployed in OpenWhisk.",
	Long:  `Remove all entities currently deployed in OpenWhisk according to manifest yml.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("undeploy called")
		log.Println("Starting OpenWhisk deployment")
		_, err := os.Stat(manifestPath)
		if err != nil {
			err = errors.New("manifest file not found.")
		}
		utils.Check(err)
		deployer, err := executeDeployer(manifestPath, false)
		utils.Check(err)
		if deployer.InteractiveChoice {
			log.Println("Uneployment complete")
		}
	},
}

func init() {
	RootCmd.AddCommand(undeployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// undeployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// undeployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	undeployCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "", "path to manifest file")
}
