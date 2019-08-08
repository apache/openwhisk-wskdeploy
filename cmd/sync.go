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
	"github.com/apache/openwhisk-wskdeploy/utils"
	"github.com/apache/openwhisk-wskdeploy/wski18n"
	"github.com/spf13/cobra"
)

// sync represents the mechanism to sync OpenWhisk projects between client and server
var syncCmd = &cobra.Command{
	Use:        "sync",
	SuggestFor: []string{"update"},
	Short:      wski18n.T(wski18n.ID_CMD_DESC_SHORT_SYNC),
	Long:       wski18n.T(wski18n.ID_CMD_DESC_LONG_SYNC),
	RunE:       SyncCmdImp,
}

func SyncCmdImp(cmd *cobra.Command, args []string) error {
	utils.Flags.Sync = true
	return Deploy(cmd)
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
