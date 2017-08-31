/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-client-go/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/deployers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	"regexp"
	"strings"
	"path"
	"path/filepath"
)

var RootCmd = &cobra.Command{
	Use:   "wskdeploy",
	Short: "A tool set to help deploy your openwhisk packages in batch.",
	Long: `A tool to deploy openwhisk packages with a manifest and/or deployment yaml file.

wskdeploy without any commands or flags deploys openwhisk package in the current directory if manifest.yaml exists.

      `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: RootCmdImp,
}

func RootCmdImp(cmd *cobra.Command, args []string) {
	Deploy()
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if utils.Flags.WithinOpenWhisk {
		err := substCmdArgs()
		if err != nil {
			utils.PrintOpenWhiskError(err.Error())
			return
		}
	}

	if err := RootCmd.Execute(); err != nil {
		log.Println(err)
		if utils.Flags.WithinOpenWhisk {
			utils.PrintOpenWhiskError(err.Error())
		} else {
			os.Exit(-1)
		}
	} else {
		if utils.Flags.WithinOpenWhisk {
			fmt.Print(`{"deploy":"success"}`) // maybe return report of what has been deployed.
		}
	}
}

func substCmdArgs() error {
	// Extract arguments from input JSON string

	// { "cmd": ".." } // space-separated arguments

	arg := os.Args[1]

    fmt.Println("arg is " + arg)
	// unmarshal the string to a JSON object
	var obj map[string]interface{}
	json.Unmarshal([]byte(arg), &obj)

	if v, ok := obj["cmd"].(string); ok {
		regex, _ := regexp.Compile("[ ]+")
		os.Args = regex.Split("wskdeploy "+strings.TrimSpace(v), -1)
	} else {
		return errors.New(wski18n.T("Missing cmd key"))
	}
	return nil
}

func init() {
	utils.Flags.WithinOpenWhisk = len(os.Getenv("__OW_API_HOST")) > 0

	cobra.OnInitialize(initConfig)

	//Initialize the supported runtime infos.
	op, err := utils.ParseOpenWhisk()
	if err == nil {
		utils.Rts = utils.ConvertToMap(op)
	} else {
		utils.Rts = utils.DefaultRts
	}
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&utils.Flags.CfgFile, "config", "", "config file (default is $HOME/.wskprops)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.Flags().StringVarP(&utils.Flags.ProjectPath, "pathpath", "p", ".", "path to serverless project")
	RootCmd.Flags().StringVarP(&utils.Flags.ManifestPath, "manifest", "m", "", "path to manifest file")
	RootCmd.Flags().StringVarP(&utils.Flags.DeploymentPath, "deployment", "d", "", "path to deployment file")
	RootCmd.PersistentFlags().BoolVarP(&utils.Flags.UseInteractive, "allow-interactive", "i", false, "allow interactive prompts")
	RootCmd.PersistentFlags().BoolVarP(&utils.Flags.UseDefaults, "allow-defaults", "a", false, "allow defaults")
	RootCmd.PersistentFlags().BoolVarP(&utils.Flags.Verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.ApiHost, "apihost", "", "", wski18n.T("whisk API HOST"))
    RootCmd.PersistentFlags().StringVarP(&utils.Flags.Namespace, "namespace", "n", "", wski18n.T("namespace"))
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.Auth, "auth", "u", "", wski18n.T("authorization `KEY`"))
	RootCmd.PersistentFlags().StringVar(&utils.Flags.ApiVersion, "apiversion", "", wski18n.T("whisk API `VERSION`"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    userHome := utils.GetHomeDirectory()
    defaultPath := path.Join(userHome, whisk.DEFAULT_LOCAL_CONFIG)
	if utils.Flags.CfgFile != "" {
        // Read the file as a wskprops file, to check if it is valid.
        _, err := whisk.ReadProps(utils.Flags.CfgFile)
        if err != nil {
            utils.Flags.CfgFile = defaultPath
            log.Println("Invalid config file detected, so by bdefault it is set to " + utils.Flags.CfgFile)
        }
	} else {
        utils.Flags.CfgFile = defaultPath
    }
}

func Deploy() error {

	whisk.SetVerbose(utils.Flags.Verbose)

	projectPath, err := filepath.Abs(utils.Flags.ProjectPath)
	utils.Check(err)

	if utils.Flags.ManifestPath == "" {
		if _, err := os.Stat(path.Join(projectPath, utils.ManifestFileNameYaml)); err == nil {
			utils.Flags.ManifestPath = path.Join(projectPath, utils.ManifestFileNameYaml)
			log.Printf("Using %s for deployment \n", utils.Flags.ManifestPath)
		} else if _, err := os.Stat(path.Join(projectPath, utils.ManifestFileNameYaml)); err == nil {
			utils.Flags.ManifestPath = path.Join(projectPath, utils.ManifestFileNameYml)
			log.Printf("Using %s for deployment", utils.Flags.ManifestPath)
		} else {
			log.Printf("Manifest file not found at path %s", projectPath)
			return errors.New("missing manifest.yaml file")
		}
	}

	if utils.Flags.DeploymentPath == "" {
		if _, err := os.Stat(path.Join(projectPath, "deployment.yaml")); err == nil {
			utils.Flags.DeploymentPath = path.Join(projectPath, utils.DeploymentFileNameYaml)
		} else if _, err := os.Stat(path.Join(projectPath, "deployment.yml")); err == nil {
			utils.Flags.DeploymentPath = path.Join(projectPath, utils.DeploymentFileNameYml)
		}
	}

	if utils.MayExists(utils.Flags.ManifestPath) {

		var deployer = deployers.NewServiceDeployer()
		deployer.ProjectPath = projectPath
		deployer.ManifestPath = utils.Flags.ManifestPath
		deployer.DeploymentPath = utils.Flags.DeploymentPath
		// perform some quick check here.
		go func() {
			deployer.Check()
		}()
		deployer.IsDefault = utils.Flags.UseDefaults

		deployer.IsInteractive = utils.Flags.UseInteractive

		// master record of any dependency that has been downloaded
		deployer.DependencyMaster = make(map[string]utils.DependencyRecord)

		whiskClient, clientConfig := deployers.NewWhiskClient(utils.Flags.CfgFile, utils.Flags.DeploymentPath, utils.Flags.ManifestPath, deployer.IsInteractive)
		deployer.Client = whiskClient
		deployer.ClientConfig = clientConfig

		err := deployer.ConstructDeploymentPlan()

		if err != nil {
			utils.Check(err)
			return err
		}

		err = deployer.Deploy()
		if err != nil {
			utils.Check(err)
			return err
		} else {
			return nil
		}

	} else {
		if utils.Flags.WithinOpenWhisk {
			utils.PrintOpenWhiskError(wski18n.T("missing manifest.yaml file"))
			return errors.New("missing manifest.yaml file")
		} else {
			log.Println("missing manifest.yaml file")
			return errors.New("missing manifest.yaml file")
		}
	}

}

func Undeploy() error {
	// TODO: Work your own magic here
	whisk.SetVerbose(utils.Flags.Verbose)

	if utils.Flags.ManifestPath == "" {
		if ok, _ := regexp.Match(utils.ManifestFileNameYml, []byte(utils.Flags.ManifestPath)); ok {
			utils.Flags.ManifestPath = path.Join(utils.Flags.ProjectPath, utils.ManifestFileNameYml)
		} else {
			utils.Flags.ManifestPath = path.Join(utils.Flags.ProjectPath, utils.ManifestFileNameYaml)
		}

	}

	if utils.Flags.DeploymentPath == "" {
		if ok, _ := regexp.Match(utils.DeploymentFileNameYml, []byte(utils.Flags.ManifestPath)); ok {
			utils.Flags.DeploymentPath = path.Join(utils.Flags.ProjectPath, utils.DeploymentFileNameYml)
		} else {
			utils.Flags.DeploymentPath = path.Join(utils.Flags.ProjectPath, utils.DeploymentFileNameYaml)
		}

	}

	if utils.FileExists(utils.Flags.ManifestPath) {

		var deployer = deployers.NewServiceDeployer()
		deployer.ProjectPath = utils.Flags.ProjectPath
		deployer.ManifestPath = utils.Flags.ManifestPath
		deployer.DeploymentPath = utils.Flags.DeploymentPath

		deployer.IsInteractive = utils.Flags.UseInteractive
		deployer.IsDefault = utils.Flags.UseDefaults

		whiskClient, clientConfig := deployers.NewWhiskClient(utils.Flags.CfgFile, utils.Flags.DeploymentPath, utils.Flags.ManifestPath, deployer.IsInteractive)
		deployer.Client = whiskClient
		deployer.ClientConfig = clientConfig

		verifiedPlan, err := deployer.ConstructUnDeploymentPlan()
		err = deployer.UnDeploy(verifiedPlan)
		if err != nil {
			utils.Check(err)
			return err
		} else {
			return nil
		}

	} else {
		log.Println("missing manifest.yaml file")
		return errors.New("missing manifest.yaml file")
	}
}
