package tests

import (
	"github.com/openwhisk/openwhisk-wskdeploy/deployers"
	"github.com/openwhisk/openwhisk-wskdeploy/parsers"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mr *deployers.ManifestReader
var ps *parsers.YAMLParser
var ms *parsers.ManifestYAML

func init() {

	sd = deployers.NewServiceDeployer()
	sd.ManifestPath = manifest_file
	mr = deployers.NewManfiestReader(sd)
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
	ms := ps.ParseManifest("../../dat/manifest6.yaml")
	err := mr.InitRootPackage(ps, ms)
	assert.Equal(t, err, nil, "Init Root Package failed")

	// TODO.
}
