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
	"log"
	"os"
	"path"

	"github.com/openwhisk/openwhisk-wskdeploy/deployers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/openwhisk/openwhisk-client-go/whisk"
)

const ManifestFileName = "manifest"
const DeploymentFileName = "deployment"

var cfgFile string
var CliVersion string
var CliBuild string

var Verbose bool
var projectPath string
var deploymentPath string
var manifestPath string
var useDefaults string
var useInteractive string

var RootCmd = &cobra.Command{
	Use:   "wskdeploy",
	Short: "A tool set to help deploy your openwhisk packages in batch.",
	Long: `wskdeploy is a tool to help deploy your packages, feeds, actions, triggers,
rules onto OpenWhisk platform in batch. The deployment is based on the manifest
and deployment yaml file.

A sample manifest yaml file is as below:

package:
  name: triggerrule
  version: 1.0
  license: Apache-2.0
  actions:
    hello:
      version: 1.0
      location: src/greeting.js
      runtime: nodejs
      inputs:
        name: string
        place: string
      outputs:
        payload: string
  triggers:
    locationUpdate:
  rules:
    myRule:
      trigger: locationUpdate
      action: hello
===========================================
A sample deployment yaml file is as below:

application:
  name: wskdeploy-samples
  namespace: guest

  packages:
    triggerrule:
      credential: 23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP

      actions:
        hello:
          inputs:
            name: Bernie
            place: Vermont
      `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

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

			err := deployer.ConstructDeploymentPlan()
			utils.Check(err)

			err = deployer.Deploy()
			utils.Check(err)

		} else {
			log.Println("missing manifest.yaml file")
		}

	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wskdeploy.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.Flags().StringVarP(&projectPath, "pathpath", "p", ".", "path to serverless project")
	RootCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "", "path to manifest file")
	RootCmd.Flags().StringVarP(&deploymentPath, "deployment", "d", "", "path to deployment file")
	RootCmd.Flags().StringVar(&useDefaults, "allow-defaults", "false", "allow defaults")
	RootCmd.Flags().StringVar(&useInteractive, "allow-interactive", "false", "allow interactive prompts")
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".wskdeploy") // name of config file (without extension)
	viper.AddConfigPath("$HOME")      // adding home directory as first search path
	viper.AutomaticEnv()              // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
