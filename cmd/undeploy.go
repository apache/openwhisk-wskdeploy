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
	"github.com/spf13/cobra"
	"github.com/openwhisk/openwhisk-wskdeploy/cmdImp"
)

// undeployCmd represents the undeploy command
var undeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploy assets from OpenWhisk",
	Long:  `Undeploy removes deployed assets from the manifest and deployment files`,
	Run: UndeployCmdImp,
}

func UndeployCmdImp(cmd *cobra.Command, args []string) {
	// Set all the parameters passed via the command to the struct of undeploy command.
	undeployParams := cmdImp.DeployParams{cmdImp.Verbose, cmdImp.ProjectPath, cmdImp.ManifestPath,
		cmdImp.DeploymentPath, cmdImp.UseDefaults, cmdImp.UseInteractive}
	// Call the implementation of wskdeploy command.
	cmdImp.Undeploy(undeployParams)
}

func init() {
	RootCmd.AddCommand(undeployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// undeployCmd.PersistentFlags().String("foo", "", "A help for foo")
	undeployCmd.PersistentFlags().StringVar(&cmdImp.CfgFile, "config", "", "config file (default is $HOME/.wskdeploy.yaml)")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// undeployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")=
	undeployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	undeployCmd.Flags().StringVarP(&cmdImp.ProjectPath, "pathpath", "p", ".", "path to serverless project")
	undeployCmd.Flags().StringVarP(&cmdImp.ManifestPath, "manifest", "m", "", "path to manifest file")
	undeployCmd.Flags().StringVarP(&cmdImp.DeploymentPath, "deployment", "d", "", "path to deployment file")
}
