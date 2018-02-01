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

package deployers

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
)

// name of directory that can contain source code
const FileSystemSourceDirectoryName = "actions"

type FileSystemReader struct {
	serviceDeployer *ServiceDeployer
}

func NewFileSystemReader(serviceDeployer *ServiceDeployer) *FileSystemReader {
	var reader FileSystemReader
	reader.serviceDeployer = serviceDeployer

	return &reader
}

func (reader *FileSystemReader) ReadProjectDirectory(manifest *parsers.YAML) ([]utils.ActionRecord, error) {

	// Inform user of what reader is doing
	dbgMsg := wski18n.T(wski18n.ID_DEBUG_PROJECT_SEARCH_X_path_X_key_X,
		map[string]interface{}{
			wski18n.KEY_PATH: reader.serviceDeployer.ProjectPath,
			wski18n.KEY_KEY:  wski18n.ACTIONS})
	wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, dbgMsg)

	projectPathCount, err := reader.getFilePathCount(reader.serviceDeployer.ProjectPath)
	actions := make([]utils.ActionRecord, 0)

	if err != nil {
		return actions, err
	}

	err = filepath.Walk(reader.serviceDeployer.ProjectPath, func(fpath string, f os.FileInfo, err error) error {
		if fpath != reader.serviceDeployer.ProjectPath {
			pathCount, err := reader.getFilePathCount(fpath)
			if err != nil {
				return wskderrors.NewFileReadError(fpath, err.Error())
			}

			if !f.IsDir() {
				if pathCount-projectPathCount == 1 || strings.HasPrefix(fpath, reader.serviceDeployer.ProjectPath+"/"+FileSystemSourceDirectoryName) {
					ext := filepath.Ext(fpath)

					// TODO(#692) do not hardcoded known extensions here, create a util that associates
					// known extensions to runtime types
					foundFile := false
					switch ext {
					case ".swift":
						foundFile = true
					case ".js":
						foundFile = true
					case ".py":
						foundFile = true
					}

					if foundFile == true {
						_, action, err := reader.CreateActionFromFile(reader.serviceDeployer.ManifestPath, fpath)
						if err != nil {
							return wskderrors.NewFileReadError(fpath, err.Error())
						}

						var record utils.ActionRecord
						record.Action = action
						record.Packagename = manifest.Package.Packagename
						record.Filepath = fpath

						actions = append(actions, record)
					}
				}
			} else if strings.HasPrefix(fpath, reader.serviceDeployer.ProjectPath+"/"+FileSystemSourceDirectoryName) {
				// Inform user of what reader is doing
				dbgMsg := wski18n.T(wski18n.ID_DEBUG_PROJECT_SEARCH_X_path_X_key_X,
					map[string]interface{}{
						wski18n.KEY_PATH: filepath.Base(fpath),
						wski18n.KEY_KEY:  wski18n.ACTION_CODE})
				wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, dbgMsg)

				// TODO(#694) It appears there is no code to do what the debug message suggests
				// TODO(#694) continued: i.e., searching for Action code...

			} else {
				return filepath.SkipDir
			}

		}
		return err
	})

	if err != nil {
		return actions, err
	}

	return actions, nil

}

func (reader *FileSystemReader) CreateActionFromFile(manipath, filePath string) (string, *whisk.Action, error) {
	ext := filepath.Ext(filePath)
	baseName := filepath.Base(filePath)
	name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	action := new(whisk.Action)

	// TODO(#692) same TODO as above, that is create util for matching exts. to runtimes in common util pkg.
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

		dat, err := new(utils.ContentReader).LocalReader.ReadLocal(filePath)
		if err != nil {
			return name, action, wskderrors.NewFileReadError(filePath, err.Error())
		}

		action.Exec = new(whisk.Exec)
		code := string(dat)
		action.Exec.Code = &code
		action.Exec.Kind = kind
		action.Name = name
		pub := false
		action.Publish = &pub
		return name, action, nil
	}

	// TODO(#691) create new named error
	// If the action is not supported, we better to return an error.
	return "", nil, errors.New("Unsupported action type.")
}

func (reader *FileSystemReader) getFilePathCount(path string) (int, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return 0, err
	}

	pathList := strings.Split(absPath, "/")
	return len(pathList) - 1, nil
}

func (reader *FileSystemReader) SetFileActions(actions []utils.ActionRecord) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, fileAction := range actions {
		existAction, exists := reader.serviceDeployer.Deployment.Packages[fileAction.Packagename].Actions[fileAction.Action.Name]

		if exists == true {
			if existAction.Filepath == fileAction.Filepath || existAction.Filepath == "" {
				// we're adding a filesystem detected action so just updated code and filepath if needed
				existAction.Action.Exec.Code = fileAction.Action.Exec.Code
				existAction.Filepath = fileAction.Filepath
			} else {
				// TODO(#691) create new named error
				// Action exists, but references two different sources
				return errors.New("Conflict detected for action named " + existAction.Action.Name + ". Found two locations for source file: " + existAction.Filepath + " and " + fileAction.Filepath)
			}
		} else {
			// not a new action so to actions in package
			reader.serviceDeployer.Deployment.Packages[fileAction.Packagename].Actions[fileAction.Action.Name] = fileAction
		}
	}

	return nil
}
