// +build unit

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

package parsers

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var manifest_yaml = "../tests/usecases/helloworld/manifest.yaml"
var manifestfile1 = "../tests/dat/manifest1.yaml"
var manifestfile3 = "../tests/dat/manifest3.yaml"
var manifestfile4 = "../tests/dat/manifest4.yaml"
var manifestfile5 = "../tests/dat/manifest5.yaml"
var manifestfile6 = "../tests/dat/manifest6.yaml"
var testfile1 = "../tests/dat/deploy1.yaml"
var testfile2 = "../tests/dat/deploy2.yaml"
var testfile3 = "../tests/dat/deploy3.yaml"
var testfile4 = "../tests/dat/deploy4.yaml"

func TestParseManifestYAML(t *testing.T) {
	data, err := ioutil.ReadFile(manifest_yaml)
	if err != nil {
		panic(err)
	}

	var manifest ManifestYAML
	err = NewYAMLParser().Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}
	//get and verify package name
	assert.Equal(t, "helloworld", manifest.Package.Packagename, "Get package name failed.")

	count := 0
	for action_name := range manifest.Package.Actions {
		var action = manifest.Package.Actions[action_name]
		//get and verify action location
		//assert.Equal(t, "src/greeting.js", action.Location, "Get action location failed.")
		//get and verify total param number
		assert.Equal(t, 2, len(action.Inputs), "Get input param number failed.")
		count++
	}
	//get and verify action count
	assert.Equal(t, 2, count, "Get action number failed.")

}

func TestParseManifestYAML_trigger(t *testing.T) {
	data, err := ioutil.ReadFile(manifestfile3)
	if err != nil {
		panic(err)
	}

	var manifest ManifestYAML
	err = NewYAMLParser().Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 2, len(manifest.Package.Triggers), "Get trigger list failed.")
	for trigger_name := range manifest.Package.Triggers {
		var trigger = manifest.Package.Triggers[trigger_name]
		switch trigger_name {
		case "trigger1":
		case "trigger2":
			assert.Equal(t, "myfeed", trigger.Feed, "Get trigger feed name failed.")
		default:
			t.Error("Get trigger name failed")
		}
	}
}

func TestParseManifestYAML_rule(t *testing.T) {
	data, err := ioutil.ReadFile(manifestfile4)
	if err != nil {
		panic(err)
	}

	var manifest ManifestYAML
	err = NewYAMLParser().Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 1, len(manifest.Package.Rules), "Get trigger list failed.")
	for rule_name := range manifest.Package.Rules {
		var rule = manifest.Package.Rules[rule_name]
		switch rule_name {
		case "rule1":
			assert.Equal(t, "trigger1", rule.Trigger, "Get trigger name failed.")
			assert.Equal(t, "hellpworld", rule.Action, "Get action name failed.")
			assert.Equal(t, "true", rule.Rule, "Get rule expression failed.")
		default:
			t.Error("Get rule name failed")
		}
	}
}

func TestParseManifestYAML_feed(t *testing.T) {
	data, err := ioutil.ReadFile(manifestfile5)
	if err != nil {
		panic(err)
	}

	var manifest ManifestYAML
	err = NewYAMLParser().Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 1, len(manifest.Package.Feeds), "Get feed list failed.")
	for feed_name := range manifest.Package.Feeds {
		var feed = manifest.Package.Feeds[feed_name]
		switch feed_name {
		case "feed1":
			assert.Equal(t, "https://my.company.com/services/eventHub", feed.Location, "Get feed location failed.")
			assert.Equal(t, "my_credential", feed.Credential, "Get feed credential failed.")
			assert.Equal(t, 2, len(feed.Operations), "Get operations number failed.")
			for operation_name := range feed.Operations {
				switch operation_name {
				case "operation1":
				case "operation2":
				default:
					t.Error("Get feed operation name failed")
				}
			}
		default:
			t.Error("Get feed name failed")
		}
	}
}

func TestParseManifestYAML_param(t *testing.T) {
	data, err := ioutil.ReadFile(manifestfile6)
	if err != nil {
		panic(err)
	}

	var manifest ManifestYAML
	err = NewYAMLParser().Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 1, len(manifest.Package.Actions), "Get action list failed.")
	for action_name := range manifest.Package.Actions {
		var action = manifest.Package.Actions[action_name]
		switch action_name {
		case "action1":
			for param_name := range action.Inputs {
				var param = action.Inputs[param_name]
				switch param_name {
				case "inline1":
					assert.Equal(t, "{ \"key\": true }", param.Value, "Get param value failed.")
				case "inline2":
					assert.Equal(t, "Just a string", param.Value, "Get param value failed.")
				case "inline3":
					assert.Equal(t, nil, param.Value, "Get param value failed.")
				case "inline4":
					assert.Equal(t, true, param.Value, "Get param value failed.")
				case "inline5":
					assert.Equal(t, 42, param.Value, "Get param value failed.")
				case "inline6":
					assert.Equal(t, -531, param.Value, "Get param value failed.")
				case "inline7":
					assert.Equal(t, 432.432E-43, param.Value, "Get param value failed.")
				case "inline8":
					assert.Equal(t, "[ true, null, \"boo\", { \"key\": 0 }]", param.Value, "Get param value failed.")
				case "inline9":
					assert.Equal(t, false, param.Value, "Get param value failed.")
				case "inline0":
					assert.Equal(t, 456.423, param.Value, "Get param value failed.")
				case "inlin10":
					assert.Equal(t, nil, param.Value, "Get param value failed.")
				case "inlin11":
					assert.Equal(t, true, param.Value, "Get param value failed.")
				case "expand1":
					assert.Equal(t, nil, param.Value, "Get param value failed.")
				case "expand2":
					assert.Equal(t, true, param.Value, "Get param value failed.")
				case "expand3":
					assert.Equal(t, false, param.Value, "Get param value failed.")
				case "expand4":
					assert.Equal(t, 15646, param.Value, "Get param value failed.")
				case "expand5":
					assert.Equal(t, "{ \"key\": true }", param.Value, "Get param value failed.")
				case "expand6":
					assert.Equal(t, "[ true, null, \"boo\", { \"key\": 0 }]", param.Value, "Get param value failed.")
				case "expand7":
					assert.Equal(t, nil, param.Value, "Get param value failed.")
				default:
					t.Error("Get param name failed")
				}
			}
		default:
			t.Error("Get action name failed")
		}
	}
}

func TestParseDeploymentYAML_Application(t *testing.T) {
	//var deployment utils.DeploymentYAML
	mm := NewYAMLParser()
	deployment := mm.ParseDeployment(testfile1)

	//get and verify application name
	assert.Equal(t, "wskdeploy-samples", deployment.Application.Name, "Get application name failed.")
	assert.Equal(t, "/wskdeploy/samples/", deployment.Application.Namespace, "Get application namespace failed.")
	assert.Equal(t, "user-credential", deployment.Application.Credential, "Get application credential failed.")
	assert.Equal(t, "https://172.17.0.1/api", deployment.Application.BaseUrl, "Get application base url failed.")
}

func TestParseDeploymentYAML_Package(t *testing.T) {
	//var deployment utils.DeploymentYAML
	mm := NewYAMLParser()
	deployment := mm.ParseDeployment(testfile2)

	assert.Equal(t, 1, len(deployment.Application.Packages), "Get package list failed.")
	for pkg_name := range deployment.Application.Packages {
		assert.Equal(t, "test_package", pkg_name, "Get package name failed.")
		var pkg = deployment.Application.Packages[pkg_name]
		assert.Equal(t, "http://abc.com/bbb", pkg.Function, "Get package function failed.")
		assert.Equal(t, "/wskdeploy/samples/test", pkg.Namespace, "Get package namespace failed.")
		assert.Equal(t, "12345678ABCDEF", pkg.Credential, "Get package credential failed.")
		assert.Equal(t, 1, len(pkg.Inputs), "Get package input list failed.")
		//get and verify inputs
		for param_name, param := range pkg.Inputs {
			assert.Equal(t, "value", param.Value, "Get input value failed.")
			assert.Equal(t, "param", param_name, "Get input param name failed.")
		}
	}
}

func TestParseDeploymentYAML_Action(t *testing.T) {
	mm := NewYAMLParser()
	deployment := mm.ParseDeployment(testfile2)

	for pkg_name := range deployment.Application.Packages {

		var pkg = deployment.Application.Packages[pkg_name]
		for action_name := range pkg.Actions {
			assert.Equal(t, "hello", action_name, "Get action name failed.")
			var action = pkg.Actions[action_name]
			assert.Equal(t, "/wskdeploy/samples/test/hello", action.Namespace, "Get action namespace failed.")
			assert.Equal(t, "12345678ABCDEF", action.Credential, "Get action credential failed.")
			assert.Equal(t, 1, len(action.Inputs), "Get package input list failed.")
			//get and verify inputs
			for param_name, param := range action.Inputs {
				switch param.Value.(type) {
				case string:
					assert.Equal(t, "name", param_name, "Get input param name failed.")
					assert.Equal(t, "Bernie", param.Value, "Get input value failed.")
				default:
					t.Error("Get input value type failed.")
				}
			}
		}
	}
}

func TestComposeWskPackage(t *testing.T) {
	mm := NewYAMLParser()
	deployment := mm.ParseDeployment(testfile2)
	manifest := mm.ParseManifest(manifestfile1)

	pkglist := deployment.Application.GetPackageList()
	for _, pkg := range pkglist {
		wskpkg := pkg.ComposeWskPackage()
		assert.Equal(t, "test_package", wskpkg.Name, "Get package name failed.")
		assert.Equal(t, "/wskdeploy/samples/test", wskpkg.Namespace, "Get package namespace failed.")
	}

	pkg := manifest.Package
	wskpkg := pkg.ComposeWskPackage()
	assert.Equal(t, "helloworld", wskpkg.Name, "Get package name failed.")
	assert.Equal(t, "1.0", wskpkg.Version, "Get package version failed.")
}

func TestComposeWskTrigger(t *testing.T) {
	mm := NewYAMLParser()
	deployment := mm.ParseDeployment(testfile4)
	manifest := mm.ParseManifest(manifestfile3)

	pkg := deployment.Application.GetPackageList()[0]
	for _, trigger := range pkg.GetTriggerList() {
		//temporarily add the nil to make test pass, as we plan refactor the parser as well as test codes.
		wsktrigger := trigger.ComposeWskTrigger(nil)
		assert.Equal(t, "hello-trigger", wsktrigger.Name, "Get trigger name failed.")
		assert.Equal(t, "/wskdeploy/samples/test/hello-trigger", wsktrigger.Namespace, "Get trigger namespace failed.")
	}

	pkg = manifest.Package
	for _, trigger := range pkg.GetTriggerList() {
		//temporarily add the nil to make test pass, as we plan refactor the parser as well as test codes.
		wsktrigger := trigger.ComposeWskTrigger(nil)
		switch wsktrigger.Name {
		case "trigger1":
		case "trigger2":
		default:
			t.Error("Get trigger name failed")
		}
	}
}

func TestComposeWskRule(t *testing.T) {
	mm := NewYAMLParser()
	manifest := mm.ParseManifest(manifestfile4)

	pkg := manifest.Package
	for _, rule := range pkg.GetRuleList() {
		wskrule := rule.ComposeWskRule()
		switch wskrule.Name {
		case "rule1":
			assert.Equal(t, "trigger1", wskrule.Trigger, "Get rule trigger failed.")
			assert.Equal(t, "hellpworld", wskrule.Action, "Get rule action failed.")
		default:
			t.Error("Get rule name failed")
		}
	}
}
