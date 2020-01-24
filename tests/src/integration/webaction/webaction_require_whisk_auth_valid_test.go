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
	"github.com/apache/openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestRequireWhiskAuthAnnotation(t *testing.T) {
	wskdeploy := common.NewWskdeploy()
	_, err := wskdeploy.DeployManifestPathOnly(manifestPathValidTests)
	assert.Equal(t, nil, err, "Failed to deploy 'require-whisk-auth' annotations based on the manifest file.")

	// artificial 1 second delay to allow API/swagger creation
	time.Sleep(1 * time.Second)

	_, err2 := wskdeploy.UndeployManifestPathOnly(manifestPathValidTests)
	assert.Equal(t, nil, err2, "Failed to undeploy 'require-whisk-auth' annotations based on the manifest file.")
}

var (
	manifestPathValidTests = os.Getenv("GOPATH") + "/src/github.com/apache/openwhisk-wskdeploy/tests/src/integration/webaction/manifest_require_whisk_auth_valid.yaml"
)
