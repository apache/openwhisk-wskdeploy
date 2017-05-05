// +build unit

package deployers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var sd *ServiceDeployer
var dr *DeploymentReader
var deployment_file = "../tests/usecases/openstack/deployment.yaml"
var manifest_file = "../tests/usecases/openstack/manifest.yaml"

func init() {
	sd = NewServiceDeployer()
	sd.DeploymentPath = deployment_file
	sd.ManifestPath = manifest_file
	sd.Check()
	dr = NewDeploymentReader(sd)
}

// Check DeploymentReader could handle deployment yaml successfully.
func TestDeploymentReader_HandleYaml(t *testing.T) {
	dr.HandleYaml()
	assert.Equal(t, "JiraBackupSolution", dr.DeploymentDescriptor.Application.Package.Packagename, "DeploymentReader handle deployment yaml failed.")
}

func TestDeployerCheck(t *testing.T) {
	sd := NewServiceDeployer()
	sd.DeploymentPath = "../tests/usecases/badyaml/deployment.yaml"
	sd.ManifestPath = "../tests/usecases/badyaml/manifest.yaml"
	// The system will exit thus the test will fail.
	// sd.Check()
}
