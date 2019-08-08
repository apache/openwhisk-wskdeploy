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
	"os"
	"path/filepath"
	"strings"

	"github.com/apache/openwhisk-wskdeploy/wski18n"
	"github.com/apache/openwhisk-wskdeploy/wskprint"
)

const PATH_WILDCARD = "*"
const ONE_DIR_UP = "../"

func NewZipWritter(src string, des string, include [][]string, exclude []string, manifestFilePath string) *ZipWritter {
	zw := &ZipWritter{
		src:              src,
		des:              des,
		include:          include,
		exclude:          exclude,
		excludedFiles:    make(map[string]bool, 0),
		manifestFilePath: manifestFilePath,
	}
	return zw
}

type ZipWritter struct {
	src              string
	des              string
	include          [][]string
	exclude          []string
	excludedFiles    map[string]bool
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
	var verboseMsg string

	if err != nil {
		return err
	}

	if zw.excludedFiles[filepath.Clean(path)] {
		verboseMsg = wski18n.T(wski18n.ID_VERBOSE_ZIP_EXCLUDING_FILE_X_path_X,
			map[string]interface{}{
				wski18n.KEY_PATH: path,
			})
		wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, verboseMsg)
		return nil
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
	verboseMsg = wski18n.T(wski18n.ID_VERBOSE_ZIP_ADDING_FILE_X_path_X,
		map[string]interface{}{
			wski18n.KEY_PATH: path,
		})
	wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, verboseMsg)
	return nil
}

func (zw *ZipWritter) buildIncludeMetadata() ([]Include, error) {
	var includeInfo []Include
	var listOfSourceFiles []string
	var err error
	var verboseMsg string

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
			verboseMsg = wski18n.T(wski18n.ID_VERBOSE_ZIP_INCLUDE_SOURCE_PATH_X_path_X,
				map[string]interface{}{
					wski18n.KEY_PATH: includeData[0],
				})
			wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, verboseMsg)
		} else if len(includeData) == 2 {
			i.source = filepath.Join(zw.manifestFilePath, includeData[0])
			i.destination = zw.src + "/" + includeData[1]
			verboseMsg = wski18n.T(wski18n.ID_VERBOSE_ZIP_INCLUDE_SOURCE_PATH_X_path_X_DESTINATION_PATH_X_dest_X,
				map[string]interface{}{
					wski18n.KEY_PATH:        includeData[0],
					wski18n.KEY_DESTINATION: includeData[1],
				})
		} else {
			if len(includeData) == 0 {
				verboseMsg = wski18n.T(wski18n.ID_VERBOSE_INVALID_INCLUDE_ENTRY,
					map[string]interface{}{
						wski18n.KEY_INCLUDE: "",
					})
				wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, verboseMsg)
			} else {
				for index, d := range includeData {
					includeData[index] = "\"" + d + "\""
				}
				includeEntry := strings.Join(includeData, ", ")
				verboseMsg = wski18n.T(wski18n.ID_VERBOSE_INVALID_INCLUDE_ENTRY,
					map[string]interface{}{
						wski18n.KEY_INCLUDE: includeEntry,
					})
				wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, verboseMsg)
			}
			continue
		}

		// set destDir to the destination location
		// check if its a file than change it to the Dir of destination file
		destDir := i.destination
		if isFilePath(destDir) {
			destDir = filepath.Dir(destDir)
		}
		// trim path wildcard "*" from the destination path as if it has any
		destDirs := strings.Split(destDir, PATH_WILDCARD)
		destDir = destDirs[0]

		// retrieve the name of all files matching pattern or nil if there is no matching file
		// listOfSourceFiles will hold a list of files matching patterns such as
		// actions/* or actions/libs/* or actions/libs/*/utils.js or actions/*/*/utils.js
		if listOfSourceFiles, err = filepath.Glob(i.source); err != nil {
			return includeInfo, err
		}

		// handle the scenarios where included path is something similar to actions/common/*.js
		// or actions/libs/* or actions/libs/*/utils.js
		// and destination is set to libs/ or libs/* or ./libs/* or libs/*/utils.js or libs/ or ./libs/
		if strings.ContainsAny(i.source, PATH_WILDCARD) {
			wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, wski18n.T(wski18n.ID_VERBOSE_LIST_OF_FILES_MATCHING_PATTERN))
			for _, file := range listOfSourceFiles {
				var relPath string
				if relPath, err = filepath.Rel(i.source, file); err != nil {
					return includeInfo, err
				}
				relPath = strings.TrimLeft(relPath, ONE_DIR_UP)
				j := Include{
					source:      file,
					destination: filepath.Join(destDir, relPath),
				}
				includeInfo = append(includeInfo, j)
				zw.excludedFiles[j.source] = false
				verboseMsg = wski18n.T(wski18n.ID_VERBOSE_ZIP_INCLUDE_SOURCE_PATH_X_path_X_DESTINATION_PATH_X_dest_X,
					map[string]interface{}{
						wski18n.KEY_PATH:        j.source,
						wski18n.KEY_DESTINATION: j.destination,
					})
				wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, verboseMsg)
			}
			// handle scenarios where included path is something similar to actions/common/utils.js
			// and destination is set to ./common/ i.e. no file name specified in the destination
		} else {
			if f, err := isFile(i.source); err == nil && f {
				if _, file := filepath.Split(i.destination); len(file) == 0 {
					_, sFile := filepath.Split(i.source)
					i.destination = i.destination + sFile
				}
			}
			// append just parsed include info to the list for further processing
			wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, wski18n.T(wski18n.ID_VERBOSE_LIST_OF_FILES_MATCHING_PATTERN))
			verboseMsg = wski18n.T(wski18n.ID_VERBOSE_ZIP_INCLUDE_SOURCE_PATH_X_path_X_DESTINATION_PATH_X_dest_X,
				map[string]interface{}{
					wski18n.KEY_PATH:        i.source,
					wski18n.KEY_DESTINATION: i.destination,
				})
			wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, verboseMsg)
			includeInfo = append(includeInfo, i)
			zw.excludedFiles[i.source] = false
		}
	}
	return includeInfo, nil
}

func (zw *ZipWritter) buildExcludeMetadata() error {
	var err error
	for _, exclude := range zw.exclude {
		exclude = filepath.Join(zw.manifestFilePath, exclude)
		if err = zw.findExcludedIncludedFiles(exclude, true); err != nil {
			return err
		}
	}
	return err
}

func (zw *ZipWritter) findExcludedIncludedFiles(functionPath string, flag bool) error {
	var err error
	var files []string
	var excludedFiles []string
	var f bool

	if !strings.HasSuffix(functionPath, PATH_WILDCARD) {
		functionPath = filepath.Join(functionPath, PATH_WILDCARD)
	}
	if excludedFiles, err = filepath.Glob(functionPath); err != nil {
		return err
	}
	for _, file := range excludedFiles {
		err = filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			return err
		}
	}
	for _, file := range files {
		if f, err = isFile(file); err != nil {
			return err
		} else if f {
			zw.excludedFiles[file] = flag
		} else {
			if err = zw.findExcludedIncludedFiles(file, flag); err != nil {
				return err
			}
		}
	}
	return err
}

func (zw *ZipWritter) Zip() error {

	var zipFile *os.File
	var err error
	var fileInfo os.FileInfo
	var verboseMsg string

	// create zip file e.g. greeting.zip
	if zipFile, err = os.Create(zw.des); err != nil {
		return err
	}
	defer zipFile.Close()

	verboseMsg = wski18n.T(wski18n.ID_VERBOSE_CREATING_ZIP_FILE_X_path_X,
		map[string]interface{}{
			wski18n.KEY_PATH: zipFile.Name(),
		})
	wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, verboseMsg)

	// creating a new zip writter for greeting.zip
	zw.zipWritter = zip.NewWriter(zipFile)

	// build a map of file names and bool indicating whether the file is included or excluded
	// iterate over the directory specified in "function", find the list of files and mark them as not excluded
	if err = zw.findExcludedIncludedFiles(zw.src, false); err != nil {
		return err
	}

	if err = zw.buildExcludeMetadata(); err != nil {
		return err
	}

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
		if filepath.Clean(i.source) != filepath.Clean(i.destination) {
			verboseMsg = wski18n.T(wski18n.ID_VERBOSE_DELETING_FILE_X_path_X,
				map[string]interface{}{
					wski18n.KEY_PATH: i.destination,
				})
			wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, verboseMsg)
			os.RemoveAll(i.destination)
		}
	}

	return nil
}
