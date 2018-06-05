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

package runtimes

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
	"path/filepath"
	"runtime"
)

const (
	NODEJS_FILE_EXTENSION   = "js"
	SWIFT_FILE_EXTENSION    = "swift"
	PYTHON_FILE_EXTENSION   = "py"
	JAVA_FILE_EXTENSION     = "java"
	JAR_FILE_EXTENSION      = "jar"
	PHP_FILE_EXTENSION      = "php"
	ZIP_FILE_EXTENSION      = "zip"
	HTTP_CONTENT_TYPE_KEY   = "Content-Type"
	HTTP_CONTENT_TYPE_VALUE = "application/json; charset=UTF-8"
	RUNTIME_NOT_SPECIFIED   = "NOT SPECIFIED"
	BLACKBOX                = "blackbox"
	RUNTIMES_FILE_NAME      = "runtimes.json"
	HTTPS                   = "https://"
)

// Structs used to denote the OpenWhisk Runtime information
type Limit struct {
	Apm       uint `json:"actions_per_minute"`
	Tpm       uint `json:"triggers_per_minute"`
	ConAction uint `json:"concurrent_actions"`
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
var FileRuntimeExtensionsMap map[string]string

// We could get the openwhisk info from bluemix through running the command
// `curl -k https://openwhisk.ng.bluemix.net`
// hard coding it here in case of network unavailable or failure.
func ParseOpenWhisk(apiHost string) (op OpenWhiskInfo, err error) {
	url := HTTPS + apiHost
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set(HTTP_CONTENT_TYPE_KEY, HTTP_CONTENT_TYPE_VALUE)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	var netTransport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	var netClient = &http.Client{
		Timeout:   time.Second * utils.DEFAULT_HTTP_TIMEOUT,
		Transport: netTransport,
	}

	res, err := netClient.Do(req)
	if err != nil {
		// TODO() create an error
		errString := wski18n.T(wski18n.ID_ERR_RUNTIMES_GET_X_err_X,
			map[string]interface{}{"err": err.Error()})
		whisk.Debug(whisk.DbgWarn, errString)
		if utils.Flags.Strict {
			errMessage := wski18n.T(wski18n.ID_ERR_RUNTIME_PARSER_ERROR,
				map[string]interface{}{wski18n.KEY_ERR: err.Error()})
			err = wskderrors.NewRuntimeParserError(errMessage)
			return
		}
	}

	if res != nil {
		defer res.Body.Close()
	}

	// Local openwhisk deployment sometimes only returns "application/json" as the content type
	if err != nil || !strings.Contains(HTTP_CONTENT_TYPE_VALUE, res.Header.Get(HTTP_CONTENT_TYPE_KEY)) {
		stdout := wski18n.T(wski18n.ID_MSG_UNMARSHAL_LOCAL)
		wskprint.PrintOpenWhiskInfo(stdout)
		runtimeDetails := readRuntimes()
		if runtimeDetails != nil {
			err = json.Unmarshal(runtimeDetails, &op)
			if err != nil {
				errMessage := wski18n.T(wski18n.ID_ERR_RUNTIME_PARSER_ERROR,
					map[string]interface{}{wski18n.KEY_ERR: err.Error()})
				err = wskderrors.NewRuntimeParserError(errMessage)
			}
		}
	} else {
		b, _ := ioutil.ReadAll(res.Body)
		if b != nil && len(b) > 0 {
			stdout := wski18n.T(wski18n.ID_MSG_UNMARSHAL_NETWORK_X_url_X,
				map[string]interface{}{"url": url})
			wskprint.PrintOpenWhiskInfo(stdout)
			err = json.Unmarshal(b, &op)
			if err != nil {
				errMessage := wski18n.T(wski18n.ID_ERR_RUNTIME_PARSER_ERROR,
					map[string]interface{}{wski18n.KEY_ERR: err.Error()})
				err = wskderrors.NewRuntimeParserError(errMessage)
			}
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

func FileRuntimeExtensions(op OpenWhiskInfo) (rte map[string]string) {
	rte = make(map[string]string)

	for k, v := range op.Runtimes {
		for i := range v {
			if !v[i].Deprecated {
				if strings.Contains(k, NODEJS_FILE_EXTENSION) {
					rte[v[i].Kind] = NODEJS_FILE_EXTENSION
				} else if strings.Contains(k, PYTHON_FILE_EXTENSION) {
					rte[v[i].Kind] = PYTHON_FILE_EXTENSION
				} else if strings.Contains(k, SWIFT_FILE_EXTENSION) {
					rte[v[i].Kind] = SWIFT_FILE_EXTENSION
				} else if strings.Contains(k, PHP_FILE_EXTENSION) {
					rte[v[i].Kind] = PHP_FILE_EXTENSION
				} else if strings.Contains(k, JAVA_FILE_EXTENSION) {
					rte[v[i].Kind] = JAVA_FILE_EXTENSION
				}
			}
		}
	}
	return
}

func CheckRuntimeConsistencyWithFileExtension(ext string, runtime string) bool {
	rt := FileExtensionRuntimeKindMap[ext]
	for _, v := range SupportedRunTimes[rt] {
		if runtime == v {
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

func readRuntimes() []byte {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	runtimesFileWithPath := filepath.Join(basepath, RUNTIMES_FILE_NAME)
	file, readErr := ioutil.ReadFile(runtimesFileWithPath)
	if readErr != nil {
		wskprint.PrintlnOpenWhiskWarning(readErr.Error())
		return nil
	}
	return file
}
