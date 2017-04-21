// +build integration

package tests

import (
	"testing"
	"github.com/openwhisk/openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
	"os"
)

var wskprops = common.GetWskprops()

// TODO: write the integration against openwhisk
func TestTriggerRule(t *testing.T) {
	os.Setenv("__OW_API_HOST", wskprops.APIHost)
	wskdeploy := common.NewWskdeploy()
	_, err := wskdeploy.Deploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest and deployment files.")
	_, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files.")
}

var (
	manifestPath = os.Getenv("GOPATH") + "/src/github.com/openwhisk/openwhisk-wskdeploy/tests/src/integration/triggerrule/manifest.yml"
	deploymentPath = os.Getenv("GOPATH") + "/src/github.com/openwhisk/openwhisk-wskdeploy/tests/src/integration/triggerrule/deployment.yml"
)

