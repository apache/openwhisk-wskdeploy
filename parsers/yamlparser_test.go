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
	"testing"
)

var manifestfile_val_pkg = "../tests/dat/manifest_validate_package_grammar.yaml"
var manifestfile_val_tar = "../tests/dat/manifest_validate_trigger_action_rule_grammar.yaml"
var manifestfile3 = "../tests/dat/manifest3.yaml"
var manifestfile4 = "../tests/dat/manifest4.yaml"
var manifestfile5 = "../tests/dat/manifest5.yaml"
var manifestfile6 = "../tests/dat/manifest6.yaml"
var deploymentfile_data_app = "../tests/dat/deployment_data_application.yaml"
var deploymentfile_data_app_pkg = "../tests/dat/deployment_data_application_package.yaml"
var testfile3 = "../tests/dat/deploy3.yaml"
var testfile4 = "../tests/dat/deploy4.yaml"


func TestComposeWskPackage(t *testing.T) {
	mm := NewYAMLParser()
    deployment, _ := mm.ParseDeployment(deploymentfile_data_app_pkg)
	manifest, _ := mm.ParseManifest(manifestfile_val_pkg)

	pkglist := deployment.Application.GetPackageList()
	for _, pkg := range pkglist {
		wskpkg := pkg.ComposeWskPackage()
		assert.Equal(t, "test_package", wskpkg.Name, "Get package name failed.")
		assert.Equal(t, "/wskdeploy/samples/test", wskpkg.Namespace, "Get package namespace failed.")
	}

	for n, p := range manifest.Packages{
		wskpkg := p.ComposeWskPackage()
		assert.Equal(t, "helloworld", n, "Get package name failed.")
		assert.Equal(t, "1.0", wskpkg.Version, "Get package version failed.")
	}
}

func TestComposeWskTrigger(t *testing.T) {
	mm := NewYAMLParser()
    deployment, _ := mm.ParseDeployment(testfile4)
	manifest, _ := mm.ParseManifest(manifestfile3)

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
	manifest, _ := mm.ParseManifest(manifestfile4)

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

func TestGetActionList(t *testing.T) {
	mm := NewYAMLParser()
	manifest, _ := mm.ParseManifest(manifestfile_val_tar)
	pkg := manifest.Packages["manifest2"]
	actions := pkg.GetActionList()
	assert.Equal(t,3, len(actions), "Get action list failed.")
}

func TestGetTriggerList(t *testing.T) {
	mm := NewYAMLParser()
	manifest, _ := mm.ParseManifest(manifestfile_val_tar)
	pkg := manifest.Packages["manifest2"]
	triggers := pkg.GetTriggerList()
	assert.Equal(t,2, len(triggers), "Get trigger list failed.")
}

func TestGetRuleList(t *testing.T) {
	mm := NewYAMLParser()
	manifest, _ := mm.ParseManifest(manifestfile_val_tar)
	pkg := manifest.Packages["manifest2"]
	rules := pkg.GetRuleList()
	assert.Equal(t,3, len(rules), "Get trigger list failed.")
}

func TestGetFeedList(t *testing.T) {
	mm := NewYAMLParser()
	manifest, _ := mm.ParseManifest(manifestfile_val_tar)
	pkg := manifest.Packages["manifest2"]
	feeds := pkg.GetFeedList()
	assert.Equal(t,4, len(feeds), "Get feed list failed.")
}

func TestGetApisList(t *testing.T) {
	mm := NewYAMLParser()
	manifest, _ := mm.ParseManifest(manifestfile_val_tar)
	pkg := manifest.Packages["manifest2"]
	apis := pkg.GetApis()
	assert.Equal(t,5, len(apis), "Get api list failed.")
}
