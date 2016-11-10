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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/openwhisk/wskdeploy/utils"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
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

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "wskdeploy",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var manifestPath = path.Join(projectPath, "serverless.yml")
		fmt.Println("Searching for manifest on path ", manifestPath)
		if _, err := os.Stat(manifestPath); err == nil {
			fmt.Println("Found severless manifest")

			dat, err := ioutil.ReadFile(manifestPath)
			utils.Check(err)
			//fmt.Println(string(dat))

			var manifest utils.Manifest

			err = yaml.Unmarshal(dat, &manifest)
			utils.Check(err)

			if manifest.Provider.Name != "openwhisk" {
				execErr := executeServerless()
				utils.Check(execErr)
			} else {
				execErr := executeDeployer(projectPath)
				utils.Check(execErr)
			}
		} else {
			fmt.Println("No manfiest files found.")
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
	RootCmd.Flags().StringVarP(&projectPath, "path", "p", ".", "path to serverless project")
	RootCmd.Flags().StringVarP(&manifestPath, "manifest", "m", ".", "path to manifest file")
	RootCmd.Flags().StringVarP(&deploymentPath, "deployment", "d", ".", "path to deployment file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".wskdeploy") // name of config file (without extension)
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func executeDeployer(manifestPath string) error {
	userHome := utils.GetHomeDirectory()
	propPath := path.Join(userHome, ".wskprops")
	deployer.LoadConfiguration(propPath)
	deployer.ReadDirectory(manifestPath)
	deployer.DeployActions()

	return nil
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
