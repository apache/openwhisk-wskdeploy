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
	"archive/zip"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"crypto/tls"
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/hokaccha/go-prettyjson"
	"io/ioutil"
	"net/http"
	"path"
    "time"
)

const (
    HTTP_TIMEOUT = 30
    DEFAULT_PROJECT_PATH = "."
)
// ActionRecord is a container to keep track of
// a whisk action struct and a location filepath we use to
// map files and manifest declared actions
type ActionRecord struct {
	Action      *whisk.Action
	Packagename string
	Filepath    string
}

type TriggerRecord struct {
	Trigger     *whisk.Trigger
	Packagename string
}

type RuleRecord struct {
	Rule        *whisk.Rule
	Packagename string
}

// Utility to convert hostname to URL object
func GetURLBase(host string) (*url.URL, error) {

	urlBase := fmt.Sprintf("%s/api", host)
	url, err := url.Parse(urlBase)

	if len(url.Scheme) == 0 || len(url.Host) == 0 {
		urlBase = fmt.Sprintf("https://%s/api", host)
		url, err = url.Parse(urlBase)
	}

	return url, err
}

func GetHomeDirectory() string {
	usr, err := user.Current()
	if err != nil {
        return ""
    }

	return usr.HomeDir
}

// Potentially complex structures(such as DeploymentApplication, DeploymentPackage)
// could implement those interface which is convenient for put, get subtract in
// containers etc.
type Comparable interface {
	HashCode() uint32
	Equals() bool
}

func IsFeedAction(trigger *whisk.Trigger) (string, bool) {
	for _, annotation := range trigger.Annotations {
		if annotation.Key == "feed" {
			return annotation.Value.(string), true
		}
	}

	return "", false
}

func IsJSON(s string) (interface{}, bool) {
	var js interface{}
	if json.Unmarshal([]byte(s), &js) == nil {
		return js, true
	}
	return nil, false

}

func PrettyJSON(j interface{}) (string, error) {
	formatter := prettyjson.NewFormatter()
	bytes, err := formatter.Marshal(j)
	if err != nil {
        return "", err
    }
	return string(bytes), nil
}

// Common utilities

// Prompt for user input
func Ask(reader *bufio.Reader, question string, def string) string {
	fmt.Print(question + " (" + def + "): ")
	answer, _ := reader.ReadString('\n')
	len := len(answer)
	if len == 1 {
		return def
	}
	return answer[:len-1]
}

// Test if a string
func isValidEnvironmentVar(value string) bool {

	// A valid Env. variable should start with '$' (dollar) char.
	// AND have at least 1 additional character after it.
	if value != "" && len(value) > 1 && strings.HasPrefix(value, "$") {
		return true
	}
	return false
}

// Get the env variable value by key.
// Get the env variable if the key is start by $
func GetEnvVar(key interface{}) interface{} {
	// Assure the key itself is not nil
	if key == nil {
		return nil
	}

	if reflect.TypeOf(key).String() == "string" {
		if isValidEnvironmentVar(key.(string)) {
			// retrieve the value of the env. var. from the host system.
			envkey := strings.Split(key.(string), "$")[1]
			value := os.Getenv(envkey)
			if value == "" {
				// Issue a warning to the user (verbose) that env. var. was not found
				// (i.e., and empty string was returned).
				fmt.Println("WARNING: Missing Environment Variable " + envkey + ".")
			}
			return value
		}

		// The key was not a valid env. variable, simply return it as the value itself (of type string)
		return key.(string)
	}
	return key
}

var kindToJSON []string = []string{"", "boolean", "integer", "integer", "integer", "integer", "integer", "integer", "integer", "integer",
	"integer", "integer", "integer", "number", "number", "number", "number", "array", "", "", "", "object", "", "", "string", "", ""}

// Gets JSON type name
func GetJSONType(j interface{}) string {
	fmt.Print(reflect.TypeOf(j).Kind())
	return kindToJSON[reflect.TypeOf(j).Kind()]
}

func NewZipWritter(src, des string) *ZipWritter {
	zw := &ZipWritter{src: src, des: des}
	return zw
}

type ZipWritter struct {
	src        string
	des        string
	zipWritter *zip.Writer
}

func (zw *ZipWritter) zipFile(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !f.Mode().IsRegular() || f.Size() == 0 {
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fileName := strings.TrimPrefix(path, zw.src+"/")
	wr, err := zw.zipWritter.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(wr, file)
	if err != nil {
		return err
	}
	return nil
}

func (zw *ZipWritter) Zip() error {
	// create zip file
	zipFile, err := os.Create(zw.des)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zw.zipWritter = zip.NewWriter(zipFile)
	err = filepath.Walk(zw.src, zw.zipFile)
	if err != nil {
		return nil
	}
	err = zw.zipWritter.Close()
	if err != nil {
		return err
	}
	return nil
}

// below codes is from wsk cli with tiny adjusts.
func GetExec(artifact string, kind string, isDocker bool, mainEntry string) (*whisk.Exec, error) {
	var err error
	var code string
	var content []byte
	var exec *whisk.Exec

	ext := filepath.Ext(artifact)
	exec = new(whisk.Exec)

	if !isDocker || ext == ".zip" {
		content, err = new(ContentReader).ReadLocal(artifact)
		if err != nil {
            return nil, err
        }
		code = string(content)
		exec.Code = &code
	}

	if len(kind) > 0 {
		exec.Kind = kind
	} else if isDocker {
		exec.Kind = "blackbox"
		if ext != ".zip" {
			exec.Image = artifact
		} else {
			exec.Image = "openwhisk/dockerskeleton"
		}
	} else if ext == ".swift" {
		exec.Kind = "swift:default"
	} else if ext == ".js" {
		exec.Kind = "nodejs:default"
	} else if ext == ".py" {
		exec.Kind = "python:default"
	} else if ext == ".jar" {
		exec.Kind = "java:default"
		exec.Code = nil
	} else {
		if ext == ".zip" {
			return nil, zipKindError()
		} else {
			return nil, extensionError(ext)
		}
	}

	// Error if entry point is not specified for Java
	if len(mainEntry) != 0 {
		exec.Main = mainEntry
	} else {
		if exec.Kind == "java" {
			return nil, javaEntryError()
		}
	}

	// Base64 encode the zip file content
	if ext == ".zip" {
		code = base64.StdEncoding.EncodeToString([]byte(code))
		exec.Code = &code
	}

	return exec, nil
}

func zipKindError() error {
	errMsg := wski18n.T("creating an action from a .zip artifact requires specifying the action kind explicitly")

	return errors.New(errMsg)
}

func extensionError(extension string) error {
	errMsg := wski18n.T(
		"'{{.name}}' is not a supported action runtime",
		map[string]interface{}{
			"name": extension,
		})

	return errors.New(errMsg)
}

func javaEntryError() error {
	errMsg := wski18n.T("Java actions require --main to specify the fully-qualified name of the main class")

	return errors.New(errMsg)
}

//for web action support, code from wsk cli with tiny adjustments
const WEB_EXPORT_ANNOT = "web-export"
const RAW_HTTP_ANNOT = "raw-http"
const FINAL_ANNOT = "final"

func WebAction(webMode string, annotations whisk.KeyValueArr, entityName string, fetch bool) (whisk.KeyValueArr, error) {
	switch strings.ToLower(webMode) {
	case "yes":
		fallthrough
	case "true":
		return webActionAnnotations(fetch, annotations, entityName, addWebAnnotations)
	case "no":
		fallthrough
	case "false":
		return webActionAnnotations(fetch, annotations, entityName, deleteWebAnnotations)
	case "raw":
		return webActionAnnotations(fetch, annotations, entityName, addRawAnnotations)
	default:
		return nil, errors.New(webMode)
	}
}

type WebActionAnnotationMethod func(annotations whisk.KeyValueArr) whisk.KeyValueArr

func webActionAnnotations(
	fetchAnnotations bool,
	annotations whisk.KeyValueArr,
	entityName string,
	webActionAnnotationMethod WebActionAnnotationMethod) (whisk.KeyValueArr, error) {
	if annotations != nil || !fetchAnnotations {
		annotations = webActionAnnotationMethod(annotations)
	}

	return annotations, nil
}

func addWebAnnotations(annotations whisk.KeyValueArr) whisk.KeyValueArr {
	annotations = deleteWebAnnotationKeys(annotations)
	annotations = addKeyValue(WEB_EXPORT_ANNOT, true, annotations)
	annotations = addKeyValue(RAW_HTTP_ANNOT, false, annotations)
	annotations = addKeyValue(FINAL_ANNOT, true, annotations)

	return annotations
}

func deleteWebAnnotations(annotations whisk.KeyValueArr) whisk.KeyValueArr {
	annotations = deleteWebAnnotationKeys(annotations)
	annotations = addKeyValue(WEB_EXPORT_ANNOT, false, annotations)
	annotations = addKeyValue(RAW_HTTP_ANNOT, false, annotations)
	annotations = addKeyValue(FINAL_ANNOT, false, annotations)

	return annotations
}

func addRawAnnotations(annotations whisk.KeyValueArr) whisk.KeyValueArr {
	annotations = deleteWebAnnotationKeys(annotations)
	annotations = addKeyValue(WEB_EXPORT_ANNOT, true, annotations)
	annotations = addKeyValue(RAW_HTTP_ANNOT, true, annotations)
	annotations = addKeyValue(FINAL_ANNOT, true, annotations)

	return annotations
}

func deleteWebAnnotationKeys(annotations whisk.KeyValueArr) whisk.KeyValueArr {
	annotations = deleteKey(WEB_EXPORT_ANNOT, annotations)
	annotations = deleteKey(RAW_HTTP_ANNOT, annotations)
	annotations = deleteKey(FINAL_ANNOT, annotations)

	return annotations
}

func deleteKey(key string, keyValueArr whisk.KeyValueArr) whisk.KeyValueArr {
	for i := 0; i < len(keyValueArr); i++ {
		if keyValueArr[i].Key == key {
			keyValueArr = append(keyValueArr[:i], keyValueArr[i+1:]...)
			break
		}
	}

	return keyValueArr
}

func addKeyValue(key string, value interface{}, keyValueArr whisk.KeyValueArr) whisk.KeyValueArr {
	keyValue := whisk.KeyValue{
		Key:   key,
		Value: value,
	}

	return append(keyValueArr, keyValue)
}

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

	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

    var netTransport = &http.Transport{
        TLSClientConfig: tlsConfig,
    }

    var netClient = &http.Client{
        Timeout: time.Second * HTTP_TIMEOUT,
        Transport: netTransport,
    }

	res, err := netClient.Do(req)
    if err != nil {
        errString := wski18n.T("Failed to get the supported runtimes from OpenWhisk service: {{.err}}.\n",
            map[string]interface{}{"err": err.Error()})
        whisk.Debug(whisk.DbgWarn, errString)
    }

    if res != nil {
        defer res.Body.Close()
    }

	// Local openwhisk deployment sometimes only returns "application/json" as the content type
	if err != nil || !strings.Contains(ct, res.Header.Get("Content-Type")) {
        stdout := wski18n.T("Start to unmarshal Openwhisk info from local values.\n")
        whisk.Debug(whisk.DbgInfo, stdout)
		err = json.Unmarshal(runtimeInfo, &op)
	} else {
		b, _ := ioutil.ReadAll(res.Body)
		if b != nil && len(b) > 0 {
            stdout := wski18n.T("Unmarshal Openwhisk info from internet.\n")
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
			if (!v[i].Deprecated) {
				rt[k] = append(rt[k], v[i].Kind)
			}
		}
	}
	return
}

var runtimeInfo = []byte(`{
  "support": {
    "github": "https://github.com/apache/incubator-openwhisk/issues",
    "slack": "http://slack.openwhisk.org"
  },
  "description": "OpenWhisk",
  "api_paths": ["/api/v1"],
  "runtimes": {
    "nodejs": [{
      "image": "openwhisk/nodejsaction:latest",
      "deprecated": true,
      "requireMain": false,
      "default": false,
      "attached": false,
      "kind": "nodejs"
    }, {
      "image": "openwhisk/nodejs6action:latest",
      "deprecated": false,
      "requireMain": false,
      "default": true,
      "attached": false,
      "kind": "nodejs:6"
    }],
    "java": [{
      "image": "openwhisk/java8action:latest",
      "deprecated": false,
      "requireMain": true,
      "default": true,
      "attached": true,
      "kind": "java"
    }],
    "php": [{
      "image": "openwhisk/action-php-v7.1:latest",
      "deprecated": false,
      "requireMain": false,
      "default": true,
      "attached": false,
      "kind": "php:7.1"
    }],
    "python": [{
      "image": "openwhisk/python2action:latest",
      "deprecated": false,
      "requireMain": false,
      "default": false,
      "attached": false,
      "kind": "python"
    }, {
      "image": "openwhisk/python2action:latest",
      "deprecated": false,
      "requireMain": false,
      "default": true,
      "attached": false,
      "kind": "python:2"
    }, {
      "image": "openwhisk/python3action:latest",
      "deprecated": false,
      "requireMain": false,
      "default": false,
      "attached": false,
      "kind": "python:3"
    }],
    "swift": [{
      "image": "openwhisk/swiftaction:latest",
      "deprecated": true,
      "requireMain": false,
      "default": false,
      "attached": false,
      "kind": "swift"
    }, {
      "image": "openwhisk/swift3action:latest",
      "deprecated": false,
      "requireMain": false,
      "default": true,
      "attached": false,
      "kind": "swift:3"
    }, {
      "image": "openwhisk/action-swift-v3.1.1:latest",
      "deprecated": false,
      "requireMain": false,
      "default": false,
      "attached": false,
      "kind": "swift:3.1.1"
    }]
  },
  "limits": {
    "actions_per_minute": 5000,
    "triggers_per_minute": 5000,
    "concurrent_actions": 1000
  }
  }
`)

var Rts map[string][]string

var DefaultRts = map[string][]string{
	"nodejs": {"nodejs", "nodejs:6"},
	"java":   {"java"},
	"php":    {"php:7.1"},
	"python": {"python", "python:2", "python:3"},
	"swift":  {"swift", "swift:3", "swift:3.1.1"},
}

func CheckExistRuntime(rtname string, rts map[string][]string) bool {
	for _, v := range rts {
		for i := range v {
			if rtname == v[i] {
				return true
			}
		}
	}
	return false
}

func GetManifestFilePath(projectPath string) string {
	if _, err := os.Stat(path.Join(projectPath, ManifestFileNameYaml)); err == nil {
		return path.Join(projectPath, ManifestFileNameYaml)
	} else if _, err := os.Stat(path.Join(projectPath, ManifestFileNameYml)); err == nil {
		return path.Join(projectPath, ManifestFileNameYml)
	} else {
		return ""
	}
}

func GetDeploymentFilePath(projectPath string) string {
	if _, err := os.Stat(path.Join(projectPath, DeploymentFileNameYaml)); err == nil {
		return path.Join(projectPath, DeploymentFileNameYaml)
	} else if _, err := os.Stat(path.Join(projectPath, DeploymentFileNameYml)); err == nil {
		return path.Join(projectPath, DeploymentFileNameYml)
	} else {
		return ""
	}
}

// name of manifest and deployment files
const ManifestFileNameYaml = "manifest.yaml"
const ManifestFileNameYml = "manifest.yml"
const DeploymentFileNameYaml = "deployment.yaml"
const DeploymentFileNameYml = "deployment.yml"
