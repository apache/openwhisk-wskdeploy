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
	"log"
	"path"

	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/deployers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
	"github.com/spf13/cobra"
)

// undeployCmd represents the undeploy command
var undeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		whisk.SetVerbose(Verbose)

		if manifestPath == "" {
			manifestPath = path.Join(projectPath, ManifestFileName+".yaml")
		}

		if deploymentPath == "" {
			deploymentPath = path.Join(projectPath, DeploymentFileName+".yaml")
		}

		if utils.FileExists(manifestPath) {

			var deployer = deployers.NewServiceDeployer()
			deployer.ProjectPath = projectPath
			deployer.ManifestPath = manifestPath
			deployer.DeploymentPath = deploymentPath

			userHome := utils.GetHomeDirectory()
			propPath := path.Join(userHome, ".wskprops")

			whiskClient, clientConfig := deployers.NewWhiskClient(propPath, deploymentPath)
			deployer.Client = whiskClient
			deployer.ClientConfig = clientConfig

			verifiedPlan, err := deployer.ConstructUnDeploymentPlan()
			utils.Check(err)

			err = deployer.UnDeploy(verifiedPlan)
			utils.Check(err)

		} else {
			log.Println("missing manifest.yaml file")
		}
	},
}

func init() {
	RootCmd.AddCommand(undeployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// undeployCmd.PersistentFlags().String("foo", "", "A help for foo")
	undeployCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wskdeploy.yaml)")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// undeployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	undeployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	undeployCmd.Flags().StringVarP(&projectPath, "pathpath", "p", ".", "path to serverless project")
	undeployCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "", "path to manifest file")
	undeployCmd.Flags().StringVarP(&deploymentPath, "deployment", "d", "", "path to deployment file")
	undeployCmd.Flags().StringVar(&useDefaults, "allow-defaults", "false", "allow defaults")
	undeployCmd.Flags().StringVar(&useInteractive, "allow-interactive", "false", "allow interactive prompts")
	undeployCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

}
