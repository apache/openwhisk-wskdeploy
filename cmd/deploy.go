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
	"github.com/openwhisk/wsktool/utils"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
)

var deployCmdPath string

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		var manifestPath = path.Join(deployCmdPath, "serverless.yml")
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
				//execErr := executeOpenWhisk(manifest, deployCmdPath)
				execErr := executeDeployer(deployCmdPath)
				utils.Check(execErr)
			}
		} else {
			fmt.Println("No manfiest files found.")
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
	deployCmd.Flags().StringVarP(&deployCmdPath, "path", "p", ".", "path to serverless project")
	deployCmd.Flags().StringVarP(&deployCmdPath, "manifest", "m", ".", "path to manifest file")
	deployCmd.Flags().StringVarP(&deployCmdPath, "deployment", "d", ".", "path to deployment file")
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

	os.Chdir(deployCmdPath)
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		return &utils.ServerlessErr{execErr.Error()}
	}

	return nil
}
