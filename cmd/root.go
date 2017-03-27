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
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/openwhisk/openwhisk-wskdeploy/deployers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"regexp"

	"encoding/json"
	"errors"
	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-client-go/wski18n"
	"strings"
)

var RootCmd = &cobra.Command{
	Use:   "wskdeploy",
	Short: "A tool set to help deploy your openwhisk packages in batch.",
	Long: `A tool to deploy openwhisk packages with a manifest and/or deployment yaml file.

wskdeploy without any commands or flags deploys openwhisk package in the current directory if manifest.yaml exists.

      `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		whisk.SetVerbose(Verbose)

		projectPath, err := filepath.Abs(projectPath)
		utils.Check(err)

		if manifestPath == "" {
			if ok, _ := regexp.Match(ManifestFileNameYml, []byte(manifestPath)); ok {
				manifestPath = path.Join(projectPath, ManifestFileNameYml)
			} else {
				manifestPath = path.Join(projectPath, ManifestFileNameYaml)
			}

		}

		if deploymentPath == "" {
			if ok, _ := regexp.Match(DeploymentFileNameYml, []byte(manifestPath)); ok {
				deploymentPath = path.Join(projectPath, DeploymentFileNameYml)
			} else {
				deploymentPath = path.Join(projectPath, DeploymentFileNameYaml)
			}

		}

		if utils.MayExists(manifestPath) {

			var deployer = deployers.NewServiceDeployer()
			deployer.ProjectPath = projectPath
			deployer.ManifestPath = manifestPath
			deployer.DeploymentPath = deploymentPath
			// perform some quick check here.
			go func() {
				deployer.Check()
			}()
			deployer.IsDefault = useDefaults

			deployer.IsInteractive = useInteractive

			propPath := ""
			if !utils.Flags.WithinOpenWhisk {
				userHome := utils.GetHomeDirectory()
				propPath = path.Join(userHome, ".wskprops")
			}
			whiskClient, clientConfig := deployers.NewWhiskClient(propPath, deploymentPath, deployer.IsInteractive)
			deployer.Client = whiskClient
			deployer.ClientConfig = clientConfig

			err := deployer.ConstructDeploymentPlan()
			utils.Check(err)

			err = deployer.Deploy()
			utils.Check(err)

		} else {
			if utils.Flags.WithinOpenWhisk {
				utils.PrintOpenWhiskError(wski18n.T("missing manifest.yaml file"))
			} else {
				log.Println("missing manifest.yaml file")
			}
		}

	},
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
	RootCmd.PersistentFlags().BoolVarP(&useInteractive, "allow-interactive", "i", !utils.Flags.WithinOpenWhisk, "allow interactive prompts")
	RootCmd.PersistentFlags().BoolVarP(&useDefaults, "allow-defaults", "a", false, "allow defaults")
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.ApiHost, "apihost", "", "", wski18n.T("whisk API HOST"))
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.Auth, "auth", "u", "", wski18n.T("authorization `KEY`"))
	RootCmd.PersistentFlags().StringVar(&utils.Flags.ApiVersion, "apiversion", "", wski18n.T("whisk API `VERSION`"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// enable ability to specify config file via flag
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
