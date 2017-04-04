// +build integration

package tests

import (
	"github.com/openwhisk/openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)


var wskprops = common.GetWskprops()

// TODO: write the integration against openwhisk
func TestDependency(t *testing.T) {
	os.Setenv("__OW_API_HOST", wskprops.APIHost)
	wskdeploy := common.NewWskdeploy()
	_, err := wskdeploy.Deploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest and deployment files.")
	_, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files.")
}

var (
	manifestPath   = os.Getenv("GOPATH") + "/src/github.com/openwhisk/openwhisk-wskdeploy/tests/src/integration/dependency/manifest.yaml"
	deploymentPath = ""
)