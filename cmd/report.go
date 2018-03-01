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
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/deployers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
	"github.com/spf13/cobra"
	"os"
	"path"
	"sync"
)

var wskpropsPath string

var client *whisk.Client
var wg sync.WaitGroup

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:        "report",
	SuggestFor: []string{"list"},
	Short:      wski18n.T(wski18n.ID_CMD_DESC_SHORT_REPORT),
	RunE: func(cmd *cobra.Command, args []string) error {
		if wskpropsPath != "" {
			config, _ := deployers.NewWhiskConfig(wskpropsPath, utils.Flags.DeploymentPath, utils.Flags.ManifestPath, false)
			client, _ := deployers.CreateNewClient(config)
			return printDeploymentInfo(client)
		} else {
			//default to ~/.wskprops
			userHome := utils.GetHomeDirectory()
			// TODO() we should not only use const. for config files like .wskprops, but have a dedicated
			// set of functions in its own package to interact with it as a resource
			propPath := path.Join(userHome, ".wskprops")
			config, _ := deployers.NewWhiskConfig(propPath, utils.Flags.DeploymentPath, utils.Flags.ManifestPath, false)
			client, _ := deployers.CreateNewClient(config)
			return printDeploymentInfo(client)
		}
	},
}

func init() {
	RootCmd.AddCommand(reportCmd)

	// TODO() REMOVE this flag... the flag -config exists already
	reportCmd.Flags().StringVarP(&wskpropsPath, "wskproppath", "w",
		path.Join(os.Getenv("HOME"), ".wskprops"),
		wski18n.T(wski18n.ID_CMD_FLAG_CONFIG))
}

func printDeploymentInfo(client *whisk.Client) error {
	//We currently list packages, actions, triggers, rules.
	wg.Add(4)

	wskprint.PrintlnOpenWhiskInfo(wski18n.T(wski18n.ID_MSG_DEPLOYMENT_REPORT))
	// we set the default package list options
	pkgoptions := &whisk.PackageListOptions{false, 0, 0, 0, false}
	packages, _, err := client.Packages.List(pkgoptions)
	if err != nil {
		return err
	}

	// list all packages under current namespace.
	go func() {
		defer wg.Done()
		printList(packages)
	}()

	// list all the actions under all the packages.
	go func() error {
		defer wg.Done()
		acnoptions := &whisk.ActionListOptions{0, 0, false}
		for _, pkg := range packages {
			actions, _, err := client.Actions.List(pkg.Name, acnoptions)
			if err != nil {
				return err
			}
			printActionList(actions)
		}
		return nil
	}()

	// list all the triggers under current namespace.
	go func() error {
		defer wg.Done()
		troptions := &whisk.TriggerListOptions{0, 0, false}
		_, _, err := client.Triggers.List(troptions)
		if err != nil {
			return err
		}
		//printTriggerList(triggers)
		return nil
	}()

	// list all the rules under current namespace.
	go func() error {
		defer wg.Done()
		roptions := &whisk.RuleListOptions{0, 0, false}
		rules, _, err := client.Rules.List(roptions)
		if err != nil {
			return err
		}
		printRuleList(rules)
		return nil
	}()

	wg.Wait()

	return nil
}

// From below the codes are from the whisk-client package.
// http://github.com/openwhisk/openwhisk.git
func printList(collection interface{}) {
	switch collection := collection.(type) {
	case []whisk.Action:
		printActionList(collection)
	// TODO()
	//case []whisk.Trigger
	//	printTriggerList(collection)
	case []whisk.Package:
		printPackageList(collection)
	case []whisk.Rule:
		printRuleList(collection)
	case []whisk.Namespace:
		printNamespaceList(collection)
	case []whisk.Activation:
		printActivationList(collection)
	}
}

// TODO() i18n private / shared never translated
func printRuleList(rules []whisk.Rule) {
	wskprint.PrintlnOpenWhiskInfoTitle(wski18n.RULES)

	for _, rule := range rules {
		publishState := wski18n.T("private")
		if *rule.Publish {
			publishState = wski18n.T("shared")
		}
		output := fmt.Sprintf("%-70s %s\n",
			fmt.Sprintf("/%s/%s", rule.Namespace, rule.Name),
			publishState)
		wskprint.PrintlnOpenWhiskInfo(output)
	}
}

// TODO() i18n private / shared never translated
func printPackageList(packages []whisk.Package) {
	wskprint.PrintlnOpenWhiskInfoTitle(wski18n.PACKAGES)
	for _, xPackage := range packages {
		publishState := wski18n.T("private")
		if *xPackage.Publish {
			publishState = wski18n.T("shared")
		}
		output := fmt.Sprintf("%-70s %s\n",
			fmt.Sprintf("/%s/%s", xPackage.Namespace, xPackage.Name),
			publishState)
		wskprint.PrintlnOpenWhiskInfo(output)
	}
}

// TODO() i18n private / shared never translated
func printActionList(actions []whisk.Action) {
	wskprint.PrintlnOpenWhiskInfoTitle(wski18n.ACTIONS)
	for _, action := range actions {
		publishState := wski18n.T("private")
		if *action.Publish {
			publishState = wski18n.T("shared")
		}
		kind := getValueString(action.Annotations, "exec")
		output := fmt.Sprintf("%-70s %s %s\n",
			fmt.Sprintf("/%s/%s", action.Namespace, action.Name),
			publishState,
			kind)
		wskprint.PrintlnOpenWhiskInfo(output)
	}
}

/*
func printTriggerList(triggers whisk.Trigger) {
	//fmt.Fprintf(color.Output, "%s\n", boldString("triggers"))
	wskprint.PrintlnOpenWhiskInfoTitle(wski18n.TRIGGERS)
	for _, trigger := range triggers {
		publishState := wski18n.T("private")
		if trigger.Publish {
			publishState = wski18n.T("shared")
		}
		fmt.Printf("%-70s %s\n", fmt.Sprintf("/%s/%s", trigger.Namespace, trigger.Name), publishState)
	}
}*/

func getValueString(keyValueArr whisk.KeyValueArr, key string) string {
	var value interface{}
	var res string

	value = getValue(keyValueArr, key)
	castedValue, canCast := value.(string)

	if canCast {
		res = castedValue
	}

	// TODO() This may be too much for end-user debug/trace
	//dbgMsg := fmt.Sprintf("keyValueArr[%v]: key=[%s] value=[%v]\n",  keyValueArr, key, res)
	//wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, dbgMsg )

	return res
}

func getValue(keyValueArr whisk.KeyValueArr, key string) interface{} {
	var res interface{}

	for i := 0; i < len(keyValueArr); i++ {
		if keyValueArr[i].Key == key {
			res = keyValueArr[i].Value
			break
		}
	}

	// TODO() This may be too much for end-user debug/trace
	//dbgMsg := fmt.Sprintf("keyValueArr[%v]: key=[%s] value=[%v]\n",  keyValueArr, key, res)
	//wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, dbgMsg )

	return res
}

func printNamespaceList(namespaces []whisk.Namespace) {
	wskprint.PrintlnOpenWhiskInfo(wski18n.NAMESPACES)
	for _, namespace := range namespaces {
		output := fmt.Sprintf("%s\n", namespace.Name)
		wskprint.PrintlnOpenWhiskInfo(output)
	}
}

func printActivationList(activations []whisk.Activation) {
	wskprint.PrintlnOpenWhiskInfo(wski18n.ACTIVATIONS)
	for _, activation := range activations {
		output := fmt.Sprintf("%s %20s\n", activation.ActivationID, activation.Name)
		wskprint.PrintlnOpenWhiskInfo(output)
	}
}
