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
	"path"
)

var manifestPath string
var deploymentPath string
var useInteractive string
var useDefaults string

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy entities to OpenWhisk",
	Long: `Deploy entities defined according to Manifest file
to the OpenWhisk platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("deploy called")
		log.Println("Starting OpenWhisk deployment")
		_, err := os.Stat(manifestPath)
		if err != nil {
			err = errors.New("manifest file not found.")
		}
		utils.Check(err)
		deployer, err := executeDeployer(manifestPath)
		utils.Check(err)
		if deployer.InteractiveChoice {
			log.Println("Deployment complete")
		}

	},
}

func init() {
	RootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deployCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "", "path to manifest file")
	deployCmd.Flags().StringVarP(&deploymentPath, "deployment", "d", "", "path to deployment file")
	deployCmd.Flags().StringVar(&useDefaults, "allow-defaults", "false", "allow defaults")
	deployCmd.Flags().StringVar(&useInteractive, "allow-interactive", "true", "allow interactive prompts")
}

func executeDeployer(manifestPath string) (*utils.ServiceDeployer, error) {
	userHome := utils.GetHomeDirectory()
	propPath := path.Join(userHome, ".wskprops")
	// from propPath and deploymentPath get all the information
	// and we return the information back for later usage if necessary
	client, utils.ClientConfig = utils.NewClient(propPath, deploymentPath)

	deployer := utils.NewServiceDeployer()
	deployer.Client = client

	isInteractive, err := utils.GetBoolFromString(useInteractive)
	utils.Check(err)

	isDefault, err := utils.GetBoolFromString(useDefaults)
	utils.Check(err)

	deployer.IsInteractive = isInteractive
	deployer.IsDefault = isDefault
	deployer.ManifestPath = manifestPath
	deployer.ProjectPath = projectPath
	deployer.DeploymentPath = deploymentPath

	err = deployer.ConstructDeploymentPlan()
	if err != nil {
		return nil, err
	}

	err = deployer.Deploy()
	if err != nil {
		return nil, err
	}

	return deployer, nil
}
