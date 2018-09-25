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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/hokaccha/go-prettyjson"
)

const (
	DEFAULT_HTTP_TIMEOUT = 60
	DEFAULT_PROJECT_PATH = "."
	HTTP_FILE_EXTENSION  = "http"
	// name of manifest and deployment files
	ManifestFileNameYaml   = "manifest.yaml"
	ManifestFileNameYml    = "manifest.yml"
	DeploymentFileNameYaml = "deployment.yaml"
	DeploymentFileNameYml  = "deployment.yml"
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

func GetHomeDirectory() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return usr.HomeDir
}

// Potentially complex structures(such as DeploymentProject, DeploymentPackage)
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

func PrettyJSON(j interface{}) (string, error) {
	formatter := prettyjson.NewFormatter()
	bytes, err := formatter.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
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

// agnostic util reader to fetch content from web or local path or potentially other places.
type ContentReader struct {
	URLReader
	LocalReader
}

type URLReader struct {
}

func (urlReader *URLReader) ReadUrl(url string) (content []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return content, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	} else {
		defer resp.Body.Close()
	}
	return b, nil
}

type LocalReader struct {
}

func (localReader *LocalReader) ReadLocal(path string) ([]byte, error) {
	cont, err := ioutil.ReadFile(path)
	return cont, err
}

func Read(url string) ([]byte, error) {
	if strings.HasPrefix(url, HTTP_FILE_EXTENSION) {
		return new(ContentReader).URLReader.ReadUrl(url)
	} else {
		return new(ContentReader).LocalReader.ReadLocal(url)
	}
}

func GetJSONFromStrings(content []string, keyValueFormat bool) (interface{}, error) {
	var data map[string]interface{}
	var res interface{}

	for i := 0; i < len(content); i++ {
		dc := json.NewDecoder(strings.NewReader(content[i]))
		dc.UseNumber()
		if err := dc.Decode(&data); err != nil {
			return whisk.KeyValueArr{}, err
		}
	}

	if keyValueFormat {
		res = getKeyValueFormattedJSON(data)
	} else {
		res = data
	}

	return res, nil
}

func getKeyValueFormattedJSON(data map[string]interface{}) whisk.KeyValueArr {
	var keyValueArr whisk.KeyValueArr
	for key, value := range data {
		keyValue := whisk.KeyValue{
			Key:   key,
			Value: value,
		}
		keyValueArr = append(keyValueArr, keyValue)
	}
	return keyValueArr
}
