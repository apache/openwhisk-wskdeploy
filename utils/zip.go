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
	"github.com/davecgh/go-spew/spew"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func NewZipWritter(src string, des string, include [][]string, manifestFilePath string) *ZipWritter {
	zw := &ZipWritter{src: src, des: des, include: include, manifestFilePath: manifestFilePath}
	return zw
}

type ZipWritter struct {
	src              string
	des              string
	include          [][]string
	manifestFilePath string
	zipWritter       *zip.Writer
}

type Include struct {
	source      string
	destination string
}

func (zw *ZipWritter) zipFile(path string, f os.FileInfo, err error) error {
	var file *os.File
	var wr io.Writer
	if err != nil {
		return err
	}
	if !f.Mode().IsRegular() || f.Size() == 0 {
		return nil
	}
	if file, err = os.Open(path); err != nil {
		return err
	}
	defer file.Close()

	fileName := strings.TrimPrefix(path, zw.src+"/")
	if wr, err = zw.zipWritter.Create(fileName); err != nil {
		return err
	}

	if _, err = io.Copy(wr, file); err != nil {
		return err
	}
	return nil
}

func (zw *ZipWritter) buildIncludeMetadata() ([]Include, error) {
	var includeInfo []Include
	var files []os.FileInfo
	var err error

	// iterate over set of included files specified in manifest YAML e.g.
	// include:
	// - ["source"]
	// - ["source", "destination"]
	for _, includeData := range zw.include {
		var i Include
		// if "destination" is not specified, its considered same as "source"
		// "source" is relative to where manifest.yaml file is located
		// relative source path is converted to absolute path by appending manifest path
		// since the relative source path might not be accessible from where wskdeploy is invoked
		// "destination" is relative to the action directory, the one specified in function
		// relative path is converted to absolute path by appending function directory
		if len(includeData) == 1 {
			i.source = filepath.Join(zw.manifestFilePath, includeData[0])
			i.destination = filepath.Join(zw.src, includeData[0])
		} else if len(includeData) == 2 {
			i.source = filepath.Join(zw.manifestFilePath, includeData[0])
			i.destination = zw.src + "/" + includeData[1]
		} else {
			continue
		}

		// TODO address scenarios with these patterns:
		// TODO actions/common/*
		// TODO actions/*
		// TODO actions/common/*/

		sourceDir := i.source
		if isFilePath(sourceDir) {
			sourceDir = filepath.Dir(i.source)
		}
		destDir := i.destination
		if isFilePath(destDir) {
			destDir = filepath.Dir(destDir)
		}

		// handle the scenarios where included path is something similar to actions/common/*.js
		if strings.HasPrefix(filepath.Base(i.source), "*.") {
			if files, err = ioutil.ReadDir(sourceDir); err != nil {
				return includeInfo, err
			}
			for _, file := range files {
				var j Include
				if filepath.Ext(file.Name()) == filepath.Ext(filepath.Base(i.source)) {
					j.source = filepath.Join(filepath.Dir(i.source), file.Name())
					j.destination = filepath.Join(destDir, file.Name())
					includeInfo = append(includeInfo, j)
				}
			}
			// handle scenarios where included path is something similar to actions/common/utils.js
			// and destination is set to ./common/ i.e. no file name specified in the destination
		} else if strings.HasSuffix(i.source, "*") {
			if files, err = ioutil.ReadDir(sourceDir); err != nil {
				return includeInfo, err
			}
			for _, file := range files {
				spew.Dump(file.Name())
			}
		} else {
			if f, err := isFile(i.source); err == nil && f {
				if _, file := filepath.Split(i.destination); len(file) == 0 {
					_, sFile := filepath.Split(i.source)
					i.destination = i.destination + sFile
					includeInfo = append(includeInfo, i)
				} else {
					includeInfo = append(includeInfo, i)
				}
			} else {
				// append just parsed include info to the list for further processing
				includeInfo = append(includeInfo, i)
			}
		}
	}
	return includeInfo, nil
}

func (zw *ZipWritter) Zip() error {

	var zipFile *os.File
	var err error
	var fileInfo os.FileInfo

	// create zip file e.g. greeting.zip
	if zipFile, err = os.Create(zw.des); err != nil {
		return err
	}
	defer zipFile.Close()

	// creating a new zip writter for greeting.zip
	zw.zipWritter = zip.NewWriter(zipFile)

	// walk file system rooted at the directory specified in "function"
	// walk over each file and dir under root directory e.g. function: actions/greeting
	// add actions/greeting/index.js and actions/greeting/package.json to zip file
	if err = filepath.Walk(zw.src, zw.zipFile); err != nil {
		return nil
	}

	// maintain a list of included files and/or directories with their destination
	var includeInfo []Include
	includeInfo, err = zw.buildIncludeMetadata()
	if err != nil {
		return err
	}

	for _, i := range includeInfo {
		if i.source != i.destination {
			// now determine whether the included item is file or dir
			// it could list something like this as well, "actions/common/*.js"
			if fileInfo, err = os.Stat(i.source); err != nil {
				return err
			}

			// if the included item is a directory, call a function to copy the
			// entire directory recursively including its subdirectories and files
			if fileInfo.Mode().IsDir() {
				if err = copyDir(i.source, i.destination); err != nil {
					return err
				}
				// if the included item is a file, call a function to copy the file
				// along with its path by creating the parent directories
			} else if fileInfo.Mode().IsRegular() {
				if err = copyFile(i.source, i.destination); err != nil {
					return err
				}
			}
		}
		// add included item into zip file greeting.zip
		if err = filepath.Walk(i.destination, zw.zipFile); err != nil {
			return nil
		}
	}

	// now close the zip file greeting.zip as all the included items
	// are added into the zip file along with the action root dir
	if err = zw.zipWritter.Close(); err != nil {
		return err
	}

	// and its safe to delete the files/directories which we copied earlier
	// to include them in the zip file greeting.zip
	for _, i := range includeInfo {
		os.RemoveAll(i.destination)
	}

	return nil
}
