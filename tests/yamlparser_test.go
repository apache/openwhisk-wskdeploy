package tests

import (
	"github.com/openwhisk/wskdeploy/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var manifest_yaml = "testcases/helloworld/manifest.yaml"
var manifestfile1 = "dat/manifest1.yaml"
var manifestfile3 = "dat/manifest3.yaml"
var manifestfile4 = "dat/manifest4.yaml"
var testfile1 = "dat/deploy1.yaml"
var testfile2 = "dat/deploy2.yaml"
var testfile3 = "dat/deploy3.yaml"
var testfile4 = "dat/deploy4.yaml"

func TestParseManifestYAML(t *testing.T) {
	data, err := ioutil.ReadFile(manifest_yaml)
	if err != nil {
		panic(err)
	}

	var manifest utils.ManifestYAML
	err = utils.NewYAMLParser().Unmarshal(data, &manifest)
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
		//get and verify param type
		for param_name := range action.Inputs {
			var param = action.Inputs[param_name]
			switch param.(type) {
			case string:
				assert.Equal(t, "string", param, "Get input param type failed.")
			case map[string]string:
				var param_map = param.(map[string]string)
				assert.Equal(t, "string", param_map["type"], "Get input param type failed.")
			}
		}
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

	var manifest utils.ManifestYAML
	err = utils.NewYAMLParser().Unmarshal(data, &manifest)
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

	var manifest utils.ManifestYAML
	err = utils.NewYAMLParser().Unmarshal(data, &manifest)
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

func TestParseDeploymentYAML_Application(t *testing.T) {
	//var deployment utils.DeploymentYAML
	mm := utils.NewYAMLParser()
	deployment := mm.ParseDeployment(testfile1)

	//get and verify application name
	assert.Equal(t, "wskdeploy-samples", deployment.Application.Name, "Get application name failed.")
	assert.Equal(t, "/wskdeploy/samples/", deployment.Application.Namespace, "Get application namespace failed.")
}

func TestParseDeploymentYAML_Package(t *testing.T) {
	//var deployment utils.DeploymentYAML
	mm := utils.NewYAMLParser()
	deployment := mm.ParseDeployment(testfile2)

	assert.Equal(t, 1, len(deployment.Application.Packages), "Get package list failed.")
	for pkg_name := range deployment.Application.Packages {
		assert.Equal(t, "test_package", pkg_name, "Get package name failed.")
		var pkg = deployment.Application.Packages[pkg_name]
		assert.Equal(t, "http://abc.com/bbb", pkg.Function, "Get package function failed.")
		assert.Equal(t, "12345678ABCDEF", pkg.PackageCredential, "Get package credential failed.")
		assert.Equal(t, "/wskdeploy/samples/test", pkg.Namespace, "Get package namespace failed.")
		assert.Equal(t, "12345678ABCDEF", pkg.Credential, "Get package credential failed.")
		assert.Equal(t, 1, len(pkg.Inputs), "Get package input list failed.")
		//get and verify inputs
		for param_name, value := range pkg.Inputs {
			assert.Equal(t, "value", value, "Get input value failed.")
			assert.Equal(t, "param", param_name, "Get input param name failed.")
		}
	}
}

func TestParseDeploymentYAML_Action(t *testing.T) {
	mm := utils.NewYAMLParser()
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
			for param_name, value := range action.Inputs {
				switch value.(type) {
				case string:
					assert.Equal(t, "name", param_name, "Get input param name failed.")
					assert.Equal(t, "Bernie", value, "Get input value failed.")
				default:
					t.Error("Get input value type failed.")
				}
			}
		}
	}
}

func TestComposeWskPackage(t *testing.T) {
	mm := utils.NewYAMLParser()
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
	mm := utils.NewYAMLParser()
	deployment := mm.ParseDeployment(testfile4)
	manifest := mm.ParseManifest(manifestfile3)

	pkg := deployment.Application.GetPackageList()[0]
	for _, trigger := range pkg.GetTriggerList() {
		wsktrigger := trigger.ComposeWskTrigger()
		assert.Equal(t, "hello-trigger", wsktrigger.Name, "Get trigger name failed.")
		assert.Equal(t, "/wskdeploy/samples/test/hello-trigger", wsktrigger.Namespace, "Get trigger namespace failed.")
	}

	pkg = manifest.Package
	for _, trigger := range pkg.GetTriggerList() {
		wsktrigger := trigger.ComposeWskTrigger()
		switch wsktrigger.Name {
		case "trigger1":
		case "trigger2":
		default:
			t.Error("Get trigger name failed")
		}
	}
}

func TestComposeWskRule(t *testing.T) {
	mm := utils.NewYAMLParser()
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
