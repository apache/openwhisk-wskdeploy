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
)


var wskprops = common.GetWskprops()

// TODO: write the integration against openwhisk
func TestDependency(t *testing.T) {
    wskdeploy := common.NewWskdeploy()
    _, err := wskdeploy.Deploy(manifestPath, deploymentPath)
    assert.Equal(t, nil, err, "Failed to deploy based on the manifest and deployment files.")

    wskcli := common.NewWskCLI()
    _, err = wskcli.Invoke("openwhisk-app/hello")
    assert.Equal(t, nil, err, "Failed to invoke openwhisk-app/hello")
    _, err = wskcli.Invoke("openwhisk-app/hello", "--param", "name", "Amy")
    assert.Equal(t, nil, err, "Failed to invoke openwhisk-app/hello")
    _, err = wskcli.Invoke("openwhisk-app/hello", "--param", "name", "Amy", "--param", "place", "California")
    assert.Equal(t, nil, err, "Failed to invoke openwhisk-app/hello")
    _, err = wskcli.Invoke("openwhisk-app/helloworld", "--param", "name", "Bob", "--param", "place", "New York")
    assert.Equal(t, nil, err, "Failed to invoke openwhisk-app/helloworld")

    _, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
    assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files.")

}

var (
    manifestPath   = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/dependency/manifest.yaml"
    deploymentPath = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/dependency/deployment.yaml"
)
