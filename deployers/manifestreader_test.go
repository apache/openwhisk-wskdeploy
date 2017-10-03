// +build unit

/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package deployers

import (
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mr *ManifestReader
var ps *parsers.YAMLParser
var ms *parsers.YAML

func init() {

	sd = NewServiceDeployer()
	sd.ManifestPath = manifest_file
	mr = NewManfiestReader(sd)
	ps = parsers.NewYAMLParser()
	ms, _ = ps.ParseManifest(manifest_file)
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
	ms, _ := ps.ParseManifest("../tests/dat/manifest6.yaml")
	err := mr.InitRootPackage(ps, ms)
	assert.Equal(t, err, nil, "Init Root Package failed")

	// TODO.
}
