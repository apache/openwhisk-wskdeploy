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

package integration

import (
    "testing"
    "github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/triggerrule"
    "github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/zipaction"
    "github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/jaraction"
    "github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/flagstests"
)

func TestJarAction(t *testing.T) {
    jaraction.RunTestJarAction(t)
}

func TestTriggerRule(t *testing.T) {
    triggerrule.RunTestTriggerRule(t)
}

func TestZipAction(t *testing.T) {
    zipaction.RunTestZipAction(t)
}

func TestFlags(t *testing.T) {
    flagstests.RunTestSupportProjectPath(t)
    flagstests.RunTestSupportProjectPathTrailingSlash(t)
    flagstests.RunTestSupportManifestYamlPath(t)
    flagstests.RunTestSupportManifestYmlPath(t)
    flagstests.RunTestSupportManifestYamlDeployment(t)
    flagstests.RunTestSupportManifestYmlDeployment(t)
}
