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
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type GitReader struct {
	Name string // the name of the dependency
	Url  string // pkg repo location, e.g. github.com/user/repo
	//BaseRepo    string	// base url of the git repo, e.g. github.com/user/repo
	//SubFolder   string	// subfolder of the package under BaseUrl
	Version     string
	ProjectPath string // The root folder of all dependency packages, e.g. src_project_path/Packages
}

func NewGitReader(projectName string, record DependencyRecord) *GitReader {
	var gitReader GitReader

	gitReader.Name = projectName
	gitReader.Url = record.BaseRepo
	gitReader.Version = record.Version
	gitReader.ProjectPath = record.ProjectPath

	return &gitReader

}

func (reader *GitReader) CloneDependency() error {
	zipFileName := reader.Name + "." + reader.Version + ".zip"
	zipFilePath := reader.Url + "/zipball" + "/" + reader.Version

	os.MkdirAll(reader.ProjectPath, os.ModePerm)
	output, err := os.Create(path.Join(reader.ProjectPath, zipFileName))
	if err != nil {
		return err
	}
	defer output.Close()

	response, err := http.Get(zipFilePath)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return err
	}

	zipReader, err := zip.OpenReader(path.Join(reader.ProjectPath, zipFileName))
	if err != nil {
		return err
	}

	u, err := url.Parse(reader.Url)
	team, _ := path.Split(u.Path)

	team = strings.TrimPrefix(team, "/")
	team = strings.TrimSuffix(team, "/")

	for _, file := range zipReader.File {
		path := filepath.Join(reader.ProjectPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	rootDir := filepath.Join(reader.ProjectPath, zipReader.File[0].Name)
	depPath := filepath.Join(reader.ProjectPath, reader.Name+"-"+reader.Version)
	//if the folder exists, remove it at first
	if _, err := os.Stat(depPath); err == nil {
		os.Remove(depPath)
	}
	os.Rename(rootDir, depPath)
	os.Remove(filepath.Join(reader.ProjectPath, zipFileName))

	return nil
}
