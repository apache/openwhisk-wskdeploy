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
	"bufio"
	"fmt"
	"os"

	"github.com/openwhisk/openwhisk-wskdeploy/parsers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an action, feed, trigger or rule to the manifest",
}

// action represents the `add action` command
var actionCmd = &cobra.Command{
	Use:   "action",
	Short: "add action to the manifest file and create default directory structure.",
	Run: func(cmd *cobra.Command, args []string) {
		maniyaml := parsers.ReadOrCreateManifest()

		reader := bufio.NewReader(os.Stdin)
		action := parsers.Action{}

		for {
			action.Name = utils.Ask(reader, "Name", "")

			// Check action name is unique
			if _, ok := maniyaml.Package.Actions[action.Name]; !ok {
				break
			}
			fmt.Print(action.Name + " is already used. Pick another action name\n")
		}

		action.Runtime = utils.Ask(reader, "Runtime", "nodejs:6")
		maniyaml.Package.Actions[action.Name] = action

		// Create directory structure before update manifest, as a way
		// to check the action name is a valid path name
		err := os.MkdirAll("actions/"+action.Name, 0777)
		utils.Check(err)

		parsers.Write(maniyaml, "manifest.yaml")
	},
}

// trigger represents the `add trigger` command
var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "add trigger to the manifest file.",
	Run: func(cmd *cobra.Command, args []string) {
		maniyaml := parsers.ReadOrCreateManifest()

		reader := bufio.NewReader(os.Stdin)
		trigger := parsers.Trigger{}

		for {
			trigger.Name = utils.Ask(reader, "Name", "")

			// Check trigger name is unique
			if _, ok := maniyaml.Package.Triggers[trigger.Name]; !ok {
				break
			}
			fmt.Print(trigger.Name + " is already used. Pick another trigger name\n")
		}

		trigger.Feed = utils.Ask(reader, "Feed", "")
		maniyaml.Package.Triggers[trigger.Name] = trigger

		parsers.Write(maniyaml, "manifest.yaml")
	},
}

// rule represents the `add rule` command
var ruleCmd = &cobra.Command{
	Use:   "rule",
	Short: "add rule to the manifest file.",
	Run: func(cmd *cobra.Command, args []string) {
		maniyaml := parsers.ReadOrCreateManifest()

		reader := bufio.NewReader(os.Stdin)
		rule := parsers.Rule{}

		for {
			rule.Rule = utils.Ask(reader, "Rule Name", "")

			// Check rule name is unique
			if _, ok := maniyaml.Package.Triggers[rule.Rule]; !ok {
				break
			}
			fmt.Print(rule.Rule + " is already used. Pick another rule name\n")
		}

		rule.Action = utils.Ask(reader, "Action", "")
		rule.Trigger = utils.Ask(reader, "Trigger", "")
		maniyaml.Package.Rules[rule.Rule] = rule

		parsers.Write(maniyaml, "manifest.yaml")
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.AddCommand(actionCmd)
	addCmd.AddCommand(triggerCmd)
	addCmd.AddCommand(ruleCmd)
}
