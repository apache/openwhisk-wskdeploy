// TODO(749) Rewrite test to use "packages" schema

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
	//	"fmt"
	"github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/common"
	//	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	//	"testing"
)

var wskprops = common.GetWskprops()

func composeDeployFiles(count int) (manifestStr string, deploymentStr string) {
	manifestStr = `package:
  name: TestSequencesCreation
  actions:
`
	deploymentStr = `project:
  name: TestSequencesCreationApp
  packages:
    TestSequencesCreation:
      actions:
`
	sequenceStr := `  sequences:
    validate-sequence:
      actions: `

	for i := 1; i < count+1; i++ {
		manifestStr = manifestStr + "    func" + strconv.Itoa(i) + ":" +
			`
      function: actions/function.js
      runtime: nodejs:6
      inputs:
        functionID: string
        visited: string
      outputs:
        visited: string
`
		sequenceStr = sequenceStr + "func" + strconv.Itoa(i) + ","
		deploymentStr = deploymentStr + "        func" + strconv.Itoa(i) + ":" +
			`
          inputs:
            functionID: ` + strconv.Itoa(i) +
			`
            visited:
`
	}
	manifestStr = manifestStr + strings.TrimRight(sequenceStr, ",")
	return
}

func _createTmpfile(data string, filename string) (f *os.File, err error) {
	dir, _ := os.Getwd()
	tmpfile, err := ioutil.TempFile(dir, filename)
	if err != nil {
		return nil, err
	}
	_, err = tmpfile.Write([]byte(data))
	if err != nil {
		return tmpfile, err
	}
	return tmpfile, nil
}

// TODO(749) - Rewrite to work with "packages" key/schema
//func TestValidateSequenceCreation(t *testing.T) {
//	count := 10
//	wskdeploy := common.NewWskdeploy()
//	for i := 1; i < count+1; i++ {
//		maniData, deplyData := composeDeployFiles(i + 1)
//		tmpManifile, err := _createTmpfile(maniData, "sequence_test_mani_")
//		tmpDeplyfile, err := _createTmpfile(deplyData, "sequence_test_deply_")
//		if err != nil {
//			assert.Fail(t, "Failed to create temp file")
//		}
//
//		fmt.Printf("Deploying sequence %d\n:", i)
//		_, err = wskdeploy.Deploy(tmpManifile.Name(), tmpDeplyfile.Name())
//		assert.Equal(t, nil, err, "Failed to deploy sequence.")
//		_, err = wskdeploy.Undeploy(tmpManifile.Name(), tmpDeplyfile.Name())
//		assert.Equal(t, nil, err, "Failed to undeploy sequence.")
//
//		tmpManifile.Close()
//		tmpDeplyfile.Close()
//		os.Remove(tmpManifile.Name())
//		os.Remove(tmpDeplyfile.Name())
//	}
//}
