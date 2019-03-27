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
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:        "version",
	SuggestFor: []string{"edition", "release"},
	Short:      wski18n.T(wski18n.ID_CMD_DESC_SHORT_VERSION),
	Run: func(cmd *cobra.Command, args []string) {
		wskprint.PrintlnOpenWhiskOutput(
			// Note: no need to translate the following string
			// TODO(#767) - Flags.CliVersion are not set during build
			fmt.Sprintf("wskdeploy version: %s", utils.Flags.CliVersion))
		wskprint.PrintlnOpenWhiskOutput(
			fmt.Sprintf("wskdeploy git commit: %s", utils.Flags.CliGitCommit))
		wskprint.PrintlnOpenWhiskOutput(
			fmt.Sprintf("wskdeploy build date: %s", utils.Flags.CliBuildDate))
	},
}
