package deployers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/parsers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
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

func (reader *FileSystemReader) ReadProjectDirectory(manifest *parsers.ManifestYAML) ([]utils.ActionRecord, error) {

	fmt.Println("Inspecting project directory for actions....")

	projectPathCount, err := reader.getFilePathCount(reader.serviceDeployer.ProjectPath)
	utils.Check(err)

	actions := make([]utils.ActionRecord, 0)

	err = filepath.Walk(reader.serviceDeployer.ProjectPath, func(fpath string, f os.FileInfo, err error) error {
		if fpath != reader.serviceDeployer.ProjectPath {
			pathCount, err := reader.getFilePathCount(fpath)
			utils.Check(err)

			if !f.IsDir() {
				if pathCount-projectPathCount == 1 || strings.HasPrefix(fpath, reader.serviceDeployer.ProjectPath+"/"+FileSystemSourceDirectoryName) {
					ext := filepath.Ext(fpath)

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
						utils.Check(err)

						var record utils.ActionRecord
						record.Action = action
						record.Packagename = manifest.Package.Packagename
						record.Filepath = fpath

						actions = append(actions, record)
					}
				}
			} else if strings.HasPrefix(fpath, reader.serviceDeployer.ProjectPath+"/"+FileSystemSourceDirectoryName) {
				fmt.Println("Searching directory " + filepath.Base(fpath) + " for action source code.")
			} else {
				return filepath.SkipDir
			}

		}
		return err
	})

	if err != nil {
		return nil, err
	}

	return actions, nil

}

func (reader *FileSystemReader) CreateActionFromFile(manipath, filePath string) (string, *whisk.Action, error) {
	ext := filepath.Ext(filePath)
	baseName := filepath.Base(filePath)
	name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	action := new(whisk.Action)

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
		utils.Check(err)

		action.Exec = new(whisk.Exec)
		code := string(dat)
		action.Exec.Code = &code
		action.Exec.Kind = kind
		action.Name = name
		pub := false
		action.Publish = &pub
		return name, action, nil
	}
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
