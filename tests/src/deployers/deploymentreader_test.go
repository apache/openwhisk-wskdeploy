package tests

import (
	"github.com/openwhisk/openwhisk-wskdeploy/deployers"
	"github.com/stretchr/testify/assert"
	"testing"
)

var sd *deployers.ServiceDeployer
var dr *deployers.DeploymentReader
var deployment_file = "../../../tests/usecases/openstack/deployment.yaml"
var manifest_file = "../../../tests/usecases/openstack/manifest.yaml"

func init() {
	sd = deployers.NewServiceDeployer()
	sd.DeploymentPath = deployment_file
	sd.ManifestPath = manifest_file
	sd.Check()
	dr = deployers.NewDeploymentReader(sd)
}

// Check DeploymentReader could handle deployment yaml successfully.
func TestDeploymentReader_HandleYaml(t *testing.T) {
	dr.HandleYaml()
	assert.Equal(t, "JiraBackupSolution", dr.DeploymentDescriptor.Application.Package.Packagename, "DeploymentReader handle deployment yaml failed.")
}

func TestDeployerCheck(t *testing.T) {
	sd := deployers.NewServiceDeployer()
	sd.DeploymentPath = "../../../tests/usecases/badyaml/deployment.yaml"
	sd.ManifestPath = "../../../tests/usecases/badyaml/manifest.yaml"
	// The system will exit thus the test will fail.
	// sd.Check()
}
