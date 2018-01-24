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
	"os"
	"github.com/spf13/cobra"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:	"add",
	SuggestFor: []string {"insert"},
	Short:	wski18n.T(wski18n.ID_CMD_DESC_SHORT_ADD),
}

// action represents the `add action` command
var actionCmd = &cobra.Command{
	Use:   "action",
	Short: wski18n.T(wski18n.ID_CMD_DESC_SHORT_ADD_X_key_X,
		map[string]interface{}{wski18n.KEY_KEY: parsers.YAML_KEY_ACTION}),
	RunE: func(cmd *cobra.Command, args []string) error {
		maniyaml, err := parsers.ReadOrCreateManifest()
        if err != nil {
            return err
        }

	reader := bufio.NewReader(os.Stdin)
	action := parsers.Action{}

	for {
		action.Name = utils.Ask(reader, wski18n.NAME_ACTION, "")

		// Check action name is unique
		if _, ok := maniyaml.Package.Actions[action.Name]; !ok {
			break
		}

		warnMsg := wski18n.T(wski18n.ID_WARN_ENTITY_NAME_EXISTS_X_key_X_name_X,
			map[string]interface{}{
				wski18n.KEY_KEY: parsers.YAML_KEY_ACTION,
				wski18n.KEY_NAME: action.Name})
		wskprint.PrintOpenWhiskWarning(warnMsg)
	}

	// TODO() use dynamic/programmatic way to get default runtime (not hardcoded)
	// TODO() And List all supported runtime names (values) (via API)
	action.Runtime = utils.Ask(reader, wski18n.NAME_RUNTIME, "nodejs:6")
	maniyaml.Package.Actions[action.Name] = action

	// Create directory structure before update manifest, as a way
	// to check the action name is a valid path name
	err = os.MkdirAll("actions/"+action.Name, 0777)

        if err != nil {
            return err
        }

		return parsers.Write(maniyaml, utils.ManifestFileNameYaml)
	},
}

// trigger represents the `add trigger` command
var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: wski18n.T(wski18n.ID_CMD_DESC_SHORT_ADD_X_key_X,
		map[string]interface{}{wski18n.KEY_KEY: parsers.YAML_KEY_TRIGGER}),
	RunE: func(cmd *cobra.Command, args []string) error {
		maniyaml, err := parsers.ReadOrCreateManifest()
        if err != nil {
            return err
        }

		reader := bufio.NewReader(os.Stdin)
		trigger := parsers.Trigger{}

		for {
			trigger.Name = utils.Ask(reader, wski18n.NAME_TRIGGER, "")

			// Check trigger name is unique
			if _, ok := maniyaml.Package.Triggers[trigger.Name]; !ok {
				break
			}

			warnMsg := wski18n.T(wski18n.ID_WARN_ENTITY_NAME_EXISTS_X_key_X_name_X,
				map[string]interface{}{
					wski18n.KEY_KEY: parsers.YAML_KEY_TRIGGER,
					wski18n.KEY_NAME: trigger.Name})
			wskprint.PrintOpenWhiskWarning(warnMsg)
		}

		trigger.Feed = utils.Ask(reader, wski18n.NAME_FEED, "")
		maniyaml.Package.Triggers[trigger.Name] = trigger

		return parsers.Write(maniyaml, utils.ManifestFileNameYaml)
	},
}

// rule represents the `add rule` command
var ruleCmd = &cobra.Command{
	Use:   "rule",
	Short: wski18n.T(wski18n.ID_CMD_DESC_SHORT_ADD_X_key_X,
		map[string]interface{}{wski18n.KEY_KEY: parsers.YAML_KEY_RULE}),
	RunE: func(cmd *cobra.Command, args []string) error {
		maniyaml, err := parsers.ReadOrCreateManifest()
        if err != nil {
            return err
        }

		reader := bufio.NewReader(os.Stdin)
		rule := parsers.Rule{}

		for {
			rule.Rule = utils.Ask(reader, wski18n.NAME_RULE, "")

			// Check rule name is unique
			if _, ok := maniyaml.Package.Triggers[rule.Rule]; !ok {
				break
			}

			warnMsg := wski18n.T(wski18n.ID_WARN_ENTITY_NAME_EXISTS_X_key_X_name_X,
				map[string]interface{}{
					wski18n.KEY_KEY: parsers.YAML_KEY_RULE,
					wski18n.KEY_NAME: rule.Name})
			wskprint.PrintOpenWhiskWarning(warnMsg)
		}

		rule.Action = utils.Ask(reader, wski18n.NAME_ACTION, "")
		rule.Trigger = utils.Ask(reader, wski18n.NAME_TRIGGER, "")
		maniyaml.Package.Rules[rule.Rule] = rule

		return parsers.Write(maniyaml, utils.ManifestFileNameYaml)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.AddCommand(actionCmd)
	addCmd.AddCommand(triggerCmd)
	addCmd.AddCommand(ruleCmd)
}
