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
)

// undeployCmd represents the undeploy command
// TODO() i18n the short/long descriptions
var undeployCmd = &cobra.Command{
	Use:        "undeploy",
	SuggestFor: []string{"remove"},
	Short:      "Undeploy assets from OpenWhisk",
	Long:       `Undeploy removes deployed assets from the manifest and deployment files`,
	RunE:       UndeployCmdImp,
}

func UndeployCmdImp(cmd *cobra.Command, args []string) error {
	return Undeploy(cmd)
}

func init() {
	RootCmd.AddCommand(undeployCmd)
}
