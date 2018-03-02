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
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/deployers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
	"github.com/spf13/cobra"
)

var stderr = ""
var stdout = ""

// Whisk Deploy has root command: wskdeploy
// wskdeploy is being created using Cobra Library
var RootCmd = &cobra.Command{
	Use:           "wskdeploy",
	SilenceErrors: true,
	SilenceUsage:  true,
	Short:         wski18n.T(wski18n.ID_CMD_DESC_SHORT_ROOT),
	Long:          wski18n.T(wski18n.ID_CMD_DESC_LONG_ROOT),
	RunE:          RootCmdImp,
}

func RootCmdImp(cmd *cobra.Command, args []string) error {
	return Deploy()
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if utils.Flags.WithinOpenWhisk {
		err := substCmdArgs()
		if err != nil {
			wskprint.PrintOpenWhiskFromError(err)
			return
		}
	}

	if err := RootCmd.Execute(); err != nil {
		wskprint.PrintOpenWhiskFromError(err)
		os.Exit(-1)
	} else {
		if utils.Flags.WithinOpenWhisk {
			// TODO() Why are we printing success here?
			// TODO() maybe return report of what has been deployed.
			wskprint.PrintlnOpenWhiskSuccess(wski18n.T(wski18n.ID_MSG_DEPLOYMENT_SUCCEEDED))
		}
	}
}

// This function is only used when wskdeploy is being called as an Action and its input
// (i.e., command and arguments) is JSON data (map).
func substCmdArgs() error {
	// Extract arguments from input JSON string
	// { "cmd": ".." } // space-separated arguments

	arg := os.Args[1]

	// TODO() Move to proper status output/debug/trace
	fmt.Println("arg is " + arg)
	// unmarshal the string to a JSON object
	var obj map[string]interface{}
	json.Unmarshal([]byte(arg), &obj)

	if v, ok := obj["cmd"].(string); ok {
		regex, _ := regexp.Compile("[ ]+")
		os.Args = regex.Split("wskdeploy "+strings.TrimSpace(v), -1)
	} else {
		return errors.New(wski18n.T(wski18n.ID_ERR_JSON_MISSING_KEY_CMD))
	}
	return nil
}

func init() {
	// TODO() move Env var. to some global const
	utils.Flags.WithinOpenWhisk = len(os.Getenv("__OW_API_HOST")) > 0

	cobra.OnInitialize(initConfig)

	// Defining Persistent Flags of Whisk Deploy Root command (wskdeploy)
	// Persistent flags are global in terms of its availability and acceptable
	// with any other Whisk Deploy command e.g. undeploy, export, etc.
	// TODO() Publish command, not completed
	// TODO() Report command, not completed
	RootCmd.PersistentFlags().StringVar(&utils.Flags.CfgFile, "config", "", wski18n.T(wski18n.ID_CMD_FLAG_CONFIG))
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.ProjectPath, "project", "p", ".", wski18n.T(wski18n.ID_CMD_FLAG_PROJECT))
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.ManifestPath, "manifest", "m", "", wski18n.T(wski18n.ID_CMD_FLAG_MANIFEST))
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.DeploymentPath, "deployment", "d", "", wski18n.T(wski18n.ID_CMD_FLAG_DEPLOYMENT))
	RootCmd.PersistentFlags().BoolVarP(&utils.Flags.Strict, "strict", "s", false, wski18n.T(wski18n.ID_CMD_FLAG_STRICT))
	RootCmd.PersistentFlags().BoolVarP(&utils.Flags.UseInteractive, "allow-interactive", "i", false, wski18n.T(wski18n.ID_CMD_FLAG_INTERACTIVE))
	RootCmd.PersistentFlags().BoolVarP(&utils.Flags.Verbose, "verbose", "v", false, wski18n.T(wski18n.ID_CMD_FLAG_VERBOSE))
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.ApiHost, "apihost", "", "", wski18n.T(wski18n.ID_CMD_FLAG_API_HOST))
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.Namespace, "namespace", "n", "", wski18n.T(wski18n.ID_CMD_FLAG_NAMESPACE))
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.Auth, "auth", "u", "", wski18n.T(wski18n.ID_CMD_FLAG_AUTH_KEY))
	RootCmd.PersistentFlags().StringVar(&utils.Flags.ApiVersion, "apiversion", "", wski18n.T(wski18n.ID_CMD_FLAG_API_VERSION))
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.Key, "key", "k", "", wski18n.T(wski18n.ID_CMD_FLAG_KEY_FILE))
	RootCmd.PersistentFlags().StringVarP(&utils.Flags.Cert, "cert", "c", "", wski18n.T(wski18n.ID_CMD_FLAG_CERT_FILE))
	RootCmd.PersistentFlags().BoolVarP(&utils.Flags.Managed, "managed", "", false, wski18n.T(wski18n.ID_CMD_FLAG_MANAGED))
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
			warn := wski18n.T(wski18n.ID_WARN_CONFIG_INVALID_X_path_X,
				map[string]interface{}{
					wski18n.KEY_PATH: utils.Flags.CfgFile})
			wskprint.PrintOpenWhiskWarning(warn)
		}

	} else {
		utils.Flags.CfgFile = defaultPath
	}
}

func setSupportedRuntimes(apiHost string) {
	op, error := utils.ParseOpenWhisk(apiHost)
	if error == nil {
		utils.SupportedRunTimes = utils.ConvertToMap(op)
		utils.DefaultRunTimes = utils.DefaultRuntimes(op)
		utils.FileExtensionRuntimeKindMap = utils.FileExtensionRuntimes(op)
		utils.FileRuntimeExtensionsMap = utils.FileRuntimeExtensions(op)
	}
}

func Deploy() error {

	whisk.SetVerbose(utils.Flags.Verbose)
	// Verbose mode is the only mode for wskdeploy to turn on all the debug messages,
	// so set Verbose mode (and also debug mode) to true.
	whisk.SetDebug(utils.Flags.Verbose)

	project_Path := strings.TrimSpace(utils.Flags.ProjectPath)
	if len(project_Path) == 0 {
		project_Path = utils.DEFAULT_PROJECT_PATH
	}
	projectPath, _ := filepath.Abs(project_Path)

	// TODO() identical code block below; please create function both can share
	if utils.Flags.ManifestPath == "" {
		if _, err := os.Stat(path.Join(projectPath, utils.ManifestFileNameYaml)); err == nil {
			utils.Flags.ManifestPath = path.Join(projectPath, utils.ManifestFileNameYaml)
			stdout = wski18n.T(wski18n.ID_MSG_MANIFEST_DEPLOY_X_path_X,
				map[string]interface{}{wski18n.KEY_PATH: utils.Flags.ManifestPath})
		} else if _, err := os.Stat(path.Join(projectPath, utils.ManifestFileNameYml)); err == nil {
			utils.Flags.ManifestPath = path.Join(projectPath, utils.ManifestFileNameYml)
			stdout = wski18n.T(wski18n.ID_MSG_MANIFEST_DEPLOY_X_path_X,
				map[string]interface{}{wski18n.KEY_PATH: utils.Flags.ManifestPath})
		} else {
			stderr = wski18n.T(wski18n.ID_ERR_MANIFEST_FILE_NOT_FOUND_X_path_X,
				map[string]interface{}{wski18n.KEY_PATH: projectPath})
			return wskderrors.NewErrorManifestFileNotFound(projectPath, stderr)
		}
		whisk.Debug(whisk.DbgInfo, stdout)
	}

	if utils.Flags.DeploymentPath == "" {
		if _, err := os.Stat(path.Join(projectPath, utils.DeploymentFileNameYaml)); err == nil {
			utils.Flags.DeploymentPath = path.Join(projectPath, utils.DeploymentFileNameYaml)
		} else if _, err := os.Stat(path.Join(projectPath, utils.DeploymentFileNameYml)); err == nil {
			utils.Flags.DeploymentPath = path.Join(projectPath, utils.DeploymentFileNameYml)
		}
	}

	if utils.MayExists(utils.Flags.ManifestPath) {

		var deployer = deployers.NewServiceDeployer()
		deployer.ProjectPath = projectPath
		deployer.ManifestPath = utils.Flags.ManifestPath
		deployer.DeploymentPath = utils.Flags.DeploymentPath
		deployer.IsInteractive = utils.Flags.UseInteractive

		// master record of any dependency that has been downloaded
		deployer.DependencyMaster = make(map[string]utils.DependencyRecord)

		clientConfig, error := deployers.NewWhiskConfig(utils.Flags.CfgFile, utils.Flags.DeploymentPath, utils.Flags.ManifestPath, deployer.IsInteractive)
		if error != nil {
			return error
		}

		whiskClient, error := deployers.CreateNewClient(clientConfig)
		if error != nil {
			return error
		}

		deployer.Client = whiskClient
		deployer.ClientConfig = clientConfig

		// The auth, apihost and namespace have been chosen, so that we can check the supported runtimes here.
		setSupportedRuntimes(clientConfig.Host)

		err := deployer.ConstructDeploymentPlan()

		if err != nil {
			return err
		}

		err = deployer.Deploy()

		if err != nil {
			return err
		} else {
			return nil
		}

	} else {
		errString := wski18n.T(wski18n.ID_ERR_MANIFEST_FILE_NOT_FOUND_X_path_X,
			map[string]interface{}{wski18n.KEY_PATH: utils.Flags.ManifestPath})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewErrorManifestFileNotFound(utils.Flags.ManifestPath, errString)
	}

}

func Undeploy() error {

	whisk.SetVerbose(utils.Flags.Verbose)
	// Verbose mode is the only mode for wskdeploy to turn on all the debug messages, so the currenty Verbose mode
	// also set debug mode to true.
	whisk.SetDebug(utils.Flags.Verbose)

	project_Path := strings.TrimSpace(utils.Flags.ProjectPath)
	if len(project_Path) == 0 {
		project_Path = utils.DEFAULT_PROJECT_PATH
	}
	projectPath, _ := filepath.Abs(project_Path)

	if utils.Flags.ManifestPath == "" {
		if _, err := os.Stat(path.Join(projectPath, utils.ManifestFileNameYaml)); err == nil {
			utils.Flags.ManifestPath = path.Join(projectPath, utils.ManifestFileNameYaml)
			stdout = wski18n.T(wski18n.ID_MSG_MANIFEST_UNDEPLOY_X_path_X,
				map[string]interface{}{wski18n.KEY_PATH: utils.Flags.ManifestPath})
		} else if _, err := os.Stat(path.Join(projectPath, utils.ManifestFileNameYml)); err == nil {
			utils.Flags.ManifestPath = path.Join(projectPath, utils.ManifestFileNameYml)
			stdout = wski18n.T(wski18n.ID_MSG_MANIFEST_UNDEPLOY_X_path_X,
				map[string]interface{}{wski18n.KEY_PATH: utils.Flags.ManifestPath})
		} else {
			stderr = wski18n.T(wski18n.ID_ERR_MANIFEST_FILE_NOT_FOUND_X_path_X,
				map[string]interface{}{wski18n.KEY_PATH: projectPath})
			return wskderrors.NewErrorManifestFileNotFound(projectPath, stderr)
		}
		wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, stdout)
	}

	if utils.Flags.DeploymentPath == "" {
		if _, err := os.Stat(path.Join(projectPath, utils.DeploymentFileNameYaml)); err == nil {
			utils.Flags.DeploymentPath = path.Join(projectPath, utils.DeploymentFileNameYaml)
			// TODO() have a single function that conditionally (verbose) prints ALL Flags
			dbgMsg := fmt.Sprintf("%s >> [%s]: [%s]",
				wski18n.T(wski18n.ID_DEBUG_UNDEPLOYING_USING),
				wski18n.DEPLOYMENT,
				utils.Flags.DeploymentPath)
			wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, dbgMsg)

		} else if _, err := os.Stat(path.Join(projectPath, utils.DeploymentFileNameYml)); err == nil {
			utils.Flags.DeploymentPath = path.Join(projectPath, utils.DeploymentFileNameYml)
			// TODO() have a single function that conditionally (verbose) prints ALL Flags
			dbgMsg := fmt.Sprintf("%s >> [%s]: [%s]",
				wski18n.T(wski18n.ID_DEBUG_UNDEPLOYING_USING),
				wski18n.DEPLOYMENT,
				utils.Flags.DeploymentPath)
			wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, dbgMsg)
		}
	}

	if utils.FileExists(utils.Flags.ManifestPath) {

		var deployer = deployers.NewServiceDeployer()
		deployer.ProjectPath = utils.Flags.ProjectPath
		deployer.ManifestPath = utils.Flags.ManifestPath
		deployer.DeploymentPath = utils.Flags.DeploymentPath
		deployer.IsInteractive = utils.Flags.UseInteractive

		clientConfig, error := deployers.NewWhiskConfig(utils.Flags.CfgFile, utils.Flags.DeploymentPath, utils.Flags.ManifestPath, deployer.IsInteractive)
		if error != nil {
			return error
		}

		whiskClient, error := deployers.CreateNewClient(clientConfig)
		if error != nil {
			return error
		}

		deployer.Client = whiskClient
		deployer.ClientConfig = clientConfig

		// The auth, apihost and namespace have been chosen, so that we can check the supported runtimes here.
		setSupportedRuntimes(clientConfig.Host)

		verifiedPlan, err := deployer.ConstructUnDeploymentPlan()
		if err != nil {
			return err
		}

		err = deployer.UnDeploy(verifiedPlan)
		if err != nil {
			return err
		} else {
			return nil
		}

	} else {
		errString := wski18n.T(wski18n.ID_ERR_MANIFEST_FILE_NOT_FOUND_X_path_X,
			map[string]interface{}{wski18n.KEY_PATH: utils.Flags.ManifestPath})
		return wskderrors.NewErrorManifestFileNotFound(utils.Flags.ManifestPath, errString)
	}
}
