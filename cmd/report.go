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
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"path"
	"sync"
)

var wskpropsPath string

var client *whisk.Client
var wg sync.WaitGroup

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Returns summary of what's been deployed on OpenWhisk in specific namespace",
	Long: `Command helps user get an overall report about what's been deployed
on OpenWhisk with specific OpenWhisk namespace. By default it will read the wsk property file
located under current user home.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		if wskpropsPath != "" {
			client, _ = deployers.NewWhiskClient(wskpropsPath, utils.Flags.DeploymentPath, false)
		}
		userHome := utils.GetHomeDirectory()
		//default to ~/.wskprops
		propPath := path.Join(userHome, ".wskprops")
		client, _ = deployers.NewWhiskClient(propPath, utils.Flags.DeploymentPath, false)
		printDeploymentInfo(client)
	},
}

func init() {
	RootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringVarP(&wskpropsPath, "wskproppath", "w", ".", "path to wsk property file, default is to ~/.wskprops")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

var boldString = color.New(color.Bold).SprintFunc()

func printDeploymentInfo(*whisk.Client) error {
	//We currently list packages, actions, triggers, rules.
	wg.Add(4)

	fmt.Println("----==== OpenWhisk Deployment Status ====----")
	// we set the default package list options
	pkgoptions := &whisk.PackageListOptions{false, 0, 0, 0, false}
	packages, _, err := client.Packages.List(pkgoptions)
	utils.Check(err)

	// list all packages under current namespace.
	go func() {
		defer wg.Done()
		printList(packages)
	}()

	// list all the actions under all the packages.
	go func() {
		defer wg.Done()
		acnoptions := &whisk.ActionListOptions{0, 0, false}
		for _, pkg := range packages {
			actions, _, err := client.Actions.List(pkg.Name, acnoptions)
			utils.Check(err)
			printActionList(actions)
		}
	}()

	// list all the triggers under current namespace.
	go func() {
		defer wg.Done()
		troptions := &whisk.TriggerListOptions{0, 0, false}
		_, _, err := client.Triggers.List(troptions)
		utils.Check(err)
		//printTriggerList(triggers)

	}()

	// list all the rules under current namespace.
	go func() {
		defer wg.Done()
		roptions := &whisk.RuleListOptions{0, 0, false}
		rules, _, err := client.Rules.List(roptions)
		utils.Check(err)
		printRuleList(rules)
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

func printRuleList(rules []whisk.Rule) {
	fmt.Fprintf(color.Output, "%s\n", boldString("rules"))
	for _, rule := range rules {
		publishState := wski18n.T("private")
		if *rule.Publish {
			publishState = wski18n.T("shared")
		}
		fmt.Printf("%-70s %s\n", fmt.Sprintf("/%s/%s", rule.Namespace, rule.Name), publishState)
	}
}

func printPackageList(packages []whisk.Package) {
	fmt.Fprintf(color.Output, "%s\n", boldString("packages"))
	for _, xPackage := range packages {
		publishState := wski18n.T("private")
		if *xPackage.Publish {
			publishState = wski18n.T("shared")
		}
		fmt.Printf("%-70s %s\n", fmt.Sprintf("/%s/%s", xPackage.Namespace, xPackage.Name), publishState)
	}
}

func printActionList(actions []whisk.Action) {
	fmt.Fprintf(color.Output, "%s\n", boldString("actions"))
	for _, action := range actions {
		publishState := wski18n.T("private")
		if *action.Publish {
			publishState = wski18n.T("shared")
		}
		kind := getValueString(action.Annotations, "exec")
		fmt.Printf("%-70s %s %s\n", fmt.Sprintf("/%s/%s", action.Namespace, action.Name), publishState, kind)
	}
}

/*
func printTriggerList(triggers whisk.Trigger) {
	fmt.Fprintf(color.Output, "%s\n", boldString("triggers"))
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

	whisk.Debug(whisk.DbgInfo, "Got string value '%v' for key '%s'\n", res, key)

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

	whisk.Debug(whisk.DbgInfo, "Got value '%v' from '%v' for key '%s'\n", res, keyValueArr, key)

	return res
}

func printNamespaceList(namespaces []whisk.Namespace) {
	fmt.Fprintf(color.Output, "%s\n", boldString("namespaces"))
	for _, namespace := range namespaces {
		fmt.Printf("%s\n", namespace.Name)
	}
}

func printActivationList(activations []whisk.Activation) {
	fmt.Fprintf(color.Output, "%s\n", boldString("activations"))
	for _, activation := range activations {
		fmt.Printf("%s %20s\n", activation.ActivationID, activation.Name)
	}
}
