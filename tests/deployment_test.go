package tests

import (
	"github.com/openwhisk/wsktool/model"
	"io/ioutil"
	"testing"
)

var deployment_yaml = "dat/deployment.yaml"

func TestParseDeploymentYAML(t *testing.T) {
	data, err := ioutil.ReadFile(deployment_yaml)
	if err != nil {
		panic(err)
	}

	var deployment model.DeploymentYAML
	err = model.Deployer.Unmarshal(data, &deployment)
	if err != nil {
		panic(err)
	}
	var expectedPackagename = "helloworld"
	if deployment.Package.Packagename != expectedPackagename {
		t.Error("get package name failed")
	}

	var expectedDesc = "input person address"
	if deployment.Package.Actions[0].Input[1]["description"] != expectedDesc {
		t.Error("get input failed")
	}

}
