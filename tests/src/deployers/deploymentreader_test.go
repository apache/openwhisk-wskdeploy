package tests

import (
	"github.com/openwhisk/openwhisk-wskdeploy/deployers"
	"github.com/stretchr/testify/assert"
	"testing"
)

var sd *deployers.ServiceDeployer
var dr *deployers.DeploymentReader
var deployment_file = "../../../tests/usecases/openstack/deployment.yaml"

func init() {
	sd = deployers.NewServiceDeployer()
	sd.DeploymentPath = deployment_file
	dr = deployers.NewDeploymentReader(sd)
}

// Check DeploymentReader could handle deployment yaml successfully.
func TestDeploymentReader_HandleYaml(t *testing.T) {
	dr.HandleYaml()
	assert.Equal(t, dr.DeploymentDescriptor.Application.Package.Packagename, "JiraBackupSolution", "DeploymentReader handle deployment yaml failed.")
}
