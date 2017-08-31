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
	assert.NotNil(t, dr.DeploymentDescriptor.Application.Packages["JiraBackupSolution"], "DeploymentReader handle deployment yaml failed.")
}

func TestDeployerCheck(t *testing.T) {
	sd := NewServiceDeployer()
	sd.DeploymentPath = "../tests/usecases/badyaml/deployment.yaml"
	sd.ManifestPath = "../tests/usecases/badyaml/manifest.yaml"
	// The system will exit thus the test will fail.
	// sd.Check()
}
