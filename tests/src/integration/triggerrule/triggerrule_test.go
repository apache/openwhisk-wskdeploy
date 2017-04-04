// +build integration

package tests

import (
	"testing"
	"github.com/openwhisk/openwhisk-wskdeploy/cmdImp"
	"github.com/stretchr/testify/assert"
	"os"
)

// TODO: write the integration against openwhisk
func TestTriggerRule(t *testing.T) {
	// Load configuration file

	//Assign the parameters to deploy and undeploy the action
	deployParams := cmdImp.DeployParams{false, ".", manifestPath, deploymentPath, false, false}
	undeployParams := cmdImp.DeployParams{false, ".", manifestPath, deploymentPath, false, false}

	// Deploy the action
	err := cmdImp.Deploy(deployParams)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest and deployment files.")

	// Undeploy the action
	err = cmdImp.Undeploy(undeployParams)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files.")
}

var (
	manifestPath = os.Getenv("GOPATH") + "/src/github.com/openwhisk/openwhisk-wskdeploy/tests/src/integration/triggerrule/manifest.yml"
	deploymentPath = os.Getenv("GOPATH") + "/src/github.com/openwhisk/openwhisk-wskdeploy/tests/src/integration/triggerrule/deployment.yml"
)

