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

	"github.com/openwhisk/openwhisk-wskdeploy/cmdImp"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"regexp"

	"encoding/json"
	"errors"
	"strings"

	"github.com/openwhisk/openwhisk-client-go/wski18n"
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

var Deploy = cmdImp.Deploy

func RootCmdImp(cmd *cobra.Command, args []string) {
	// Set all the parameters passed via the command to the struct of wskdeploy command.
	deployParams := cmdImp.DeployParams{cmdImp.Verbose, cmdImp.ProjectPath, cmdImp.ManifestPath,
		cmdImp.DeploymentPath, cmdImp.UseDefaults, cmdImp.UseInteractive}
	// Call the implementation of wskdeploy command.
	Deploy(deployParams)

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

	RootCmd.PersistentFlags().StringVar(&cmdImp.CfgFile, "config", "", "config file (default is $HOME/.wskdeploy.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.Flags().StringVarP(&cmdImp.ProjectPath, "pathpath", "p", ".", "path to serverless project")
	RootCmd.Flags().StringVarP(&cmdImp.ManifestPath, "manifest", "m", "", "path to manifest file")
	RootCmd.Flags().StringVarP(&cmdImp.DeploymentPath, "deployment", "d", "", "path to deployment file")
	RootCmd.PersistentFlags().BoolVarP(&cmdImp.UseInteractive, "allow-interactive", "i", !utils.Flags.WithinOpenWhisk, "allow interactive prompts")
	RootCmd.PersistentFlags().BoolVarP(&cmdImp.UseDefaults, "allow-defaults", "a", false, "allow defaults")
	RootCmd.PersistentFlags().BoolVarP(&cmdImp.Verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.ApiHost, "apihost", "", "", wski18n.T("whisk API HOST"))
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.Auth, "auth", "u", "", wski18n.T("authorization `KEY`"))
	RootCmd.PersistentFlags().StringVar(&utils.Flags.ApiVersion, "apiversion", "", wski18n.T("whisk API `VERSION`"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cmdImp.CfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(cmdImp.CfgFile)
	}

	viper.SetConfigName(".wskdeploy") // name of config file (without extension)
	viper.AddConfigPath("$HOME")      // adding home directory as first search path
	viper.AutomaticEnv()              // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
