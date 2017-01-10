/*
 * Copyright 2015-2016 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"errors"
	"fmt"
	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/wskdeploy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
)

var cfgFile string
var projectPath string
var manifestPath string
var deploymentPath string
var useInteractive string
var useDefaults string
var Verbose bool
var clientConfig *whisk.Config
var CliVersion string
var CliBuild string

// RootCmd represents the base command when called without any subcommands
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

		log.Println("OpenWhisk Deploy initial configuration")
		log.Println("Project path is ", projectPath)

		var searchPath = path.Join(projectPath, "serverless.yaml")
		log.Println("Searching for serverless manifest on path ", searchPath)

		if _, err := os.Stat(searchPath); err == nil {
			log.Println("Found severless manifest")

			dat, err := ioutil.ReadFile(searchPath)
			utils.Check(err)

			var manifest utils.Manifest

			err = yaml.Unmarshal(dat, &manifest)
			utils.Check(err)

			if manifest.Provider.Name != "openwhisk" {
				log.Println("Starting Serverless deployment")
				execErr := executeServerless()
				utils.Check(execErr)
				fmt.Println("Deployment complete")
			} else {
				log.Println("Starting OpenWhisk deployment")
				deployer, err := executeDeployer(manifestPath)
				utils.Check(err)
				if deployer.InteractiveChoice {
					fmt.Println("Deployment complete")
				}
			}

		} else {
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
	RootCmd.Flags().StringVar(&useInteractive, "allow-interactive", "true", "allow interactive prompts")
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

func executeDeployer(manifestPath string) (*ServiceDeployer, error) {
	userHome := utils.GetHomeDirectory()
	propPath := path.Join(userHome, ".wskprops")
	deployer := NewServiceDeployer()

	isInteractive, err := utils.GetBoolFromString(useInteractive)
	utils.Check(err)

	isDefault, err := utils.GetBoolFromString(useDefaults)
	utils.Check(err)

	deployer.IsInteractive = isInteractive
	deployer.IsDefault = isDefault
	deployer.ManifestPath = manifestPath
	deployer.ProjectPath = projectPath
	deployer.DeploymentPath = deploymentPath
	// from propPath and deploymentPath get all the information
	// and we return the information back for later usage if necessary
	deployer.Client, clientConfig = utils.NewClient(propPath, deploymentPath)

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

// Process manifest using OpenWhisk Tool
func executeOpenWhisk(manifest utils.Manifest, manifestPath string) error {
	err := filepath.Walk(manifestPath, processPath)
	utils.Check(err)
	fmt.Println("OpenWhisk processing TBD")
	return nil
}

func processPath(path string, f os.FileInfo, err error) error {
	fmt.Println("Visited ", path)
	return nil
}

// If "serverless" is installed, then use it to process manifest
func executeServerless() error {
	//utils.Check
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		return &utils.ServerlessErr{"Please set missing environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY tokens"}
	}
	binary, lookErr := exec.LookPath(utils.ServerlessBinaryCommand)
	if lookErr != nil {
		panic(lookErr)
	}
	args := make([]string, 2)
	args[0] = utils.ServerlessBinaryCommand
	args[1] = "deploy"

	env := os.Environ()

	os.Chdir(projectPath)
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		return &utils.ServerlessErr{execErr.Error()}
	}

	return nil
}
