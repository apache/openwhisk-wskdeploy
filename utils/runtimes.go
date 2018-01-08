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

package utils

import (
	"crypto/tls"
	"encoding/json"
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
)

const NODEJS_FILE_EXTENSION = "js"
const SWIFT_FILE_EXTENSION = "swift"
const PYTHON_FILE_EXTENSION = "py"
const JAVA_FILE_EXTENSION = "java"
const JAR_FILE_EXTENSION = "jar"
const PHP_FILE_EXTENSION = "php"
const ZIP_FILE_EXTENSION = "zip"

// Structs used to denote the OpenWhisk Runtime information
type Limit struct {
	Apm       uint16 `json:"actions_per_minute"`
	Tpm       uint16 `json:"triggers_per_minute"`
	ConAction uint16 `json:"concurrent_actions"`
}

type Runtime struct {
	Image      string `json:"image"`
	Deprecated bool   `json:"deprecated"`
	ReMain     bool   `json:"requireMain"`
	Default    bool   `json:"default"`
	Attach     bool   `json:"attached"`
	Kind       string `json:"kind"`
}

type SupportInfo struct {
	Github string `json:"github"`
	Slack  string `json:"slack"`
}

type OpenWhiskInfo struct {
	Support  SupportInfo          `json:"support"`
	Desc     string               `json:"description"`
	ApiPath  []string             `json:"api_paths"`
	Runtimes map[string][]Runtime `json:"runtimes"`
	Limits   Limit                `json:"limits"`
}

var FileExtensionRuntimeKindMap map[string]string
var SupportedRunTimes map[string][]string
var DefaultRunTimes map[string]string


// We could get the openwhisk info from bluemix through running the command
// `curl -k https://openwhisk.ng.bluemix.net`
// hard coding it here in case of network unavailable or failure.
func ParseOpenWhisk(apiHost string) (op OpenWhiskInfo, err error) {
	ct := "application/json; charset=UTF-8"
	req, _ := http.NewRequest("GET", "https://"+apiHost, nil)
	req.Header.Set("Content-Type", ct)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	var netTransport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	var netClient = &http.Client{
		Timeout:   time.Second * DEFAULT_HTTP_TIMEOUT,
		Transport: netTransport,
	}

	res, err := netClient.Do(req)
	if err != nil {
		// TODD() create an error
		errString := wski18n.T(wski18n.ID_ERR_GET_RUNTIMES_X_err_X,
			map[string]interface{}{"err": err.Error()})
		whisk.Debug(whisk.DbgWarn, errString)
	}

	if res != nil {
		defer res.Body.Close()
	}

	// Local openwhisk deployment sometimes only returns "application/json" as the content type
	if err != nil || !strings.Contains(ct, res.Header.Get("Content-Type")) {
		stdout := wski18n.T(wski18n.ID_MSG_UNMARSHAL_LOCAL)
		whisk.Debug(whisk.DbgInfo, stdout)
		err = json.Unmarshal(RUNTIME_DETAILS, &op)
	} else {
		b, _ := ioutil.ReadAll(res.Body)
		if b != nil && len(b) > 0 {
			stdout := wski18n.T(wski18n.ID_MSG_UNMARSHAL_NETWORK)
			whisk.Debug(whisk.DbgInfo, stdout)
			err = json.Unmarshal(b, &op)
		}
	}
	return
}

func ConvertToMap(op OpenWhiskInfo) (rt map[string][]string) {
	rt = make(map[string][]string)
	for k, v := range op.Runtimes {
		rt[k] = make([]string, 0, len(v))
		for i := range v {
			if !v[i].Deprecated {
				rt[k] = append(rt[k], v[i].Kind)
			}
		}
	}
	return
}

func DefaultRuntimes(op OpenWhiskInfo) (rt map[string]string) {
	rt = make(map[string]string)
	for k, v := range op.Runtimes {
		for i := range v {
			if !v[i].Deprecated {
				if v[i].Default {
					rt[k] = v[i].Kind
				}
			}
		}
	}
	return
}

func FileExtensionRuntimes(op OpenWhiskInfo) (ext map[string]string) {
	ext = make(map[string]string)
	for k := range op.Runtimes {
		if strings.Contains(k, NODEJS_FILE_EXTENSION) {
			ext[NODEJS_FILE_EXTENSION] = k
		} else if strings.Contains(k, PYTHON_FILE_EXTENSION) {
			ext[PYTHON_FILE_EXTENSION] = k
		} else if strings.Contains(k, SWIFT_FILE_EXTENSION) {
			ext[SWIFT_FILE_EXTENSION] = k
		} else if strings.Contains(k, PHP_FILE_EXTENSION) {
			ext[PHP_FILE_EXTENSION] = k
		} else if strings.Contains(k, JAVA_FILE_EXTENSION) {
			ext[JAVA_FILE_EXTENSION] = k
			ext[JAR_FILE_EXTENSION] = k
		}
	}
	return
}

func CheckRuntimeConsistencyWithFileExtension (ext string, runtime string) bool {
	rt := FileExtensionRuntimeKindMap[ext]
	for _, v := range SupportedRunTimes[rt] {
		if (runtime == v) {
			return true
		}
	}
	return false
}

func CheckExistRuntime(rtname string, runtimes map[string][]string) bool {
	for _, v := range runtimes {
		for i := range v {
			if rtname == v[i] {
				return true
			}
		}
	}
	return false
}

func ListOfSupportedRuntimes(runtimes map[string][]string) (rt []string) {
	for _, v := range runtimes {
		for _, t := range v {
			rt = append(rt, t)
		}
	}
	return
}

var RUNTIME_DETAILS = []byte(`{
	"support":{
		"github":"https://github.com/apache/incubator-openwhisk/issues",
		"slack":"http://slack.openwhisk.org"
	},
	"description":"OpenWhisk",
	"api_paths":["/api/v1"],
	"runtimes":{
		"nodejs":[{
			"image":"openwhisk/nodejsaction:latest",
			"deprecated":true,
			"requireMain":false,
			"default":false,
			"attached":false,
			"kind":"nodejs"
		},{
			"image":"openwhisk/nodejs6action:latest",
			"deprecated":false,
			"requireMain":false,
			"default":true,
			"attached":false,
			"kind":"nodejs:6"
		},{
			"image":"openwhisk/action-nodejs-v8:latest",
			"deprecated":false,
			"requireMain":false,
			"default":false,
			"attached":false,
			"kind":"nodejs:8"
		}],
		"java":[{
			"image":"openwhisk/java8action:latest",
			"deprecated":false,
			"requireMain":true,
			"default":true,
			"attached":true,
			"kind":"java"
		}],
		"php":[{
			"image":"openwhisk/action-php-v7.1:latest",
			"deprecated":false,
			"requireMain":false,
			"default":true,
			"attached":false,
			"kind":"php:7.1"
		}],
		"python":[{
			"image":"openwhisk/python2action:latest",
			"deprecated":false,
			"requireMain":false,
			"default":false,
			"attached":false,
			"kind":"python"
		},{
			"image":"openwhisk/python2action:latest",
			"deprecated":false,
			"requireMain":false,
			"default":true,
			"attached":false,
			"kind":"python:2"
		},{
			"image":"openwhisk/python3action:latest",
			"deprecated":false,
			"requireMain":false,
			"default":false,
			"attached":false,
			"kind":"python:3"
		}],
		"swift":[{
			"image":"openwhisk/swiftaction:latest",
			"deprecated":true,
			"requireMain":false,
			"default":false,
			"attached":false,
			"kind":"swift"
		},{
			"image":"openwhisk/swift3action:latest",
			"deprecated":true,
			"requireMain":false,
			"default":false,
			"attached":false,
			"kind":"swift:3"
		},{
			"image":"openwhisk/action-swift-v3.1.1:latest",
			"deprecated":false,
			"requireMain":false,
			"default":true,
			"attached":false,
			"kind":"swift:3.1.1"
		}]
	},
	"limits":{
		"actions_per_minute":5000,
		"triggers_per_minute":5000,
		"concurrent_actions":1000
	}
	}
`)
