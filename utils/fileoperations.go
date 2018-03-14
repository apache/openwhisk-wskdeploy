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
	"fmt"
	"os"
	"strings"
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
		return false
	} else {
		return true
	}
}

func IsDirectory(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	} else {
		defer f.Close()
	}

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	default:
		return false
	}
}

func ReadProps(path string) (map[string]string, error) {

	props := map[string]string{}

	file, err := os.Open(path)
	if err != nil {
		// If file does not exist, just return props
		fmt.Printf("Warning: Unable to read whisk properties file '%s' (file open error: %s)\n", path, err)
		return props, err
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
		if err != nil {
			return err
		}

		_, err = file.WriteString("=")
		if err != nil {
			return err
		}

		_, err = file.WriteString(v)
		if err != nil {
			return err
		}

		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteFile(filename string, content string) error {
	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	if _, err = file.WriteString(content); err != nil {
		return err
	}

	return nil
}
