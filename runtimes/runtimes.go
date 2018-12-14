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
	"net/url"
	"strings"
	"time"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
)

const (
	NODEJS_FILE_EXTENSION   = "js"
	SWIFT_FILE_EXTENSION    = "swift"
	PYTHON_FILE_EXTENSION   = "py"
	JAVA_FILE_EXTENSION     = "java"
	JAR_FILE_EXTENSION      = "jar"
	CSHARP_FILE_EXTENSION   = "cs"
	PHP_FILE_EXTENSION      = "php"
	ZIP_FILE_EXTENSION      = "zip"
	RUBY_FILE_EXTENSION     = "rb"
	GO_FILE_EXTENSION       = "go"
	NODEJS_RUNTIME          = "nodejs"
	SWIFT_RUNTIME           = SWIFT_FILE_EXTENSION
	PYTHON_RUNTIME          = "python"
	JAVA_RUNTIME            = JAVA_FILE_EXTENSION
	DOTNET_RUNTIME          = ZIP_FILE_EXTENSION
	PHP_RUNTIME             = PHP_FILE_EXTENSION
	RUBY_RUNTIME            = "ruby"
	GO_RUNTIME              = GO_FILE_EXTENSION
	HTTP_CONTENT_TYPE_KEY   = "Content-Type"
	HTTP_CONTENT_TYPE_VALUE = "application/json; charset=UTF-8"
	RUNTIME_NOT_SPECIFIED   = "NOT SPECIFIED"
	BLACKBOX                = "blackbox"
	HTTPS                   = "https://"
)

// Structs used to denote the OpenWhisk Runtime information
type Limit struct {
	Apm       uint `json:"actions_per_minute"`
	Tpm       uint `json:"triggers_per_minute"`
	ConAction uint `json:"concurrent_actions"`
}

type Runtime struct {
	Deprecated bool   `json:"deprecated"`
	Default    bool   `json:"default"`
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
	opURL := apiHost
	_, err = url.ParseRequestURI(opURL)
	if err != nil {
		opURL = HTTPS + opURL
	}
	req, _ := http.NewRequest("GET", opURL, nil)
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
		err = json.Unmarshal(RUNTIME_DETAILS, &op)
		if err != nil {
			errMessage := wski18n.T(wski18n.ID_ERR_RUNTIME_PARSER_ERROR,
				map[string]interface{}{wski18n.KEY_ERR: err.Error()})
			err = wskderrors.NewRuntimeParserError(errMessage)
		}
	} else {
		b, _ := ioutil.ReadAll(res.Body)
		if b != nil && len(b) > 0 {
			stdout := wski18n.T(wski18n.ID_MSG_UNMARSHAL_NETWORK_X_url_X,
				map[string]interface{}{"url": opURL})
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
		if strings.Contains(k, NODEJS_RUNTIME) {
			ext[NODEJS_FILE_EXTENSION] = k
		} else if strings.Contains(k, PYTHON_RUNTIME) {
			ext[PYTHON_FILE_EXTENSION] = k
		} else if strings.Contains(k, SWIFT_RUNTIME) {
			ext[SWIFT_FILE_EXTENSION] = k
		} else if strings.Contains(k, PHP_RUNTIME) {
			ext[PHP_FILE_EXTENSION] = k
		} else if strings.Contains(k, JAVA_RUNTIME) {
			ext[JAVA_FILE_EXTENSION] = k
			ext[JAR_FILE_EXTENSION] = k
		} else if strings.Contains(k, RUBY_RUNTIME) {
			ext[RUBY_FILE_EXTENSION] = k
		} else if strings.Contains(k, GO_RUNTIME) {
			ext[GO_FILE_EXTENSION] = k
        } else if strings.Contains(k, DOTNET_RUNTIME) {
			ext[CSHARP_FILE_EXTENSION] = k
			ext[ZIP_FILE_EXTENSION] = k
        }
	}
	return
}

func FileRuntimeExtensions(op OpenWhiskInfo) (rte map[string]string) {
	rte = make(map[string]string)

	for k, v := range op.Runtimes {
		for i := range v {
			if !v[i].Deprecated {
				if strings.Contains(k, NODEJS_RUNTIME) {
					rte[v[i].Kind] = NODEJS_FILE_EXTENSION
				} else if strings.Contains(k, PYTHON_RUNTIME) {
					rte[v[i].Kind] = PYTHON_FILE_EXTENSION
				} else if strings.Contains(k, SWIFT_RUNTIME) {
					rte[v[i].Kind] = SWIFT_FILE_EXTENSION
				} else if strings.Contains(k, PHP_RUNTIME) {
					rte[v[i].Kind] = PHP_FILE_EXTENSION
				} else if strings.Contains(k, JAVA_RUNTIME) {
					rte[v[i].Kind] = JAVA_FILE_EXTENSION
				} else if strings.Contains(k, RUBY_RUNTIME) {
					rte[v[i].Kind] = RUBY_FILE_EXTENSION
				} else if strings.Contains(k, GO_RUNTIME) {
					rte[v[i].Kind] = GO_FILE_EXTENSION
				} else if strings.Contains(k, DOTNET_RUNTIME) {
                    rte[v[i].Kind] = CSHARP_FILE_EXTENSION
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

var RUNTIME_DETAILS = []byte(`{
    "support": {
        "github": "https://github.com/apache/incubator-openwhisk/issues",
        "slack": "http://slack.openwhisk.org"
    },
    "description": "OpenWhisk",
    "api_paths": [
        "/api/v1"
    ],
    "runtimes": {
        "nodejs": [
            {
                "kind": "nodejs",
                "image": {
                    "prefix": "openwhisk",
                    "name": "nodejsaction",
                    "tag": "latest"
                },
                "deprecated": true,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            },
            {
                "kind": "nodejs:6",
                "default": true,
                "image": {
                    "prefix": "openwhisk",
                    "name": "nodejs6action",
                    "tag": "latest"
                },
                "deprecated": false,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                },
                "stemCells": [{
                    "count": 2,
                    "memory": "256 MB"
                }]
            },
            {
                "kind": "nodejs:8",
                "default": false,
                "image": {
                    "prefix": "openwhisk",
                    "name": "action-nodejs-v8",
                    "tag": "latest"
                },
                "deprecated": false,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            },
            {
                "kind": "nodejs:10",
                "default": false,
                "image": {
                    "prefix": "openwhisk",
                    "name": "action-nodejs-v10",
                    "tag": "latest"
                },
                "deprecated": false,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            }
        ],
        "python": [
            {
                "kind": "python",
                "image": {
                    "prefix": "openwhisk",
                    "name": "python2action",
                    "tag": "latest"
                },
                "deprecated": false,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            },
            {
                "kind": "python:2",
                "default": true,
                "image": {
                    "prefix": "openwhisk",
                    "name": "python2action",
                    "tag": "latest"
                },
                "deprecated": false,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            },
            {
                "kind": "python:3",
                "image": {
                    "prefix": "openwhisk",
                    "name": "python3action",
                    "tag": "latest"
                },
                "deprecated": false,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            }
        ],
        "swift": [
            {
                "kind": "swift",
                "image": {
                    "prefix": "openwhisk",
                    "name": "swiftaction",
                    "tag": "latest"
                },
                "deprecated": true,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            },
            {
                "kind": "swift:3",
                "image": {
                    "prefix": "openwhisk",
                    "name": "swift3action",
                    "tag": "latest"
                },
                "deprecated": true,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            },
            {
                "kind": "swift:3.1.1",
                "image": {
                    "prefix": "openwhisk",
                    "name": "action-swift-v3.1.1",
                    "tag": "latest"
                },
                "deprecated": false,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            },
            {
                "kind": "swift:4.1",
                "default": true,
                "image": {
                    "prefix": "openwhisk",
                    "name": "action-swift-v4.1",
                    "tag": "latest"
                },
                "deprecated": false,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            }
        ],
        "java": [
            {
                "kind": "java",
                "default": true,
                "image": {
                    "prefix": "openwhisk",
                    "name": "java8action",
                    "tag": "latest"
                },
                "deprecated": false,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                },
                "requireMain": true
            }
        ],
        "php": [
            {
                "kind": "php:7.1",
                "default": false,
                "deprecated": false,
                "image": {
                    "prefix": "openwhisk",
                    "name": "action-php-v7.1",
                    "tag": "latest"
                },
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            },
            {
                "kind": "php:7.2",
                "default": true,
                "deprecated": false,
                "image": {
                    "prefix": "openwhisk",
                    "name": "action-php-v7.2",
                    "tag": "latest"
                },
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            }
        ],
        "ruby": [
            {
                "kind": "ruby:2.5",
                "default": true,
                "deprecated": false,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                },
                "image": {
                    "prefix": "openwhisk",
                    "name": "action-ruby-v2.5",
                    "tag": "latest"
                }
            }
        ],
        "go": [
            {
                "kind": "go:1.11",
                "default": true,
                "deprecated": false,
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                },
                "image": {
                    "prefix": "openwhisk",
                    "name": "actionloop-golang-v1.11",
                    "tag": "latest"
                }
            }
        ],
        "dotnet": [
            {
                "kind": "dotnet:2.2",
                "default": true,
                "deprecated": false,
                "requireMain": true,
                "image": {
                    "prefix": "openwhisk",
                    "name": "action-dotnet-v2.2",
                    "tag": "latest"
                },
                "attached": {
                    "attachmentName": "codefile",
                    "attachmentType": "text/plain"
                }
            }
        ]
    },
    "blackboxes": [
        {
            "prefix": "openwhisk",
            "name": "dockerskeleton",
            "tag": "latest"
        }
    ]
}`)
