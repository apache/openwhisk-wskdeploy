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

package parsers

import (
	"encoding/base64"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskenv"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
)

const (
	PATH_SEPERATOR  = "/"
	API             = "API"
	HTTPS           = "https"
	HTTP            = "http"
	API_VERSION     = "v1"
	WEB             = "web"
	DEFAULT_PACKAGE = "default"
)

// Read existing manifest file or create new if none exists
func ReadOrCreateManifest() (*YAML, error) {
	maniyaml := YAML{}

	if _, err := os.Stat(utils.ManifestFileNameYaml); err == nil {
		dat, _ := ioutil.ReadFile(utils.ManifestFileNameYaml)
		err := NewYAMLParser().Unmarshal(dat, &maniyaml)
		if err != nil {
			return &maniyaml, wskderrors.NewFileReadError(utils.ManifestFileNameYaml, err.Error())
		}
	}
	return &maniyaml, nil
}

// Serialize manifest to local file
func Write(manifest *YAML, filename string) error {
	output, err := NewYAMLParser().marshal(manifest)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(filename, err.Error())
	}

	f, err := os.Create(filename)
	if err != nil {
		return wskderrors.NewFileReadError(filename, err.Error())
	}
	defer f.Close()

	f.Write(output)
	return nil
}

func (dm *YAMLParser) Unmarshal(input []byte, manifest *YAML) error {
	err := yaml.UnmarshalStrict(input, manifest)
	if err != nil {
		return err
	}
	return nil
}

func (dm *YAMLParser) marshal(manifest *YAML) (output []byte, err error) {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (dm *YAMLParser) ParseManifest(manifestPath string) (*YAML, error) {
	mm := NewYAMLParser()
	maniyaml := YAML{}

	content, err := utils.Read(manifestPath)
	if err != nil {
		return &maniyaml, wskderrors.NewFileReadError(manifestPath, err.Error())
	}

	err = mm.Unmarshal(content, &maniyaml)
	if err != nil {
		return &maniyaml, wskderrors.NewYAMLParserErr(manifestPath, err)
	}
	maniyaml.Filepath = manifestPath
	manifest := ReadEnvVariable(&maniyaml)

	return manifest, nil
}

func (dm *YAMLParser) ComposeDependenciesFromAllPackages(manifest *YAML, projectPath string, filePath string) (map[string]utils.DependencyRecord, error) {
	dependencies := make(map[string]utils.DependencyRecord)
	packages := make(map[string]Package)

	if len(manifest.Packages) != 0 {
		packages = manifest.Packages
	} else {
		packages = manifest.GetProject().Packages
	}

	for n, p := range packages {
		d, err := dm.ComposeDependencies(p, projectPath, filePath, n)
		if err == nil {
			for k, v := range d {
				dependencies[k] = v
			}
		} else {
			return nil, err
		}
	}
	return dependencies, nil
}

func (dm *YAMLParser) ComposeDependencies(pkg Package, projectPath string, filePath string, packageName string) (map[string]utils.DependencyRecord, error) {

	var errorParser error
	depMap := make(map[string]utils.DependencyRecord)
	for key, dependency := range pkg.Dependencies {
		version := dependency.Version
		if version == "" {
			// TODO() interactive ask for branch, AND consider YAML specification to allow key for branch
			version = YAML_VALUE_BRANCH_MASTER
		}

		location := dependency.Location

		isBinding := false
		if utils.LocationIsBinding(location) {

			if !strings.HasPrefix(location, "/") {
				location = "/" + dependency.Location
			}

			isBinding = true
		} else if utils.LocationIsGithub(location) {

			// TODO() define const for the protocol prefix, etc.
			if !strings.HasPrefix(location, "https://") && !strings.HasPrefix(location, "http://") {
				location = "https://" + dependency.Location
				location = wskenv.InterpolateStringWithEnvVar(location).(string)
			}

			isBinding = false
		} else {
			// TODO() create new named error in wskerrors package
			return nil, errors.New(wski18n.T(wski18n.ID_ERR_DEPENDENCY_UNKNOWN_TYPE))
		}

		keyValArrParams := make(whisk.KeyValueArr, 0)
		for name, param := range dependency.Inputs {
			var keyVal whisk.KeyValue
			keyVal.Key = name

			keyVal.Value, errorParser = ResolveParameter(name, &param, filePath)

			if errorParser != nil {
				return nil, errorParser
			}

			if keyVal.Value != nil {
				keyValArrParams = append(keyValArrParams, keyVal)
			}
		}

		keyValArrAnot := make(whisk.KeyValueArr, 0)
		for name, value := range dependency.Annotations {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value = wskenv.InterpolateStringWithEnvVar(value)

			keyValArrAnot = append(keyValArrAnot, keyVal)
		}

		packDir := path.Join(projectPath, "Packages")
		depName := packageName + ":" + key
		depMap[depName] = utils.NewDependencyRecord(packDir, packageName, location, version, keyValArrParams, keyValArrAnot, isBinding)
	}

	return depMap, nil
}

func (dm *YAMLParser) ComposeAllPackages(manifest *YAML, filePath string, ma whisk.KeyValue) (map[string]*whisk.Package, error) {
	packages := map[string]*whisk.Package{}
	manifestPackages := make(map[string]Package)

	if len(manifest.Packages) != 0 {
		manifestPackages = manifest.Packages
	} else {
		manifestPackages = manifest.GetProject().Packages
	}

	if len(manifestPackages) == 0 {
		warningString := wski18n.T(
			wski18n.ID_WARN_PACKAGES_NOT_FOUND_X_path_X,
			map[string]interface{}{
				wski18n.KEY_PATH: manifest.Filepath})
		wskprint.PrintOpenWhiskWarning(warningString)
	}

	// Compose each package found in manifest
	for n, p := range manifestPackages {
		s, err := dm.ComposePackage(p, n, filePath, ma)

		if err == nil {
			packages[n] = s
		} else {
			return nil, err
		}
	}

	return packages, nil
}

func (dm *YAMLParser) ComposePackage(pkg Package, packageName string, filePath string, ma whisk.KeyValue) (*whisk.Package, error) {
	var errorParser error
	pag := &whisk.Package{}
	pag.Name = packageName
	//The namespace for this package is absent, so we use default guest here.
	pag.Namespace = pkg.Namespace
	pub := false
	pag.Publish = &pub

	//Version is a mandatory value
	//If it is an empty string, it will be set to default value
	//And print an warning message
	// TODO(#673) implement STRICT flag
	if pkg.Version == "" {
		warningString := wski18n.T(
			wski18n.ID_WARN_KEY_MISSING_X_key_X_value_X,
			map[string]interface{}{
				wski18n.KEY_KEY:   wski18n.PACKAGE_VERSION,
				wski18n.KEY_VALUE: DEFAULT_PACKAGE_VERSION})
		wskprint.PrintOpenWhiskWarning(warningString)

		warningString = wski18n.T(
			wski18n.ID_WARN_KEYVALUE_NOT_SAVED_X_key_X,
			map[string]interface{}{wski18n.KEY_KEY: wski18n.PACKAGE_VERSION})

		wskprint.PrintOpenWhiskWarning(warningString)
		pkg.Version = DEFAULT_PACKAGE_VERSION
	}

	//License is a mandatory value
	//set license to unknown if it is an empty string
	//And print an warning message
	// TODO(#673) implement STRICT flag
	if pkg.License == "" {
		warningString := wski18n.T(
			wski18n.ID_WARN_KEY_MISSING_X_key_X_value_X,
			map[string]interface{}{
				wski18n.KEY_KEY:   wski18n.PACKAGE_LICENSE,
				wski18n.KEY_VALUE: DEFAULT_PACKAGE_LICENSE})
		wskprint.PrintOpenWhiskWarning(warningString)

		warningString = wski18n.T(
			wski18n.ID_WARN_KEYVALUE_NOT_SAVED_X_key_X,
			map[string]interface{}{wski18n.KEY_KEY: wski18n.PACKAGE_VERSION})

		wskprint.PrintOpenWhiskWarning(warningString)

		pkg.License = DEFAULT_PACKAGE_LICENSE
	} else {
		utils.CheckLicense(pkg.License)
	}

	//set parameters
	keyValArr := make(whisk.KeyValueArr, 0)
	for name, param := range pkg.Inputs {
		var keyVal whisk.KeyValue
		keyVal.Key = name

		keyVal.Value, errorParser = ResolveParameter(name, &param, filePath)

		if errorParser != nil {
			return nil, errorParser
		}

		if keyVal.Value != nil {
			keyValArr = append(keyValArr, keyVal)
		}
	}

	if len(keyValArr) > 0 {
		pag.Parameters = keyValArr
	}

	// set Package Annotations
	listOfAnnotations := make(whisk.KeyValueArr, 0)
	for name, value := range pkg.Annotations {
		var keyVal whisk.KeyValue
		keyVal.Key = name
		keyVal.Value = wskenv.InterpolateStringWithEnvVar(value)
		listOfAnnotations = append(listOfAnnotations, keyVal)
	}
	if len(listOfAnnotations) > 0 {
		pag.Annotations = append(pag.Annotations, listOfAnnotations...)
	}

	// add Managed Annotations if this is Managed Deployment
	if utils.Flags.Managed {
		pag.Annotations = append(pag.Annotations, ma)
	}

	// "default" package is a reserved package name
	// and in this case wskdeploy deploys openwhisk entities under
	// /namespace instead of /namespace/package
	if strings.ToLower(pag.Name) == DEFAULT_PACKAGE {
		wskprint.PrintlnOpenWhiskInfo(wski18n.T(wski18n.ID_MSG_DEFAULT_PACKAGE))
	}

	return pag, nil
}

func (dm *YAMLParser) ComposeSequencesFromAllPackages(namespace string, mani *YAML, ma whisk.KeyValue) ([]utils.ActionRecord, error) {
	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)

	manifestPackages := make(map[string]Package)

	if len(mani.Packages) != 0 {
		manifestPackages = mani.Packages
	} else {
		manifestPackages = mani.GetProject().Packages
	}

	for n, p := range manifestPackages {
		s, err := dm.ComposeSequences(namespace, p.Sequences, n, ma)
		if err == nil {
			s1 = append(s1, s...)
		} else {
			return nil, err
		}
	}
	return s1, nil
}

func (dm *YAMLParser) ComposeSequences(namespace string, sequences map[string]Sequence, packageName string, ma whisk.KeyValue) ([]utils.ActionRecord, error) {
	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)

	for key, sequence := range sequences {
		wskaction := new(whisk.Action)
		wskaction.Exec = new(whisk.Exec)
		wskaction.Exec.Kind = YAML_KEY_SEQUENCE
		actionList := strings.Split(sequence.Actions, ",")

		var components []string
		for _, a := range actionList {
			act := strings.TrimSpace(a)
			if !strings.ContainsRune(act, '/') && !strings.HasPrefix(act, packageName+"/") &&
				strings.ToLower(packageName) != DEFAULT_PACKAGE {
				act = path.Join(packageName, act)
			}
			components = append(components, path.Join("/"+namespace, act))
		}

		wskaction.Exec.Components = components
		wskaction.Name = key
		pub := false
		wskaction.Publish = &pub
		wskaction.Namespace = namespace

		keyValArr := make(whisk.KeyValueArr, 0)
		for name, value := range sequence.Annotations {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value = wskenv.InterpolateStringWithEnvVar(value)

			keyValArr = append(keyValArr, keyVal)
		}

		if len(keyValArr) > 0 {
			wskaction.Annotations = keyValArr
		}

		// appending managed annotations if its a managed deployment
		if utils.Flags.Managed {
			wskaction.Annotations = append(wskaction.Annotations, ma)
		}

		record := utils.ActionRecord{Action: wskaction, Packagename: packageName, Filepath: key}
		s1 = append(s1, record)
	}
	return s1, nil
}

func (dm *YAMLParser) ComposeActionsFromAllPackages(manifest *YAML, filePath string, ma whisk.KeyValue) ([]utils.ActionRecord, error) {
	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)
	manifestPackages := make(map[string]Package)

	if len(manifest.Packages) != 0 {
		manifestPackages = manifest.Packages
	} else {
		manifestPackages = manifest.GetProject().Packages
	}

	for n, p := range manifestPackages {
		a, err := dm.ComposeActions(filePath, p.Actions, n, ma)
		if err == nil {
			s1 = append(s1, a...)
		} else {
			return nil, err
		}
	}
	return s1, nil
}

func (dm *YAMLParser) ComposeActions(filePath string, actions map[string]Action, packageName string, ma whisk.KeyValue) ([]utils.ActionRecord, error) {

	var errorParser error
	var ext string
	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)

	for key, action := range actions {
		splitFilePath := strings.Split(filePath, string(os.PathSeparator))
		// set the name of the action (which is the key)
		action.Name = key

		// Create action data object with CLI
		wskaction := new(whisk.Action)
		wskaction.Exec = new(whisk.Exec)

		/*
		 *  Action.Function
		 */
		//set action.Function to action.Location
		//because Location is deprecated in Action entity
		if action.Function == "" && action.Location != "" {
			action.Function = action.Location
		}

		//bind action, and exposed URL
		if action.Function != "" {

			filePath := strings.TrimRight(filePath, splitFilePath[len(splitFilePath)-1]) + action.Function

			if utils.IsDirectory(filePath) {
				// TODO() define ext as const
				zipName := filePath + ".zip"
				err := utils.NewZipWritter(filePath, zipName).Zip()
				if err != nil {
					return nil, err
				}
				// TODO() do not use defer in a loop, resource leaks possible
				defer os.Remove(zipName)
				// TODO(): support docker and main entry as did by go cli?
				wskaction.Exec, err = utils.GetExec(zipName, action.Runtime, false, "")
				if err != nil {
					return nil, err
				}
			} else {
				ext = path.Ext(filePath)
				// drop the "." from file extension
				if len(ext) > 0 && ext[0] == '.' {
					ext = ext[1:]
				}

				// determine default runtime for the given file extension
				var kind string
				r := utils.FileExtensionRuntimeKindMap[ext]
				kind = utils.DefaultRunTimes[r]

				// produce an error when a runtime could not be derived from the action file extension
				// and its not explicitly specified in the manifest YAML file
				// and action source is not a zip file
				if len(kind) == 0 && len(action.Runtime) == 0 && ext != utils.ZIP_FILE_EXTENSION {
					errMessage := wski18n.T(wski18n.ID_ERR_RUNTIME_MISMATCH_X_runtime_X_ext_X_action_X,
						map[string]interface{}{
							wski18n.KEY_RUNTIME:   action.Runtime,
							wski18n.KEY_EXTENTION: ext,
							wski18n.KEY_ACTION:    action.Name})
					return nil, wskderrors.NewInvalidRuntimeError(errMessage,
						splitFilePath[len(splitFilePath)-1], action.Name,
						action.Runtime,
						utils.ListOfSupportedRuntimes(utils.SupportedRunTimes))
				}

				wskaction.Exec.Kind = kind

				action.Function = filePath
				dat, err := utils.Read(filePath)
				if err != nil {
					return s1, err
				}
				code := string(dat)
				if ext == utils.ZIP_FILE_EXTENSION || ext == utils.JAR_FILE_EXTENSION {
					code = base64.StdEncoding.EncodeToString([]byte(dat))
				}
				if ext == utils.ZIP_FILE_EXTENSION && len(action.Runtime) == 0 {
					errMessage := wski18n.T(wski18n.ID_ERR_RUNTIME_INVALID_X_runtime_X_action_X,
						map[string]interface{}{
							wski18n.KEY_RUNTIME: action.Runtime,
							wski18n.KEY_ACTION:  action.Name})
					return nil, wskderrors.NewInvalidRuntimeError(errMessage,
						splitFilePath[len(splitFilePath)-1],
						action.Name,
						action.Runtime,
						utils.ListOfSupportedRuntimes(utils.SupportedRunTimes))
				}
				wskaction.Exec.Code = &code
			}

		}

		/*
			 		 *  Action.Runtime
					 *  Perform few checks if action runtime is specified in manifest YAML file
					 *  (1) Check if specified runtime is one of the supported runtimes by OpenWhisk server
					 *  (2) Check if specified runtime is consistent with action source file extensions
					 *  Set the action runtime to match with the source file extension, if wskdeploy is not invoked in strict mode
		*/
		if action.Runtime != "" {
			if utils.CheckExistRuntime(action.Runtime, utils.SupportedRunTimes) {
				// for zip actions, rely on the runtimes from the manifest file as it can not be derived from the action source file extension
				// pick runtime from manifest file if its supported by OpenWhisk server
				if ext == utils.ZIP_FILE_EXTENSION {
					wskaction.Exec.Kind = action.Runtime
				} else {
					if utils.CheckRuntimeConsistencyWithFileExtension(ext, action.Runtime) {
						wskaction.Exec.Kind = action.Runtime
					} else {
						warnStr := wski18n.T(wski18n.ID_ERR_RUNTIME_MISMATCH_X_runtime_X_ext_X_action_X,
							map[string]interface{}{
								wski18n.KEY_RUNTIME:   action.Runtime,
								wski18n.KEY_EXTENTION: ext,
								wski18n.KEY_ACTION:    action.Name})
						wskprint.PrintOpenWhiskWarning(warnStr)

						// even if runtime is not consistent with file extension, deploy action with specified runtime in strict mode
						if utils.Flags.Strict {
							wskaction.Exec.Kind = action.Runtime
						} else {
							warnStr := wski18n.T(wski18n.ID_WARN_RUNTIME_CHANGED_X_runtime_X_action_X,
								map[string]interface{}{
									wski18n.KEY_RUNTIME: wskaction.Exec.Kind,
									wski18n.KEY_ACTION:  action.Name})
							wskprint.PrintOpenWhiskWarning(warnStr)
						}
					}
				}
			} else {
				warnStr := wski18n.T(wski18n.ID_ERR_RUNTIME_INVALID_X_runtime_X_action_X,
					map[string]interface{}{
						wski18n.KEY_RUNTIME: action.Runtime,
						wski18n.KEY_ACTION:  action.Name})
				wskprint.PrintOpenWhiskWarning(warnStr)

				if ext == utils.ZIP_FILE_EXTENSION {
					// for zip action, error out if specified runtime is not supported by
					// OpenWhisk server
					return nil, wskderrors.NewInvalidRuntimeError(warnStr,
						splitFilePath[len(splitFilePath)-1],
						action.Name,
						action.Runtime,
						utils.ListOfSupportedRuntimes(utils.SupportedRunTimes))
				} else {
					warnStr := wski18n.T(wski18n.ID_WARN_RUNTIME_CHANGED_X_runtime_X_action_X,
						map[string]interface{}{
							wski18n.KEY_RUNTIME: wskaction.Exec.Kind,
							wski18n.KEY_ACTION:  action.Name})
					wskprint.PrintOpenWhiskWarning(warnStr)
				}

			}
		}

		// we can specify the name of the action entry point using main
		if action.Main != "" {
			wskaction.Exec.Main = action.Main
		}

		/*
		 *  Action.Inputs
		 */
		keyValArr := make(whisk.KeyValueArr, 0)
		for name, param := range action.Inputs {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value, errorParser = ResolveParameter(name, &param, filePath)

			if errorParser != nil {
				return nil, errorParser
			}

			if keyVal.Value != nil {
				keyValArr = append(keyValArr, keyVal)
			}
		}

		// if we have successfully parser valid key/value parameters
		if len(keyValArr) > 0 {
			wskaction.Parameters = keyValArr
		}

		/*
		 *  Action.Outputs
		 */
		keyValArr = make(whisk.KeyValueArr, 0)
		for name, param := range action.Outputs {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value, errorParser = ResolveParameter(name, &param, filePath)

			// short circuit on error
			if errorParser != nil {
				return nil, errorParser
			}

			if keyVal.Value != nil {
				keyValArr = append(keyValArr, keyVal)
			}
		}

		// TODO{} add outputs as annotations (work to discuss officially supporting for compositions)
		if len(keyValArr) > 0 {
			// TODO() ?
			//wskaction.Annotations  // TBD
		}

		/*
		 *  Action.Annotations
		 */
		listOfAnnotations := make(whisk.KeyValueArr, 0)
		for name, value := range action.Annotations {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value = wskenv.InterpolateStringWithEnvVar(value)
			listOfAnnotations = append(listOfAnnotations, keyVal)
		}
		if len(listOfAnnotations) > 0 {
			wskaction.Annotations = append(wskaction.Annotations, listOfAnnotations...)
		}
		// add managed annotations if its marked as managed deployment
		if utils.Flags.Managed {
			wskaction.Annotations = append(wskaction.Annotations, ma)
		}

		/*
		 *  Web Export
		 */
		// Treat ACTION as a web action, a raw HTTP web action, or as a standard action based on web-export;
		// when web-export is set to yes | true, treat action as a web action,
		// when web-export is set to raw, treat action as a raw HTTP web action,
		// when web-export is set to no | false, treat action as a standard action
		if len(action.Webexport) != 0 {
			wskaction.Annotations, errorParser = utils.WebAction(filePath, action.Name, action.Webexport, listOfAnnotations, false)
			if errorParser != nil {
				return s1, errorParser
			}
		}

		/*
		 *  Action.Limits
		 */
		if action.Limits != nil {
			wsklimits := new(whisk.Limits)

			// TODO() use LIMITS_SUPPORTED in yamlparser to enumerata through instead of hardcoding
			// perhaps change into a tuple
			if utils.LimitsTimeoutValidation(action.Limits.Timeout) {
				wsklimits.Timeout = action.Limits.Timeout
			} else {
				warningString := wski18n.T(wski18n.ID_WARN_LIMIT_IGNORED_X_limit_X,
					map[string]interface{}{wski18n.KEY_LIMIT: LIMIT_VALUE_TIMEOUT})
				wskprint.PrintOpenWhiskWarning(warningString)
			}
			if utils.LimitsMemoryValidation(action.Limits.Memory) {
				wsklimits.Memory = action.Limits.Memory
			} else {
				warningString := wski18n.T(wski18n.ID_WARN_LIMIT_IGNORED_X_limit_X,
					map[string]interface{}{wski18n.KEY_LIMIT: LIMIT_VALUE_MEMORY_SIZE})
				wskprint.PrintOpenWhiskWarning(warningString)
			}
			if utils.LimitsLogsizeValidation(action.Limits.Logsize) {
				wsklimits.Logsize = action.Limits.Logsize
			} else {
				warningString := wski18n.T(wski18n.ID_WARN_LIMIT_IGNORED_X_limit_X,
					map[string]interface{}{wski18n.KEY_LIMIT: LIMIT_VALUE_LOG_SIZE})
				wskprint.PrintOpenWhiskWarning(warningString)
			}
			if wsklimits.Timeout != nil || wsklimits.Memory != nil || wsklimits.Logsize != nil {
				wskaction.Limits = wsklimits
			}

			// TODO() use LIMITS_UNSUPPORTED in yamlparser to enumerata through instead of hardcoding
			// emit warning errors if these limits are not nil
			utils.NotSupportLimits(action.Limits.ConcurrentActivations, LIMIT_VALUE_CONCURRENT_ACTIVATIONS)
			utils.NotSupportLimits(action.Limits.UserInvocationRate, LIMIT_VALUE_USER_INVOCATION_RATE)
			utils.NotSupportLimits(action.Limits.CodeSize, LIMIT_VALUE_CODE_SIZE)
			utils.NotSupportLimits(action.Limits.ParameterSize, LIMIT_VALUE_PARAMETER_SIZE)
		}

		wskaction.Name = key
		pub := false
		wskaction.Publish = &pub

		record := utils.ActionRecord{Action: wskaction, Packagename: packageName, Filepath: action.Function}
		s1 = append(s1, record)
	}

	return s1, nil

}

func (dm *YAMLParser) ComposeTriggersFromAllPackages(manifest *YAML, filePath string, ma whisk.KeyValue) ([]*whisk.Trigger, error) {
	var triggers []*whisk.Trigger = make([]*whisk.Trigger, 0)
	manifestPackages := make(map[string]Package)

	if len(manifest.Packages) != 0 {
		manifestPackages = manifest.Packages
	} else {
		manifestPackages = manifest.GetProject().Packages
	}

	for _, p := range manifestPackages {
		t, err := dm.ComposeTriggers(filePath, p, ma)
		if err == nil {
			triggers = append(triggers, t...)
		} else {
			return nil, err
		}
	}
	return triggers, nil
}

func (dm *YAMLParser) ComposeTriggers(filePath string, pkg Package, ma whisk.KeyValue) ([]*whisk.Trigger, error) {
	var errorParser error
	var t1 []*whisk.Trigger = make([]*whisk.Trigger, 0)

	for _, trigger := range pkg.GetTriggerList() {
		wsktrigger := new(whisk.Trigger)
		wsktrigger.Name = wskenv.ConvertSingleName(trigger.Name)
		wsktrigger.Namespace = trigger.Namespace
		pub := false
		wsktrigger.Publish = &pub

		// print warning information when .Source key's value is not empty
		if trigger.Source != "" {
			warningString := wski18n.T(
				wski18n.ID_WARN_KEY_DEPRECATED_X_oldkey_X_filetype_X_newkey_X,
				map[string]interface{}{
					wski18n.KEY_OLD:       YAML_KEY_SOURCE,
					wski18n.KEY_NEW:       YAML_KEY_FEED,
					wski18n.KEY_FILE_TYPE: wski18n.MANIFEST})
			wskprint.PrintOpenWhiskWarning(warningString)
		}
		if trigger.Feed == "" {
			trigger.Feed = trigger.Source
		}

		// replacing env. variables here in the trigger feed name
		// to support trigger feed with $READ_FROM_ENV_TRIGGER_FEED
		trigger.Feed = wskenv.InterpolateStringWithEnvVar(trigger.Feed).(string)

		keyValArr := make(whisk.KeyValueArr, 0)
		if trigger.Feed != "" {
			var keyVal whisk.KeyValue

			keyVal.Key = YAML_KEY_FEED
			keyVal.Value = trigger.Feed

			keyValArr = append(keyValArr, keyVal)

			wsktrigger.Annotations = keyValArr
		}

		keyValArr = make(whisk.KeyValueArr, 0)
		for name, param := range trigger.Inputs {
			var keyVal whisk.KeyValue
			keyVal.Key = name

			keyVal.Value, errorParser = ResolveParameter(name, &param, filePath)

			if errorParser != nil {
				return nil, errorParser
			}

			if keyVal.Value != nil {
				keyValArr = append(keyValArr, keyVal)
			}
		}

		if len(keyValArr) > 0 {
			wsktrigger.Parameters = keyValArr
		}

		listOfAnnotations := make(whisk.KeyValueArr, 0)
		for name, value := range trigger.Annotations {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value = wskenv.InterpolateStringWithEnvVar(value)
			listOfAnnotations = append(listOfAnnotations, keyVal)
		}
		if len(listOfAnnotations) > 0 {
			wsktrigger.Annotations = append(wsktrigger.Annotations, listOfAnnotations...)
		}

		// add managed annotations if its a managed deployment
		if utils.Flags.Managed {
			wsktrigger.Annotations = append(wsktrigger.Annotations, ma)
		}

		t1 = append(t1, wsktrigger)
	}
	return t1, nil
}

func (dm *YAMLParser) ComposeRulesFromAllPackages(manifest *YAML, ma whisk.KeyValue) ([]*whisk.Rule, error) {
	var rules []*whisk.Rule = make([]*whisk.Rule, 0)
	manifestPackages := make(map[string]Package)

	if len(manifest.Packages) != 0 {
		manifestPackages = manifest.Packages
	} else {
		manifestPackages = manifest.GetProject().Packages
	}

	for n, p := range manifestPackages {
		r, err := dm.ComposeRules(p, n, ma)
		if err == nil {
			rules = append(rules, r...)
		} else {
			return nil, err
		}
	}
	return rules, nil
}

func (dm *YAMLParser) ComposeRules(pkg Package, packageName string, ma whisk.KeyValue) ([]*whisk.Rule, error) {
	var r1 []*whisk.Rule = make([]*whisk.Rule, 0)

	for _, rule := range pkg.GetRuleList() {
		wskrule := new(whisk.Rule)
		wskrule.Name = wskenv.ConvertSingleName(rule.Name)
		//wskrule.Namespace = rule.Namespace
		pub := false
		wskrule.Publish = &pub
		wskrule.Trigger = wskenv.ConvertSingleName(rule.Trigger)
		wskrule.Action = wskenv.ConvertSingleName(rule.Action)
		act := strings.TrimSpace(wskrule.Action.(string))
		if !strings.ContainsRune(act, '/') && !strings.HasPrefix(act, packageName+"/") &&
			strings.ToLower(packageName) != DEFAULT_PACKAGE {
			act = path.Join(packageName, act)
		}
		wskrule.Action = act
		listOfAnnotations := make(whisk.KeyValueArr, 0)
		for name, value := range rule.Annotations {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value = wskenv.InterpolateStringWithEnvVar(value)
			listOfAnnotations = append(listOfAnnotations, keyVal)
		}
		if len(listOfAnnotations) > 0 {
			wskrule.Annotations = append(wskrule.Annotations, listOfAnnotations...)
		}

		// add managed annotations if its a managed deployment
		if utils.Flags.Managed {
			wskrule.Annotations = append(wskrule.Annotations, ma)
		}

		r1 = append(r1, wskrule)
	}
	return r1, nil
}

func (dm *YAMLParser) ComposeApiRecordsFromAllPackages(client *whisk.Config, manifest *YAML) ([]*whisk.ApiCreateRequest, error) {
	var requests []*whisk.ApiCreateRequest = make([]*whisk.ApiCreateRequest, 0)
	manifestPackages := make(map[string]Package)

	//if manifest.Package.Packagename != "" {
	//	return dm.ComposeApiRecords(client, manifest.Package.Packagename, manifest.Package, manifest.Filepath)
	//} else {
	if len(manifest.Packages) != 0 {
		manifestPackages = manifest.Packages
	} else {
		manifestPackages = manifest.GetProject().Packages
	}
	//}

	for packageName, p := range manifestPackages {
		r, err := dm.ComposeApiRecords(client, packageName, p, manifest.Filepath)
		if err == nil {
			requests = append(requests, r...)
		} else {
			return nil, err
		}
	}
	return requests, nil
}

/*
 * read API section from manifest file:
 * apis: # List of APIs
 *     hello-world: #API name
 *	/hello: #gateway base path
 *	    /world:   #gateway rel path
 *		greeting: get #action name: gateway method
 *
 * compose APIDoc structure from the manifest:
 * {
 *	"apidoc":{
 *      	"namespace":<namespace>,
 *      	"gatewayBasePath":"/hello",
 *      	"gatewayPath":"/world",
 *      	"gatewayMethod":"GET",
 *      	"action":{
 *         		"name":"hello",
 *			"namespace":"guest",
 *			"backendMethod":"GET",
 *			"backendUrl":<url>,
 *			"authkey":<auth>
 *		}
 * 	}
 * }
 */
func (dm *YAMLParser) ComposeApiRecords(client *whisk.Config, packageName string, pkg Package, manifestPath string) ([]*whisk.ApiCreateRequest, error) {
	var requests []*whisk.ApiCreateRequest = make([]*whisk.ApiCreateRequest, 0)

	if pkg.Apis != nil {
		// verify APIGW_ACCESS_TOKEN is set before composing APIs
		// until this point, we dont know whether APIs are specified in manifest or not
		if len(client.ApigwAccessToken) == 0 {
			return nil, wskderrors.NewWhiskClientInvalidConfigError(wski18n.ID_MSG_CONFIG_MISSING_APIGW_ACCESS_TOKEN)
		}
	}

	for apiName, apiDoc := range pkg.Apis {
		for gatewayBasePath, gatewayBasePathMap := range apiDoc {
			// append "/" to the gateway base path if its missing
			if !strings.HasPrefix(gatewayBasePath, PATH_SEPERATOR) {
				gatewayBasePath = PATH_SEPERATOR + gatewayBasePath
			}
			for gatewayRelPath, gatewayRelPathMap := range gatewayBasePathMap {
				// append "/" to the gateway relative path if its missing
				if !strings.HasPrefix(gatewayRelPath, PATH_SEPERATOR) {
					gatewayRelPath = PATH_SEPERATOR + gatewayRelPath
				}
				for actionName, gatewayMethod := range gatewayRelPathMap {
					// verify that the action is defined under actions sections
					if _, ok := pkg.Actions[actionName]; !ok {
						return nil, wskderrors.NewYAMLFileFormatError(manifestPath,
							wski18n.T(wski18n.ID_ERR_API_MISSING_ACTION_X_action_X_api_X,
								map[string]interface{}{
									wski18n.KEY_ACTION: actionName,
									wski18n.KEY_API:    apiName}))
					} else {
						// verify that the action is defined as web action
						// web-export set to any of [true, yes, raw]
						if !utils.IsWebAction(pkg.Actions[actionName].Webexport) {
							return nil, wskderrors.NewYAMLFileFormatError(manifestPath,
								wski18n.T(wski18n.ID_ERR_API_MISSING_WEB_ACTION_X_action_X_api_X,
									map[string]interface{}{
										wski18n.KEY_ACTION: actionName,
										wski18n.KEY_API:    apiName}))
						} else {
							request := new(whisk.ApiCreateRequest)
							request.ApiDoc = new(whisk.Api)
							request.ApiDoc.GatewayBasePath = gatewayBasePath
							// is API verb is valid, it must be one of (GET, PUT, POST, DELETE)
							request.ApiDoc.GatewayRelPath = gatewayRelPath
							if _, ok := whisk.ApiVerbs[strings.ToUpper(gatewayMethod)]; !ok {
								return nil, wskderrors.NewInvalidAPIGatewayMethodError(manifestPath,
									gatewayBasePath+gatewayRelPath,
									gatewayMethod,
									dm.getGatewayMethods())
							}
							request.ApiDoc.GatewayMethod = strings.ToUpper(gatewayMethod)
							request.ApiDoc.Namespace = client.Namespace
							request.ApiDoc.ApiName = apiName
							request.ApiDoc.Id = strings.Join([]string{API, request.ApiDoc.Namespace, request.ApiDoc.GatewayRelPath}, ":")
							// set action of an API Doc
							request.ApiDoc.Action = new(whisk.ApiAction)
							if packageName == DEFAULT_PACKAGE {
								request.ApiDoc.Action.Name = actionName
							} else {
								request.ApiDoc.Action.Name = packageName + PATH_SEPERATOR + actionName
							}
							url := []string{HTTPS + ":" + PATH_SEPERATOR, client.Host, strings.ToLower(API),
								API_VERSION, WEB, client.Namespace, packageName, actionName + "." + HTTP}
							request.ApiDoc.Action.Namespace = client.Namespace
							request.ApiDoc.Action.BackendUrl = strings.Join(url, PATH_SEPERATOR)
							request.ApiDoc.Action.BackendMethod = gatewayMethod
							request.ApiDoc.Action.Auth = client.AuthToken
							// add a newly created ApiCreateRequest object to a list of requests
							requests = append(requests, request)
						}
					}
				}
			}
		}
	}
	return requests, nil
}

func (dm *YAMLParser) getGatewayMethods() []string {
	methods := []string{}
	for k := range whisk.ApiVerbs {
		methods = append(methods, k)
	}
	return methods
}
