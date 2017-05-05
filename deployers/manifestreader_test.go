// +build unit

package deployers

import (
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mr *ManifestReader
var ps *parsers.YAMLParser
var ms *parsers.ManifestYAML

func init() {

	sd = NewServiceDeployer()
	sd.ManifestPath = manifest_file
	mr = NewManfiestReader(sd)
	ps = parsers.NewYAMLParser()
	ms = ps.ParseManifest(manifest_file)
}

// Test could parse Manifest file successfully
func TestManifestReader_ParseManifest(t *testing.T) {
	_, _, err := mr.ParseManifest()
	assert.Equal(t, err, nil, "New ManifestReader failed")
}

// Test could Init root package successfully.
func TestManifestReader_InitRootPackage(t *testing.T) {
	err := mr.InitRootPackage(ps, ms)
	assert.Equal(t, err, nil, "Init Root Package failed")
}

// Test Parameters
func TestManifestReader_param(t *testing.T) {
	ms := ps.ParseManifest("../tests/dat/manifest6.yaml")
	err := mr.InitRootPackage(ps, ms)
	assert.Equal(t, err, nil, "Init Root Package failed")

	// TODO.
}
