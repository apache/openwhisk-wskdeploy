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
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
)

func MayExists(file string) bool {
	if strings.HasPrefix(file, "http") {
		return true // to avoid multiple fetches
	} else {
		return FileExists(file)
	}
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		err = errors.New("File not found.")
		return false
	} else {
		return true
	}
}

func IsDirectory(filePath string) bool {
	f, err := os.Open(filePath)
	Check(err)

	defer f.Close()

	fi, err := f.Stat()
	Check(err)

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	default:
		return false
	}
}

func CreateActionFromFile(manipath, filePath string) (*whisk.Action, error) {
	ext := path.Ext(filePath)
	baseName := path.Base(filePath)
	//check if the file if from local or from web
	//currently only consider http
	islocal := !strings.HasPrefix(filePath, "http")
	name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	action := new(whisk.Action)
	//better refactor this
	if islocal {
		splitmanipath := strings.Split(manipath, string(os.PathSeparator))
		filePath = strings.TrimRight(manipath, splitmanipath[len(splitmanipath)-1]) + filePath
	}
	// process source code files
	if ext == ".swift" || ext == ".js" || ext == ".py" {

		kind := "nodejs:default"

		switch ext {
		case ".swift":
			kind = "swift:default"
		case ".js":
			kind = "nodejs:default"
		case ".py":
			kind = "python"
		}

		var dat []byte
		var err error

		if islocal {
			dat, err = new(ContentReader).LocalReader.ReadLocal(filePath)
		} else {
			dat, err = new(ContentReader).URLReader.ReadUrl(filePath)
		}

		code := string(dat)
		pub := false
		Check(err)
		action.Exec = new(whisk.Exec)
		action.Exec.Code = &code
		action.Exec.Kind = kind
		action.Name = name
		action.Publish = &pub
		return action, nil
		//dat, err := new(ContentReader).URLReader.ReadUrl(filePath)
		//Check(err)

	}
	// If the action is not supported, we better to return an error.
	return nil, errors.New("Unsupported action type.")
}

func ReadProps(path string) (map[string]string, error) {

	props := map[string]string{}

	file, err := os.Open(path)
	if err != nil {
		// If file does not exist, just return props
		fmt.Printf("Warning: Unable to read whisk properties file '%s' (file open error: %s)\n", path, err)
		return props, nil
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	props = map[string]string{}
	for _, line := range lines {
		kv := strings.Split(line, "=")
		if len(kv) != 2 {
			// Invalid format; skip
			continue
		}
		props[kv[0]] = kv[1]
	}

	return props, nil

}

func WriteProps(path string, props map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Warning: Unable to create whisk properties file '%s' (file create error: %s)\n", path, err)
		return err
	}
	defer file.Close()

	for k, v := range props {
		_, err := file.WriteString(k)
		Check(err)

		_, err = file.WriteString("=")
		Check(err)

		_, err = file.WriteString(v)
		Check(err)

		_, err = file.WriteString("\n")
		Check(err)
	}
	return nil
}

// Load configuration will load properties from a file
func LoadConfiguration(propPath string) ([]string, error) {
	props, err := ReadProps(propPath)
	Check(err)
	Namespace := props["NAMESPACE"]
	Apihost := props["APIHOST"]
	Authtoken := props["AUTH"]
	return []string{Namespace, Apihost, Authtoken}, nil
}
