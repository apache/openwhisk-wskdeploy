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
	"fmt"
	"github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

/* *
 * Please configure following env. variables in order to run this integration test:
 * 	BLUEMIX_APIHOST
 *	BLUEMIX_NAMESPACE
 * 	BLUEMIX_AUTH
 *	MESSAGEHUB_USER
 *	MESSAGEHUB_PASSWORD
 */
func TestAlarmTrigger(t *testing.T) {
	os.Setenv("MESSAGEHUB_ADMIN_HOST", "https://kafka-admin-prod01.messagehub.services.us-south.bluemix.net:443")
	os.Setenv("KAFKA_BROKERS_SASL", "[\"kafka01-prod01.messagehub.services.us-south.bluemix.net:9093\", \"kafka02-prod01.messagehub.services.us-south.bluemix.net:9093\", \"kafka03-prod01.messagehub.services.us-south.bluemix.net:9093\"]")
	os.Setenv("SRC_TOPIC", "in-topic")
	os.Setenv("DEST_TOPIC", "out-topic")

	wskprops := common.GetWskpropsFromEnvVars(common.BLUEMIX_APIHOST, common.BLUEMIX_NAMESPACE, common.BLUEMIX_AUTH)
	err := common.ValidateWskprops(wskprops)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Wsk properties are not properly configured, so tests are skipped.")
	} else {
		wskdeploy := common.NewWskdeploy()
		_, err := wskdeploy.DeployWithCredentials(manifestPath, deploymentPath, wskprops)
		assert.Equal(t, nil, err, "Failed to deploy the manifest file.")
		_, err = wskdeploy.UndeployWithCredentials(manifestPath, deploymentPath, wskprops)
		assert.Equal(t, nil, err, "Failed to undeploy the manifest file.")
	}
}

var (
	manifestPath   = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/message-hub/manifest.yaml"
	deploymentPath = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/message-hub/deployment.yaml"
)
