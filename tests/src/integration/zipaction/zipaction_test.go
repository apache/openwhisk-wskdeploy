// +build integration

package tests

import (
	"github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)


var wskprops = common.GetWskprops()


func TestZipAction(t *testing.T) {
	os.Setenv("__OW_API_HOST", wskprops.APIHost)
	wskdeploy := common.NewWskdeploy()
	_, err := wskdeploy.Deploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest and deployment files.")
	_, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files.")
}

var (
	manifestPath   = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/zipaction/manifest.yml"
	deploymentPath = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/zipaction/deployment.yml"
)
