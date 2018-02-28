// +build integration

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

package tests

import (
	"github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestDefaultPackage(t *testing.T) {
	path := "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/defaultpackage/"
	manifestPath := os.Getenv("GOPATH") + path + "manifest.yaml"
	wskdeploy := common.NewWskdeploy()
	_, err := wskdeploy.DeployManifestPathOnly(manifestPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest file.")
	_, err = wskdeploy.UndeployManifestPathOnly(manifestPath)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest file.")

	time.Sleep(10 * time.Second)

	manifestPath = os.Getenv("GOPATH") + path + "manifest-with-project.yaml"
	_, err = wskdeploy.DeployManifestPathOnly(manifestPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest file.")
	_, err = wskdeploy.UndeployManifestPathOnly(manifestPath)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest file.")

}
