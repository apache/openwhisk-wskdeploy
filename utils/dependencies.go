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
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
)

const (
	GITHUB       = "github"
	WHISK_SYSTEM = "whisk.system"
)

type DependencyRecord struct {
	ProjectPath string //root of the source codes of dependent projects, e.g. src_project_path/Packages
	Packagename string //name of the package
	Location    string //location
	Version     string //version
	Parameters  whisk.KeyValueArr
	Annotations whisk.KeyValueArr
	IsBinding   bool
	BaseRepo    string
	SubFolder   string
}

func NewDependencyRecord(projectPath string,
	packagename string,
	location string,
	version string,
	parameters whisk.KeyValueArr,
	annotations whisk.KeyValueArr,
	isBinding bool) DependencyRecord {
	var record DependencyRecord
	record.ProjectPath = projectPath
	record.Packagename = packagename
	record.Location = location
	record.Version = version
	record.Parameters = parameters
	record.Annotations = annotations
	record.IsBinding = isBinding
	//split url to BaseUrl and SubFolder
	if !record.IsBinding {
		paths := strings.Split(location, "/")
		record.BaseRepo = strings.Join([]string{paths[0], paths[1], paths[2], paths[3], paths[4]}, "/")
		if len(paths) > 5 {
			record.SubFolder = strings.TrimPrefix(record.Location, record.BaseRepo)
		} else {
			record.SubFolder = ""
		}
	}

	return record
}

func LocationIsBinding(location string) bool {
	if strings.HasPrefix(location, "/"+WHISK_SYSTEM) || strings.HasPrefix(location, WHISK_SYSTEM) || strings.HasPrefix(location, "/") {
		return true
	}

	return false
}

func removeProtocol(location string) string {
	if paths := strings.SplitAfterN(location, "://", 2); len(paths) == 2 {
		return paths[1]
	}
	return location
}

func LocationIsGithub(location string) bool {
	paths := strings.SplitN(removeProtocol(location), "/", 2)
	return strings.Contains(paths[0], GITHUB)
}
