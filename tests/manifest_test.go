package tests

import (
	"github.com/openwhisk/wskdeploy/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var manifest_yaml = "testcases/helloworld/manifest.yaml"
var expectedPackagename = "helloworld"
var expectedActionLocation = "src/greeting.js"
var expectedActionNumber = 2
var expectedInputParamNumber = 2
var expectedParamType = "string"

func TestParseManifestYAML(t *testing.T) {
	data, err := ioutil.ReadFile(manifest_yaml)
	if err != nil {
		panic(err)
	}

	var manifest utils.ManifestYAML
	err = utils.NewManifestManager().Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}
	//get and verify package name
	assert.Equal(t, expectedPackagename, manifest.Package.Packagename, "Get package name failed.")

	count := 0
	for action_name := range manifest.Package.Actions {
		var action = manifest.Package.Actions[action_name]
		//get and verify action location
		//assert.Equal(t, expectedActionLocation, action.Location, "Get action location failed.")
		//get and verify total param number
		assert.Equal(t, expectedInputParamNumber, len(action.Inputs), "Get input param number failed.")
		//get and verify param type
		for param_name := range action.Inputs {
			var param = action.Inputs[param_name]
			switch param.(type) {
			case string:
				assert.Equal(t, expectedParamType, param, "Get input param type failed.")
			case map[string]string:
				var param_map = param.(map[string]string)
				assert.Equal(t, expectedParamType, param_map["type"], "Get input param type failed.")
			}
		}
		count++
	}
	//get and verify action count
	assert.Equal(t, expectedActionNumber, count, "Get action number failed.")

}
