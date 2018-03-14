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
	"os"
	"testing"

	"github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
)

/* *
 * Please configure BLUEMIX_APIHOST, BLUEMIX_NAMESPACE and BLUEMIX_AUTH on your local machine in order to run this
 * integration test.
 */
func TestRelationships(t *testing.T) {
	wskdeploy := common.NewWskdeploy()

	_, err := wskdeploy.ManagedDeployment(manifestLibPath, deploymentLibPath)
	assert.Equal(t, nil, err, "Failed to deploy the lib manifest file.")
	_, err = wskdeploy.ManagedRelationshipsDeployment(manifestExtPath, relationshipsPath)
	assert.Equal(t, nil, err, "Failed to deploy the ext manifest file.")
	_, err = wskdeploy.ManagedRelationshipsUnDeployment(manifestExtPath, relationshipsPath)
	assert.Equal(t, nil, err, "Failed to undeploy the ext manifest file.")
	_, err = wskdeploy.ManagedUndeployment(manifestLibPath, deploymentLibPath)
	assert.Equal(t, nil, err, "Failed to undeploy the lib manifest file.")
}

var (
	manifestLibPath   = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/relationships/manifest.yml"
	manifestExtPath   = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/relationships/ext_manifest.yml"
	deploymentLibPath = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/relationships/deployment.yml"
	relationshipsPath = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/relationships/relationships.yml"
)
