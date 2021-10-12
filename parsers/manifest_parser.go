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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"net/url"

	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/apache/openwhisk-wskdeploy/conductor"
	"github.com/apache/openwhisk-wskdeploy/dependencies"
	"github.com/apache/openwhisk-wskdeploy/runtimes"
	"github.com/apache/openwhisk-wskdeploy/utils"
	"github.com/apache/openwhisk-wskdeploy/webaction"
	"github.com/apache/openwhisk-wskdeploy/wskderrors"
	"github.com/apache/openwhisk-wskdeploy/wskenv"
	"github.com/apache/openwhisk-wskdeploy/wski18n"
	"github.com/apache/openwhisk-wskdeploy/wskprint"
	yamlHelper "github.com/ghodss/yaml"
)

const (
	API                   = "API"
	HTTPS                 = "https://"
	HTTP                  = "http://"
	API_VERSION           = "v1"
	WEB                   = "web"
	PATH_SEPARATOR        = "/"
	DEFAULT_PACKAGE       = "default"
	NATIVE_DOCKER_IMAGE   = "openwhisk/dockerskeleton"
	PARAM_OPENING_BRACKET = "{"
	PARAM_CLOSING_BRACKET = "}"

	DUMMY_APIGW_ACCESS_TOKEN = "DUMMY TOKEN"

	YAML_FILE_EXTENSION = "yaml"
	YML_FILE_EXTENSION  = "yml"
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

func (dm *YAMLParser) composeInputs(inputs map[string]Parameter, packageInputs PackageInputs, manifestFilePath string) (whisk.KeyValueArr, error) {
	var errorParser error
	keyValArr := make(whisk.KeyValueArr, 0)
	var inputsWithoutValue []string
	var paramsCLI interface{}

	if len(utils.Flags.Param) > 0 {
		paramsCLI, errorParser = utils.GetJSONFromStrings(utils.Flags.Param, false)
		if errorParser != nil {
			return nil, errorParser
		}
	}

	for name, param := range inputs {
		var keyVal whisk.KeyValue
		// keyvalue key is set to parameter name
		keyVal.Key = name
		// parameter on CLI takes the highest precedence such that
		// input variables gets values from CLI first
		if paramsCLI != nil {
			// check if this particular input is specified on CLI
			if v, ok := paramsCLI.(map[string]interface{})[name]; ok {
				keyVal.Value = wskenv.ConvertSingleName(v.(string))
			}
		}
		// if those inputs are not specified on CLI,
		// read their values from the manifest file
		if keyVal.Value == nil {
			keyVal.Value, errorParser = ResolveParameter(name, &param, manifestFilePath)
			if errorParser != nil {
				return nil, errorParser
			}
			if param.Type == STRING && param.Value != nil {
				if keyVal.Value == getTypeDefaultValue(param.Type) {
					if packageInputs.Inputs != nil {
						n := wskenv.GetEnvVarName(param.Value.(string))
						if v, ok := packageInputs.Inputs[n]; ok {
							keyVal.Value = v.Value.(string)
						}
					}
				}
			}
		}
		if keyVal.Value != nil {
			keyValArr = append(keyValArr, keyVal)
		}
		if param.Required && keyVal.Value == getTypeDefaultValue(param.Type) {
			inputsWithoutValue = append(inputsWithoutValue, name)
		}
	}
	if len(inputsWithoutValue) > 0 {
		errMessage := wski18n.T(wski18n.ID_ERR_REQUIRED_INPUTS_MISSING_VALUE_X_inputs_X,
			map[string]interface{}{
				wski18n.KEY_INPUTS: strings.Join(inputsWithoutValue, ", ")})
		return nil, wskderrors.NewYAMLFileFormatError(manifestFilePath, errMessage)
	}

	return keyValArr, nil
}

func (dm *YAMLParser) composeAnnotations(annotations map[string]interface{}) whisk.KeyValueArr {
	listOfAnnotations := make(whisk.KeyValueArr, 0)
	for name, value := range annotations {
		var keyVal whisk.KeyValue
		keyVal.Key = name
		value = wskenv.InterpolateStringWithEnvVar(value)
		keyVal.Value = utils.ConvertInterfaceValue(value)
		listOfAnnotations = append(listOfAnnotations, keyVal)
	}
	return listOfAnnotations
}

func (dm *YAMLParser) ComposeDependenciesFromAllPackages(manifest *YAML, projectPath string, filePath string, managedAnnotations whisk.KeyValue, packageInputs map[string]PackageInputs) (map[string]dependencies.DependencyRecord, error) {
	dependencies := make(map[string]dependencies.DependencyRecord)
	packages := make(map[string]Package)

	if len(manifest.Packages) != 0 {
		packages = manifest.Packages
	} else {
		packages = manifest.GetProject().Packages
	}

	for n, p := range packages {
		d, err := dm.ComposeDependencies(p, projectPath, filePath, n, managedAnnotations, packageInputs[n])
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

func (dm *YAMLParser) ComposeDependencies(pkg Package, projectPath string, filePath string, packageName string, managedAnnotations whisk.KeyValue, packageInputs PackageInputs) (map[string]dependencies.DependencyRecord, error) {

	depMap := make(map[string]dependencies.DependencyRecord)
	for key, dependency := range pkg.Dependencies {
		version := dependency.Version
		if len(version) == 0 {
			version = YAML_VALUE_BRANCH_MASTER
		}

		location := dependency.Location

		isBinding := false
		if dependencies.LocationIsBinding(location) {
			if !strings.HasPrefix(location, PATH_SEPARATOR) {
				location = PATH_SEPARATOR + dependency.Location
			}
			isBinding = true
		} else if dependencies.LocationIsGithub(location) {

			// TODO() define const for the protocol prefix, etc.
			_, err := url.ParseRequestURI(location)
			if err != nil {
				location = HTTPS + dependency.Location
				location = wskenv.InterpolateStringWithEnvVar(location).(string)
			}
			isBinding = false
		} else {
			// TODO() create new named error in wskerrors package
			return nil, errors.New(wski18n.T(wski18n.ID_ERR_DEPENDENCY_UNKNOWN_TYPE))
		}

		inputs, err := dm.composeInputs(dependency.Inputs, packageInputs, filePath)
		if err != nil {
			return nil, err
		}

		annotations := dm.composeAnnotations(dependency.Annotations)

		if utils.Flags.Managed || utils.Flags.Sync {
			annotations = append(annotations, managedAnnotations)
		}

		packDir := path.Join(projectPath, strings.Title(YAML_KEY_PACKAGES))
		depName := packageName + ":" + key
		depMap[depName] = dependencies.NewDependencyRecord(packDir, packageName, location, version, inputs, annotations, isBinding)
	}

	return depMap, nil
}

func (dm *YAMLParser) ComposeAllPackages(projectInputs map[string]Parameter, manifest *YAML, filePath string, managedAnnotations whisk.KeyValue) (map[string]*whisk.Package, map[string]PackageInputs, error) {
	packages := map[string]*whisk.Package{}
	manifestPackages := make(map[string]Package)
	inputs := make(map[string]PackageInputs, 0)

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
		s, params, err := dm.ComposePackage(p, n, filePath, managedAnnotations, projectInputs)
		if err != nil {
			return nil, inputs, err
		}
		packages[n] = s
		inputs[n] = PackageInputs{PackageName: n, Inputs: params}
	}

	return packages, inputs, nil
}

func (dm *YAMLParser) composePackageInputs(projectInputs map[string]Parameter, rawInputs map[string]Parameter, filepath string) (map[string]Parameter, whisk.KeyValueArr, error) {
	inputs := make(map[string]Parameter, 0)

	// package inherits all project inputs
	for n, param := range projectInputs {
		inputs[n] = param
	}

	// iterate over package inputs
	for name, i := range rawInputs {
		value, err := ResolveParameter(name, &i, filepath)
		if err != nil {
			return nil, nil, err
		}
		// if value is set to default value for its type,
		// check for input key being an env. variable itself
		if value == getTypeDefaultValue(i.Type) {
			value = wskenv.InterpolateStringWithEnvVar("${" + name + "}")
		}

		// if at this point, still value is set to default value of its type
		// check if input key is defined under Project Inputs
		if value == getTypeDefaultValue(i.Type) {
			if i.Type == STRING && i.Value != nil {
				n := wskenv.GetEnvVarName(i.Value.(string))
				if v, ok := projectInputs[n]; ok {
					value = v.Value.(string)
				}
			}
		}

		// create a Parameter object based on the package inputs
		// resolve the value using env. variables
		// if value is not specified, treat input key as an env. variable
		// check if input key is defined in environment
		// else set it to its default value based on the type for now
		// the input value will be updated if its specified in deployment
		// or on CLI using --param or --param-file
		i.Value = value
		inputs[name] = i
	}

	// create an array of Key/Value pair with inputs
	// inputs name as key and with its value if its not nil
	keyValArr := make(whisk.KeyValueArr, 0)
	for name, param := range inputs {
		var keyVal whisk.KeyValue
		keyVal.Key = name
		keyVal.Value = param.Value
		if keyVal.Value != nil {
			keyValArr = append(keyValArr, keyVal)
		}
	}

	return inputs, keyValArr, nil
}

func (dm *YAMLParser) ComposePackage(pkg Package, packageName string, filePath string, managedAnnotations whisk.KeyValue, projectInputs map[string]Parameter) (*whisk.Package, map[string]Parameter, error) {
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
	pag.Version = wskenv.ConvertSingleName(pkg.Version)

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

	// package inputs are set as package inputs of type Parameter{}
	// read all package inputs, interpolate their values using env. variables
	// check if input variable itself is an env. variable
	packageInputs, inputs, err := dm.composePackageInputs(projectInputs, pkg.Inputs, filePath)
	if err != nil {
		return nil, nil, err
	}
	if len(inputs) > 0 {
		pag.Parameters = inputs
	}

	// set Package Annotations
	listOfAnnotations := dm.composeAnnotations(pkg.Annotations)
	if len(listOfAnnotations) > 0 {
		pag.Annotations = append(pag.Annotations, listOfAnnotations...)
	}

	// add Managed Annotations if this is Managed Deployment
	if utils.Flags.Managed || utils.Flags.Sync {
		pag.Annotations = append(pag.Annotations, managedAnnotations)
	}

	// "default" package is a reserved package name
	// and in this case wskdeploy deploys openwhisk entities under
	// /namespace instead of /namespace/package
	if strings.ToLower(pag.Name) == DEFAULT_PACKAGE {
		wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, wski18n.T(wski18n.ID_MSG_DEFAULT_PACKAGE))
		// when a package is marked public with "Public: true" in manifest file
		// the package is visible to anyone and created with publish state
		// set to shared otherwise publish state is set to private
	} else if pkg.Public {
		warningMsg := wski18n.T(wski18n.ID_WARN_PACKAGE_IS_PUBLIC_X_package_X,
			map[string]interface{}{
				wski18n.KEY_PACKAGE: pag.Name})
		wskprint.PrintlnOpenWhiskWarning(warningMsg)
		pag.Publish = &(pkg.Public)
	}

	return pag, packageInputs, nil
}

func (dm *YAMLParser) ComposeSequencesFromAllPackages(namespace string, mani *YAML, manifestFilePath string, managedAnnotations whisk.KeyValue, packageInputs map[string]PackageInputs) ([]utils.ActionRecord, error) {
	var sequences []utils.ActionRecord = make([]utils.ActionRecord, 0)
	manifestPackages := make(map[string]Package)

	if len(mani.Packages) != 0 {
		manifestPackages = mani.Packages
	} else {
		manifestPackages = mani.GetProject().Packages
	}

	for n, p := range manifestPackages {
		s, err := dm.ComposeSequences(namespace, p.Sequences, n, manifestFilePath, managedAnnotations, packageInputs[n])
		if err == nil {
			sequences = append(sequences, s...)
		} else {
			return nil, err
		}
	}
	return sequences, nil
}

func (dm *YAMLParser) ComposeSequences(namespace string, sequences map[string]Sequence, packageName string, manifestFilePath string, managedAnnotations whisk.KeyValue, packageInputs PackageInputs) ([]utils.ActionRecord, error) {
	var listOfSequences []utils.ActionRecord = make([]utils.ActionRecord, 0)
	var errorParser error

	for key, sequence := range sequences {
		wskaction := new(whisk.Action)
		wskaction.Exec = new(whisk.Exec)
		wskaction.Exec.Kind = YAML_KEY_SEQUENCE
		actionList := strings.Split(sequence.Actions, ",")

		var components []string
		for _, a := range actionList {
			act := strings.TrimSpace(a)
			if !strings.ContainsRune(act, []rune(PATH_SEPARATOR)[0]) && !strings.HasPrefix(act, packageName+PATH_SEPARATOR) &&
				strings.ToLower(packageName) != DEFAULT_PACKAGE {
				act = path.Join(packageName, act)
			}
			components = append(components, path.Join(PATH_SEPARATOR+namespace, act))
		}

		wskaction.Exec.Components = components
		if i, ok := packageInputs.Inputs[wskenv.GetEnvVarName(key)]; ok {
			wskaction.Name = i.Value.(string)
		} else {
			wskaction.Name = wskenv.ConvertSingleName(key)
		}
		pub := false
		wskaction.Publish = &pub
		wskaction.Namespace = namespace

		annotations := dm.composeAnnotations(sequence.Annotations)
		if len(annotations) > 0 {
			wskaction.Annotations = annotations
		}

		// appending managed annotations if its a managed deployment
		if utils.Flags.Managed || utils.Flags.Sync {
			wskaction.Annotations = append(wskaction.Annotations, managedAnnotations)
		}

		// Web Export
		// Treat sequence as a web action, a raw HTTP web action, or as a standard action based on web-export;
		// when web-export is set to yes | true, treat sequence as a web action,
		// when web-export is set to raw, treat sequence as a raw HTTP web action,
		// when web-export is set to no | false, treat sequence as a standard action
		if len(sequence.Web) != 0 {
			wskaction.Annotations, errorParser = webaction.SetWebActionAnnotations(manifestFilePath, wskaction.Name, sequence.Web, wskaction.Annotations, false)
			if errorParser != nil {
				return nil, errorParser
			}
		}

		// TODO Add web-secure support for sequences?

		record := utils.ActionRecord{Action: wskaction, Packagename: packageName, Filepath: key}
		listOfSequences = append(listOfSequences, record)
	}
	return listOfSequences, nil
}

func (dm *YAMLParser) ComposeActionsFromAllPackages(manifest *YAML, filePath string, managedAnnotations whisk.KeyValue, packageInputs map[string]PackageInputs) ([]utils.ActionRecord, error) {
	var actions []utils.ActionRecord = make([]utils.ActionRecord, 0)
	manifestPackages := make(map[string]Package)

	if len(manifest.Packages) != 0 {
		manifestPackages = manifest.Packages
	} else {
		manifestPackages = manifest.GetProject().Packages
	}

	for n, p := range manifestPackages {
		a, err := dm.ComposeActions(filePath, p.Actions, n, managedAnnotations, packageInputs[n])
		if err == nil {
			actions = append(actions, a...)
		} else {
			return nil, err
		}
	}
	return actions, nil
}

func (dm *YAMLParser) validateActionCode(manifestFilePath string, action Action) error {
	// Check if action.Function is specified with action.Code
	// with action.Code, action.Function is not allowed
	// with action.Code, action.Runtime should be specified
	if len(action.Function) != 0 {
		err := wski18n.T(wski18n.ID_ERR_ACTION_INVALID_X_action_X,
			map[string]interface{}{
				wski18n.KEY_ACTION: action.Name})
		return wskderrors.NewYAMLFileFormatError(manifestFilePath, err)
	}
	if len(action.Runtime) == 0 {
		err := wski18n.T(wski18n.ID_ERR_ACTION_MISSING_RUNTIME_WITH_CODE_X_action_X,
			map[string]interface{}{
				wski18n.KEY_ACTION: action.Name})
		return wskderrors.NewYAMLFileFormatError(manifestFilePath, err)
	}
	return nil
}

func (dm *YAMLParser) readActionCode(manifestFilePath string, action Action) (*whisk.Exec, error) {
	exec := new(whisk.Exec)
	if err := dm.validateActionCode(manifestFilePath, action); err != nil {
		return nil, err
	}
	// validate runtime from the manifest file
	// error out if the specified runtime is not valid or not supported
	// even if runtime is invalid, deploy action with specified runtime in strict mode
	if utils.Flags.Strict {
		exec.Kind = action.Runtime
	} else if runtimes.CheckExistRuntime(action.Runtime, runtimes.SupportedRunTimes) {
		exec.Kind = action.Runtime
	} else if len(runtimes.DefaultRunTimes[action.Runtime]) != 0 {
		exec.Kind = runtimes.DefaultRunTimes[action.Runtime]
	} else {
		err := wski18n.T(wski18n.ID_ERR_RUNTIME_INVALID_X_runtime_X_action_X,
			map[string]interface{}{
				wski18n.KEY_RUNTIME: action.Runtime,
				wski18n.KEY_ACTION:  action.Name})
		return nil, wskderrors.NewYAMLFileFormatError(manifestFilePath, err)
	}
	exec.Code = &(action.Code)
	// we can specify the name of the action entry point using main
	if len(action.Main) != 0 {
		exec.Main = action.Main
	}
	return exec, nil
}

func (dm *YAMLParser) validateActionFunction(manifestFileName string, action Action, ext string, kind string) error {
	// produce an error when a runtime could not be derived from the action file extension
	// and its not explicitly specified in the manifest YAML file
	// and action source is not a zip file
	if len(action.Runtime) == 0 && len(action.Docker) == 0 && !action.Native {
		if ext == runtimes.ZIP_FILE_EXTENSION {
			errMessage := wski18n.T(wski18n.ID_ERR_RUNTIME_INVALID_X_runtime_X_action_X,
				map[string]interface{}{
					wski18n.KEY_RUNTIME: runtimes.RUNTIME_NOT_SPECIFIED,
					wski18n.KEY_ACTION:  action.Name})
			return wskderrors.NewInvalidRuntimeError(errMessage,
				manifestFileName,
				action.Name,
				runtimes.RUNTIME_NOT_SPECIFIED,
				runtimes.ListOfSupportedRuntimes(runtimes.SupportedRunTimes))
		} else if len(kind) == 0 {
			errMessage := wski18n.T(wski18n.ID_ERR_RUNTIME_ACTION_SOURCE_NOT_SUPPORTED_X_ext_X_action_X,
				map[string]interface{}{
					wski18n.KEY_EXTENSION: ext,
					wski18n.KEY_ACTION:    action.Name})
			return wskderrors.NewInvalidRuntimeError(errMessage,
				manifestFileName,
				action.Name,
				runtimes.RUNTIME_NOT_SPECIFIED,
				runtimes.ListOfSupportedRuntimes(runtimes.SupportedRunTimes))
		}
	}
	return nil
}

func (dm *YAMLParser) readActionFunction(manifestFilePath string, manifestFileName string, action Action) (string, *whisk.Exec, error) {
	var actionFilePath string
	var zipFileName string
	exec := new(whisk.Exec)
	f := wskenv.InterpolateStringWithEnvVar(action.Function)
	interpolatedActionFunction := f.(string)

	// check if action function is pointing to an URL
	// we do not support if function is pointing to remote directory
	// therefore error out if there is a combination of http/https ending in a directory
	if strings.HasPrefix(interpolatedActionFunction, HTTP) || strings.HasPrefix(interpolatedActionFunction, HTTPS) {
		if len(path.Ext(interpolatedActionFunction)) == 0 {
			err := wski18n.T(wski18n.ID_ERR_ACTION_FUNCTION_REMOTE_DIR_NOT_SUPPORTED_X_action_X_url_X,
				map[string]interface{}{
					wski18n.KEY_ACTION: action.Name,
					wski18n.KEY_URL:    interpolatedActionFunction})
			return actionFilePath, nil, wskderrors.NewYAMLFileFormatError(manifestFilePath, err)
		}
		actionFilePath = interpolatedActionFunction
	} else {
		actionFilePath = strings.TrimRight(manifestFilePath, manifestFileName) + interpolatedActionFunction
	}

	if utils.IsDirectory(actionFilePath) {
		zipFileName = actionFilePath + "." + runtimes.ZIP_FILE_EXTENSION
		err := utils.NewZipWriter(actionFilePath, zipFileName, action.Include, action.Exclude, filepath.Dir(manifestFilePath)).Zip()
		if err != nil {
			return actionFilePath, nil, err
		}
		defer os.Remove(zipFileName)
		actionFilePath = zipFileName
	}

	action.Function = actionFilePath

	// determine extension of the given action source file
	ext := filepath.Ext(actionFilePath)
	// drop the "." from file extension
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	// determine default runtime for the given file extension
	var kind string
	r := runtimes.FileExtensionRuntimeKindMap[ext]
	kind = runtimes.DefaultRunTimes[r]
	if err := dm.validateActionFunction(manifestFileName, action, ext, kind); err != nil {
		return actionFilePath, nil, err
	}
	exec.Kind = kind

	dat, err := utils.Read(actionFilePath)
	if err != nil {
		return actionFilePath, nil, err
	}
	code := string(dat)
	if ext == runtimes.ZIP_FILE_EXTENSION || ext == runtimes.JAR_FILE_EXTENSION {
		code = base64.StdEncoding.EncodeToString([]byte(dat))
	}
	exec.Code = &code

	/*
	*  Action.Runtime
	*  Perform few checks if action runtime is specified in manifest YAML file
	*  (1) Check if specified runtime is one of the supported runtimes by OpenWhisk server
	*  (2) Check if specified runtime is consistent with action source file extensions
	*  Set the action runtime to match with the source file extension, if wskdeploy is not invoked in strict mode
	 */
	if len(action.Runtime) != 0 {
		if runtimes.CheckExistRuntime(action.Runtime, runtimes.SupportedRunTimes) {
			// for zip actions, rely on the runtimes from the manifest file as it can not be derived from the action source file extension
			// pick runtime from manifest file if its supported by OpenWhisk server
			if ext == runtimes.ZIP_FILE_EXTENSION {
				exec.Kind = action.Runtime
			} else {
				if runtimes.CheckRuntimeConsistencyWithFileExtension(ext, action.Runtime) {
					exec.Kind = action.Runtime
				} else {
					warnStr := wski18n.T(wski18n.ID_ERR_RUNTIME_MISMATCH_X_runtime_X_ext_X_action_X,
						map[string]interface{}{
							wski18n.KEY_RUNTIME:   action.Runtime,
							wski18n.KEY_EXTENSION: ext,
							wski18n.KEY_ACTION:    action.Name})
					wskprint.PrintOpenWhiskWarning(warnStr)

					// even if runtime is not consistent with file extension, deploy action with specified runtime in strict mode
					if utils.Flags.Strict {
						exec.Kind = action.Runtime
					} else {
						warnStr := wski18n.T(wski18n.ID_WARN_RUNTIME_CHANGED_X_runtime_X_action_X,
							map[string]interface{}{
								wski18n.KEY_RUNTIME: exec.Kind,
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

			if ext == runtimes.ZIP_FILE_EXTENSION {
				// for zip action, error out if specified runtime is not supported by
				// OpenWhisk server
				return actionFilePath, nil, wskderrors.NewInvalidRuntimeError(warnStr,
					manifestFileName,
					action.Name,
					action.Runtime,
					runtimes.ListOfSupportedRuntimes(runtimes.SupportedRunTimes))
			} else {
				if utils.Flags.Strict {
					exec.Kind = action.Runtime
				} else {
					warnStr := wski18n.T(wski18n.ID_WARN_RUNTIME_CHANGED_X_runtime_X_action_X,
						map[string]interface{}{
							wski18n.KEY_RUNTIME: exec.Kind,
							wski18n.KEY_ACTION:  action.Name})
					wskprint.PrintOpenWhiskWarning(warnStr)
				}
			}

		}
	}
	// we can specify the name of the action entry point using main
	if len(action.Main) != 0 {
		exec.Main = action.Main
	}

	return actionFilePath, exec, nil
}

func (dm *YAMLParser) composeActionExec(manifestFilePath string, manifestFileName string, action Action) (string, *whisk.Exec, error) {
	var actionFilePath string
	exec := new(whisk.Exec)
	var err error

	if len(action.Code) != 0 {
		exec, err = dm.readActionCode(manifestFilePath, action)
		if err != nil {
			return actionFilePath, nil, err
		}
	}
	if len(action.Function) != 0 {
		actionFilePath, exec, err = dm.readActionFunction(manifestFilePath, manifestFileName, action)
		if err != nil {
			return actionFilePath, nil, err
		}
	}

	// when an action has Docker image specified,
	// set exec.Kind to "blackbox" and
	// set exec.Image to specified image e.g. dockerhub/image
	// when an action Native is set to true,
	// set exec.Image to openwhisk/skeleton
	if len(action.Docker) != 0 || action.Native {
		exec.Kind = runtimes.BLACKBOX
		if action.Native {
			exec.Image = NATIVE_DOCKER_IMAGE
		} else {
			exec.Image = wskenv.InterpolateStringWithEnvVar(action.Docker).(string)
		}
	}

	return actionFilePath, exec, err
}

func (dm *YAMLParser) validateActionLimits(limits Limits) {
	// TODO() use LIMITS_UNSUPPORTED in yamlparser to enumerate through instead of hardcoding
	// emit warning errors if these limits are not nil
	utils.NotSupportLimits(limits.ConcurrentActivations, LIMIT_VALUE_CONCURRENT_ACTIVATIONS)
	utils.NotSupportLimits(limits.UserInvocationRate, LIMIT_VALUE_USER_INVOCATION_RATE)
	utils.NotSupportLimits(limits.CodeSize, LIMIT_VALUE_CODE_SIZE)
	utils.NotSupportLimits(limits.ParameterSize, LIMIT_VALUE_PARAMETER_SIZE)
}

func (dm *YAMLParser) composeActionLimits(limits Limits) *whisk.Limits {
	dm.validateActionLimits(limits)
	wsklimits := new(whisk.Limits)
	for _, t := range LIMITS_SUPPORTED {
		switch t {
		case LIMIT_VALUE_TIMEOUT:
			if utils.LimitsTimeoutValidation(limits.Timeout) {
				wsklimits.Timeout = limits.Timeout
			} else {
				warningString := wski18n.T(wski18n.ID_WARN_LIMIT_IGNORED_X_limit_X,
					map[string]interface{}{wski18n.KEY_LIMIT: LIMIT_VALUE_TIMEOUT})
				wskprint.PrintOpenWhiskWarning(warningString)
			}
		case LIMIT_VALUE_MEMORY_SIZE:
			if utils.LimitsMemoryValidation(limits.Memory) {
				wsklimits.Memory = limits.Memory
			} else {
				warningString := wski18n.T(wski18n.ID_WARN_LIMIT_IGNORED_X_limit_X,
					map[string]interface{}{wski18n.KEY_LIMIT: LIMIT_VALUE_MEMORY_SIZE})
				wskprint.PrintOpenWhiskWarning(warningString)

			}
		case LIMIT_VALUE_LOG_SIZE:
			if utils.LimitsLogsizeValidation(limits.Logsize) {
				wsklimits.Logsize = limits.Logsize
			} else {
				warningString := wski18n.T(wski18n.ID_WARN_LIMIT_IGNORED_X_limit_X,
					map[string]interface{}{wski18n.KEY_LIMIT: LIMIT_VALUE_LOG_SIZE})
				wskprint.PrintOpenWhiskWarning(warningString)

			}
		}
	}
	if wsklimits.Timeout != nil || wsklimits.Memory != nil || wsklimits.Logsize != nil {
		return wsklimits
	}
	return nil
}

func (dm *YAMLParser) warnIfRedundantWebActionFlags(action Action) {
	// Warn user if BOTH web and web-export specified,
	// as they are redundant; defer to "web" flag and its value
	if len(action.Web) != 0 && len(action.WebExport) != 0 {
		warningString := wski18n.T(wski18n.ID_WARN_ACTION_WEB_X_action_X,
			map[string]interface{}{wski18n.KEY_ACTION: action.Name})
		wskprint.PrintOpenWhiskWarning(warningString)
	}
}

func (dm *YAMLParser) ComposeActions(manifestFilePath string, actions map[string]Action, packageName string, managedAnnotations whisk.KeyValue, packageInputs PackageInputs) ([]utils.ActionRecord, error) {

	var errorParser error
	var listOfActions []utils.ActionRecord = make([]utils.ActionRecord, 0)
	splitManifestFilePath := strings.Split(manifestFilePath, string(PATH_SEPARATOR))
	manifestFileName := splitManifestFilePath[len(splitManifestFilePath)-1]

	for actionName, action := range actions {
		var actionFilePath string

		// update the action (of type Action) to set its name
		// here key name is the action name
		if i, ok := packageInputs.Inputs[wskenv.GetEnvVarName(actionName)]; ok {
			action.Name = i.Value.(string)
		} else {
			action.Name = wskenv.ConvertSingleName(actionName)
		}

		// Create action data object from client library
		wskaction := new(whisk.Action)

		//set action.Function to action.Location
		//because Location is deprecated in Action entity
		if len(action.Function) == 0 && len(action.Location) != 0 {
			action.Function = action.Location
		}

		actionFilePath, wskaction.Exec, errorParser = dm.composeActionExec(manifestFilePath, manifestFileName, action)
		if errorParser != nil {
			return nil, errorParser
		}

		// Action.Inputs
		listOfInputs, err := dm.composeInputs(action.Inputs, packageInputs, manifestFilePath)
		if err != nil {
			return nil, err
		}
		if len(listOfInputs) > 0 {
			wskaction.Parameters = listOfInputs
		}

		// Action.Outputs
		// TODO{} add outputs as annotations (work to discuss officially supporting for compositions)
		//listOfOutputs, err := dm.composeOutputs(action.Outputs, manifestFilePath)
		//if err != nil {
		//	return nil, err
		//}
		//if len(listOfOutputs) > 0 {
		//	wskaction.Annotations = listOfOutputs
		//}

		// Action.Annotations
		// ==================
		// WARNING!  Processing of explicit Annotations MUST occur before handling of Action keys, as these
		// keys often need to check for inconsistencies (and raise errors).
		if listOfAnnotations := dm.composeAnnotations(action.Annotations); len(listOfAnnotations) > 0 {
			wskaction.Annotations = append(wskaction.Annotations, listOfAnnotations...)
		}

		// add managed annotations if its marked as managed deployment
		if utils.Flags.Managed || utils.Flags.Sync {
			wskaction.Annotations = append(wskaction.Annotations, managedAnnotations)
		}

		// Web Export (i.e., "web-export" annotation)
		// ==========
		// Treat ACTION as a web action, a raw HTTP web action, or as a standard action based on web-export;
		// when web-export is set to yes | true, treat action as a web action,
		// when web-export is set to raw, treat action as a raw HTTP web action,
		// when web-export is set to no | false, treat action as a standard action
		dm.warnIfRedundantWebActionFlags(action)
		if len(action.GetWeb()) != 0 {
			wskaction.Annotations, errorParser = webaction.SetWebActionAnnotations(
				manifestFilePath,
				action.Name,
				action.GetWeb(),
				wskaction.Annotations,
				false)
			if errorParser != nil {
				return listOfActions, errorParser
			}
		}

		// validate special action annotations such as "require-whisk-auth"
		// TODO: the Manifest parser will validate any declared APIs that ref. this action
		if wskaction.Annotations != nil {
			if webaction.HasAnnotation(&wskaction.Annotations, webaction.REQUIRE_WHISK_AUTH) {
				_, errorParser = webaction.ValidateRequireWhiskAuthAnnotationValue(
					action.Name,
					wskaction.Annotations.GetValue(webaction.REQUIRE_WHISK_AUTH))
			}
			if errorParser != nil {
				return listOfActions, errorParser
			}
		}

		// Action.Limits
		if action.Limits != nil {
			if wsklimits := dm.composeActionLimits(*(action.Limits)); wsklimits != nil {
				wskaction.Limits = wsklimits
			}
		}

		// Conductor Action
		if action.Conductor {
			wskaction.Annotations = append(wskaction.Annotations, conductor.ConductorAction())
		}

		// Set other top-level values for the action (e.g., name, version, publish, etc.)
		wskaction.Name = action.Name
		pub := false
		wskaction.Publish = &pub
		wskaction.Version = wskenv.ConvertSingleName(action.Version)

		// create a "record" of the Action relative to its package and function filepath
		// which will be used to compose the REST API calls
		record := utils.ActionRecord{Action: wskaction, Packagename: packageName, Filepath: actionFilePath}
		listOfActions = append(listOfActions, record)
	}

	return listOfActions, nil
}

func (dm *YAMLParser) ComposeTriggersFromAllPackages(manifest *YAML, filePath string, managedAnnotations whisk.KeyValue, inputs map[string]PackageInputs) ([]*whisk.Trigger, error) {
	var triggers []*whisk.Trigger = make([]*whisk.Trigger, 0)
	manifestPackages := make(map[string]Package)

	if len(manifest.Packages) != 0 {
		manifestPackages = manifest.Packages
	} else {
		manifestPackages = manifest.GetProject().Packages
	}

	for packageName, pkg := range manifestPackages {
		t, err := dm.ComposeTriggers(filePath, pkg, managedAnnotations, inputs[packageName])
		if err == nil {
			triggers = append(triggers, t...)
		} else {
			return nil, err
		}
	}
	return triggers, nil
}

func (dm *YAMLParser) ComposeTriggers(filePath string, pkg Package, managedAnnotations whisk.KeyValue, packageInputs PackageInputs) ([]*whisk.Trigger, error) {
	var errorParser error
	var listOfTriggers []*whisk.Trigger = make([]*whisk.Trigger, 0)

	for _, trigger := range pkg.GetTriggerList() {
		wsktrigger := new(whisk.Trigger)
		if i, ok := packageInputs.Inputs[wskenv.GetEnvVarName(trigger.Name)]; ok {
			wsktrigger.Name = i.Value.(string)
		} else {
			wsktrigger.Name = wskenv.ConvertSingleName(trigger.Name)
		}
		wsktrigger.Namespace = trigger.Namespace
		pub := false
		wsktrigger.Publish = &pub

		// print warning information when .Source key's value is not empty
		if len(trigger.Source) != 0 {
			warningString := wski18n.T(
				wski18n.ID_WARN_KEY_DEPRECATED_X_oldkey_X_filetype_X_newkey_X,
				map[string]interface{}{
					wski18n.KEY_OLD:       YAML_KEY_SOURCE,
					wski18n.KEY_NEW:       YAML_KEY_FEED,
					wski18n.KEY_FILE_TYPE: wski18n.MANIFEST_FILE})
			wskprint.PrintOpenWhiskWarning(warningString)
		}
		if len(trigger.Feed) == 0 {
			trigger.Feed = trigger.Source
		}

		// replacing env. variables here in the trigger feed name
		// to support trigger feed with $READ_FROM_ENV_TRIGGER_FEED
		trigger.Feed = wskenv.ConvertSingleName(trigger.Feed)

		keyValArr := make(whisk.KeyValueArr, 0)
		if len(trigger.Feed) != 0 {
			var keyVal whisk.KeyValue
			keyVal.Key = YAML_KEY_FEED
			keyVal.Value = trigger.Feed
			keyValArr = append(keyValArr, keyVal)
			wsktrigger.Annotations = keyValArr
		}

		inputs, err := dm.composeInputs(trigger.Inputs, packageInputs, filePath)
		if err != nil {
			return nil, errorParser
		}
		if len(inputs) > 0 {
			wsktrigger.Parameters = inputs
		}

		listOfAnnotations := dm.composeAnnotations(trigger.Annotations)
		if len(listOfAnnotations) > 0 {
			wsktrigger.Annotations = append(wsktrigger.Annotations, listOfAnnotations...)
		}

		// add managed annotations if its a managed deployment
		if utils.Flags.Managed || utils.Flags.Sync {
			wsktrigger.Annotations = append(wsktrigger.Annotations, managedAnnotations)
		}

		listOfTriggers = append(listOfTriggers, wsktrigger)
	}
	return listOfTriggers, nil
}

func (dm *YAMLParser) ComposeRulesFromAllPackages(manifest *YAML, managedAnnotations whisk.KeyValue, packageInputs map[string]PackageInputs) ([]*whisk.Rule, error) {
	var rules []*whisk.Rule = make([]*whisk.Rule, 0)
	manifestPackages := make(map[string]Package)

	if len(manifest.Packages) != 0 {
		manifestPackages = manifest.Packages
	} else {
		manifestPackages = manifest.GetProject().Packages
	}

	for n, p := range manifestPackages {
		r, err := dm.ComposeRules(p, n, managedAnnotations, packageInputs[n])
		if err == nil {
			rules = append(rules, r...)
		} else {
			return nil, err
		}
	}
	return rules, nil
}

func (dm *YAMLParser) ComposeRules(pkg Package, packageName string, managedAnnotations whisk.KeyValue, packageInputs PackageInputs) ([]*whisk.Rule, error) {
	var rules []*whisk.Rule = make([]*whisk.Rule, 0)

	for _, rule := range pkg.GetRuleList() {
		wskrule := new(whisk.Rule)
		if i, ok := packageInputs.Inputs[wskenv.GetEnvVarName(rule.Name)]; ok {
			wskrule.Name = i.Value.(string)
		} else {
			wskrule.Name = wskenv.ConvertSingleName(rule.Name)
		}
		//wskrule.Namespace = rule.Namespace
		pub := false
		wskrule.Publish = &pub
		if i, ok := packageInputs.Inputs[wskenv.GetEnvVarName(rule.Trigger)]; ok {
			wskrule.Trigger = i.Value.(string)
		} else {
			wskrule.Trigger = wskenv.ConvertSingleName(rule.Trigger)
		}
		if i, ok := packageInputs.Inputs[wskenv.GetEnvVarName(rule.Action)]; ok {
			wskrule.Action = i.Value.(string)
		} else {
			wskrule.Action = wskenv.ConvertSingleName(rule.Action)
		}
		act := strings.TrimSpace(wskrule.Action.(string))
		if !strings.ContainsRune(act, []rune(PATH_SEPARATOR)[0]) && !strings.HasPrefix(act, packageName+PATH_SEPARATOR) &&
			strings.ToLower(packageName) != DEFAULT_PACKAGE {
			act = path.Join(packageName, act)
		}
		wskrule.Action = act
		listOfAnnotations := dm.composeAnnotations(rule.Annotations)
		if len(listOfAnnotations) > 0 {
			wskrule.Annotations = append(wskrule.Annotations, listOfAnnotations...)
		}

		// add managed annotations if its a managed deployment
		if utils.Flags.Managed || utils.Flags.Sync {
			wskrule.Annotations = append(wskrule.Annotations, managedAnnotations)
		}

		rules = append(rules, wskrule)
	}
	return rules, nil
}

func (dm *YAMLParser) ComposeApiRecordsFromAllPackages(client *whisk.Config, manifest *YAML,
	actionrecords []utils.ActionRecord,
	sequencerecords []utils.ActionRecord) ([]*whisk.ApiCreateRequest, map[string]*whisk.ApiCreateRequestOptions, error) {
	var requests = make([]*whisk.ApiCreateRequest, 0)
	var responses = make(map[string]*whisk.ApiCreateRequestOptions, 0)
	manifestPackages := make(map[string]Package)

	if len(manifest.Packages) != 0 {
		manifestPackages = manifest.Packages
	} else {
		manifestPackages = manifest.GetProject().Packages
	}

	for packageName, p := range manifestPackages {
		r, response, err := dm.ComposeApiRecords(client, packageName, p, manifest.Filepath,
			actionrecords, sequencerecords)
		if err == nil {
			requests = append(requests, r...)
			for k, v := range response {
				responses[k] = v
			}
		} else {
			return nil, nil, err
		}
	}
	return requests, responses, nil
}

/*
 * read API section from manifest file:
 * apis: # List of APIs
 * 		hello-world: #API name
 *			/hello: #gateway base path
 *	    		/world:   #gateway rel path
 *					greeting: get #action name: gateway method
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
func (dm *YAMLParser) ComposeApiRecords(client *whisk.Config, packageName string, pkg Package, manifestPath string,
	actionrecords []utils.ActionRecord, sequencerecords []utils.ActionRecord) ([]*whisk.ApiCreateRequest, map[string]*whisk.ApiCreateRequestOptions, error) {
	var requests = make([]*whisk.ApiCreateRequest, 0)

	// supply a dummy API GW token as it is optional
	if pkg.Apis != nil && len(pkg.Apis) != 0 {
		if len(client.ApigwAccessToken) == 0 {
			warningString := wski18n.T(wski18n.ID_MSG_CONFIG_PROVIDE_DEFAULT_APIGW_ACCESS_TOKEN,
				map[string]interface{}{wski18n.KEY_DUMMY_TOKEN: DUMMY_APIGW_ACCESS_TOKEN})
			wskprint.PrintOpenWhiskWarning(warningString)
			client.ApigwAccessToken = DUMMY_APIGW_ACCESS_TOKEN
		}
	}

	requestOptions := make(map[string]*whisk.ApiCreateRequestOptions, 0)

	for apiName, apiDoc := range pkg.Apis {
		for gatewayBasePath, gatewayBasePathMap := range apiDoc {
			// Base Path
			// validate base path should not have any path parameters
			if !isGatewayBasePathValid(gatewayBasePath) {
				err := wskderrors.NewYAMLParserErr(manifestPath,
					wski18n.T(wski18n.ID_ERR_API_GATEWAY_BASE_PATH_INVALID_X_api_X,
						map[string]interface{}{wski18n.KEY_API_BASE_PATH: gatewayBasePath}))
				return requests, requestOptions, err
			}
			// append "/" to the gateway base path if its missing
			if !strings.HasPrefix(gatewayBasePath, PATH_SEPARATOR) {
				gatewayBasePath = PATH_SEPARATOR + gatewayBasePath
			}
			for gatewayRelPath, gatewayRelPathMap := range gatewayBasePathMap {
				// Relative Path
				// append "/" to the gateway relative path if its missing
				if !strings.HasPrefix(gatewayRelPath, PATH_SEPARATOR) {
					gatewayRelPath = PATH_SEPARATOR + gatewayRelPath
				}
				for actionName, gatewayMethodResponse := range gatewayRelPathMap {
					// verify that the action is defined under action records
					if _, ok := pkg.Actions[actionName]; ok {
						// verify that the action is defined as web action;
						// web or web-export set to any of [true, yes, raw]; if not,
						// we will try to add it (if no strict" flag) and warn user that we did so
						if err := webaction.TryUpdateAPIsActionToWebAction(actionrecords, packageName,
							apiName, actionName, false); err != nil {
							return requests, requestOptions, err
						}
						// verify that the sequence action is defined under sequence records
					} else if _, ok := pkg.Sequences[actionName]; ok {
						// verify that the sequence action is defined as web sequence
						// web or web-export set to any of [true, yes, raw]; if not,
						// we will try to add it (if no strict" flag) and warn user that we did so
						if err := webaction.TryUpdateAPIsActionToWebAction(sequencerecords, packageName,
							apiName, actionName, true); err != nil {
							return requests, requestOptions, err
						}
					} else {
						return nil, nil, wskderrors.NewYAMLFileFormatError(manifestPath,
							wski18n.T(wski18n.ID_ERR_API_MISSING_ACTION_OR_SEQUENCE_X_action_or_sequence_X_api_X,
								map[string]interface{}{
									wski18n.KEY_ACTION: actionName,
									wski18n.KEY_API:    apiName}))
					}

					// get the list of path parameters from relative path
					pathParameters := generatePathParameters(gatewayRelPath)

					// Check if response type is set to http for API using path parameters
					if strings.ToLower(gatewayMethodResponse.Method) != utils.HTTP_FILE_EXTENSION &&
						len(pathParameters) > 0 {
						warningString := wski18n.T(wski18n.ID_WARN_API_INVALID_RESPONSE_TYPE,
							map[string]interface{}{
								wski18n.KEY_API:               apiName,
								wski18n.KEY_API_RELATIVE_PATH: gatewayRelPath,
								wski18n.KEY_RESPONSE:          gatewayMethodResponse.Response})
						wskprint.PrintlnOpenWhiskWarning(warningString)
						gatewayMethodResponse.Response = utils.HTTP_FILE_EXTENSION
					}

					// Check if API verb is valid, it must be one of (GET, PUT, POST, DELETE)
					if _, ok := whisk.ApiVerbs[strings.ToUpper(gatewayMethodResponse.Method)]; !ok {
						return nil, nil, wskderrors.NewInvalidAPIGatewayMethodError(manifestPath,
							gatewayBasePath+gatewayRelPath,
							gatewayMethodResponse.Method,
							dm.getGatewayMethods())
					}

					apiDocActionName := actionName
					if strings.ToLower(packageName) != DEFAULT_PACKAGE {
						apiDocActionName = packageName + PATH_SEPARATOR + actionName
					}

					// set action of an API Doc
					apiDocAction := whisk.ApiAction{
						Name:      apiDocActionName,
						Namespace: client.Namespace,
						BackendUrl: strings.Join([]string{HTTPS +
							client.Host, strings.ToLower(API),
							API_VERSION, WEB, client.Namespace, packageName,
							actionName + "." + utils.HTTP_FILE_EXTENSION},
							PATH_SEPARATOR),
						BackendMethod: gatewayMethodResponse.Method,
						Auth:          client.AuthToken,
					}

					requestApiDoc := whisk.Api{
						GatewayBasePath: gatewayBasePath,
						PathParameters:  pathParameters,
						GatewayRelPath:  gatewayRelPath,
						GatewayMethod:   strings.ToUpper(gatewayMethodResponse.Method),
						Namespace:       client.Namespace,
						ApiName:         apiName,
						Id:              strings.Join([]string{API, client.Namespace, gatewayRelPath}, ":"),
						Action:          &apiDocAction,
					}

					request := whisk.ApiCreateRequest{
						ApiDoc: &requestApiDoc,
					}

					// add a newly created ApiCreateRequest object to a list of requests
					requests = append(requests, &request)

					// Create an instance of ApiCreateRequestOptions
					options := whisk.ApiCreateRequestOptions{
						ResponseType: gatewayMethodResponse.Response,
					}
					apiPath := request.ApiDoc.ApiName + " " + request.ApiDoc.GatewayBasePath +
						request.ApiDoc.GatewayRelPath + " " + request.ApiDoc.GatewayMethod
					requestOptions[apiPath] = &options
				}
			}
		}
	}
	return requests, requestOptions, nil
}

func (dm *YAMLParser) ComposeApiRecordsFromSwagger(client *whisk.Config, manifest *YAML) (*whisk.ApiCreateRequest, *whisk.ApiCreateRequestOptions, error) {
	var api *whisk.Api
	var err error
	var config string = manifest.Project.Config
	// Check to make sure the manifest has a config under project
	if config == "" {
		return nil, nil, nil
	}

	api, err = dm.parseSwaggerApi(config, client.Namespace)
	if err != nil {
		whisk.Debug(whisk.DbgError, "parseSwaggerApi() error: %s\n", err)
		errMsg := wski18n.T("Unable to parse swagger file: {{.err}}", map[string]interface{}{"err": err})
		whiskErr := whisk.MakeWskErrorFromWskError(errors.New(errMsg), err, whisk.EXIT_CODE_ERR_GENERAL,
			whisk.DISPLAY_MSG, whisk.DISPLAY_USAGE)
		return nil, nil, whiskErr
	}

	apiCreateReq := new(whisk.ApiCreateRequest)
	apiCreateReq.ApiDoc = api

	apiCreateReqOptions := new(whisk.ApiCreateRequestOptions)

	return apiCreateReq, apiCreateReqOptions, nil

}

/*
 * Read the swagger config provided by under
 * Project:
 *   config: swagger_filename.[yaml|yml]
 *
 * NOTE: This was lifted almost verbatim from openwhisk-cli/commands/api.go
 *       and as a follow up should probably be moved and updated to live in whiskclient-go
 * NOTE: This does notb verify that actions used in swagger api definition are defined in
 *        the openwhisk server or manifest.
 */
func (dm *YAMLParser) parseSwaggerApi(configfile, namespace string) (*whisk.Api, error) {
	if len(configfile) == 0 {
		whisk.Debug(whisk.DbgError, "No swagger file is specified\n")
		errMsg := wski18n.T("A swagger configuration file was not specified.")
		whiskErr := whisk.MakeWskError(errors.New(errMsg), whisk.EXIT_CODE_ERR_GENERAL,
			whisk.DISPLAY_MSG, whisk.DISPLAY_USAGE)
		return nil, whiskErr
	}

	swagger, err := ioutil.ReadFile(configfile)
	if err != nil {
		whisk.Debug(whisk.DbgError, "readFile(%s) error: %s\n", configfile, err)
		errMsg := wski18n.T("Error reading swagger file '{{.name}}': {{.err}}",
			map[string]interface{}{"name": configfile, "err": err})
		whiskErr := whisk.MakeWskErrorFromWskError(errors.New(errMsg), err, whisk.EXIT_CODE_ERR_GENERAL,
			whisk.DISPLAY_MSG, whisk.DISPLAY_USAGE)
		return nil, whiskErr
	}

	// Check if this swagger is in JSON or YAML format
	isYaml := strings.HasSuffix(configfile, YAML_FILE_EXTENSION) || strings.HasSuffix(configfile, YML_FILE_EXTENSION)
	if isYaml {
		whisk.Debug(whisk.DbgInfo, "Converting YAML formated API configuration into JSON\n")
		jsonbytes, err := yamlHelper.YAMLToJSON([]byte(swagger))
		if err != nil {
			whisk.Debug(whisk.DbgError, "yaml.YAMLToJSON() error: %s\n", err)
			errMsg := wski18n.T("Unable to parse YAML configuration file: {{.err}}", map[string]interface{}{"err": err})
			whiskErr := whisk.MakeWskError(errors.New(errMsg), whisk.EXIT_CODE_ERR_GENERAL,
				whisk.DISPLAY_MSG, whisk.NO_DISPLAY_USAGE)
			return nil, whiskErr
		}
		swagger = jsonbytes
	}

	// Parse the JSON into a swagger object
	swaggerObj := new(whisk.ApiSwagger)
	err = json.Unmarshal([]byte(swagger), swaggerObj)
	if err != nil {
		whisk.Debug(whisk.DbgError, "JSON parse of '%s' error: %s\n", configfile, err)
		errMsg := wski18n.T("Error parsing swagger file '{{.name}}': {{.err}}",
			map[string]interface{}{"name": configfile, "err": err})
		whiskErr := whisk.MakeWskErrorFromWskError(errors.New(errMsg), err, whisk.EXIT_CODE_ERR_GENERAL,
			whisk.DISPLAY_MSG, whisk.DISPLAY_USAGE)
		return nil, whiskErr
	}
	if swaggerObj.BasePath == "" || swaggerObj.SwaggerName == "" || swaggerObj.Info == nil || swaggerObj.Paths == nil {
		whisk.Debug(whisk.DbgError, "Swagger file is invalid.\n")
		errMsg := wski18n.T("Swagger file is invalid (missing basePath, info, paths, or swagger fields)")
		whiskErr := whisk.MakeWskError(errors.New(errMsg), whisk.EXIT_CODE_ERR_GENERAL,
			whisk.DISPLAY_MSG, whisk.DISPLAY_USAGE)
		return nil, whiskErr
	}

	if !strings.HasPrefix(swaggerObj.BasePath, "/") {
		whisk.Debug(whisk.DbgError, "Swagger file basePath is invalid.\n")
		errMsg := wski18n.T("Swagger file basePath must start with a leading slash (/)")
		whiskErr := whisk.MakeWskError(errors.New(errMsg), whisk.EXIT_CODE_ERR_GENERAL,
			whisk.DISPLAY_MSG, whisk.DISPLAY_USAGE)
		return nil, whiskErr
	}

	api := new(whisk.Api)
	api.Namespace = namespace
	api.Swagger = string(swagger)

	return api, nil
}

func (dm *YAMLParser) getGatewayMethods() []string {
	methods := []string{}
	for k := range whisk.ApiVerbs {
		methods = append(methods, k)
	}
	return methods
}

func doesPathParamExist(params []string, param string) bool {
	for _, e := range params {
		if e == param {
			return true
		}
	}
	return false
}

func getPathParameterNames(path string) []string {
	var pathParameters []string
	pathElements := strings.Split(path, PATH_SEPARATOR)
	for _, e := range pathElements {
		paramName := getParamName(e)
		if len(paramName) > 0 {
			if !doesPathParamExist(pathParameters, paramName) {
				pathParameters = append(pathParameters, paramName)
			}
		}
	}
	return pathParameters
}

func generatePathParameters(relativePath string) []whisk.ApiParameter {
	pathParams := []whisk.ApiParameter{}

	pathParamNames := getPathParameterNames(relativePath)
	for _, paramName := range pathParamNames {
		param := whisk.ApiParameter{Name: paramName, In: "path", Required: true, Type: "string",
			Description: wski18n.T("Default description for '{{.name}}'", map[string]interface{}{"name": paramName})}
		pathParams = append(pathParams, param)
	}

	return pathParams
}

func isParam(param string) bool {
	if strings.HasPrefix(param, PARAM_OPENING_BRACKET) &&
		strings.HasSuffix(param, PARAM_CLOSING_BRACKET) {
		return true
	}
	return false
}

func getParamName(param string) string {
	paramName := ""
	if isParam(param) {
		paramName = param[1 : len(param)-1]
	}
	return paramName
}

func isGatewayBasePathValid(basePath string) bool {
	// return false if base path is empty string
	if len(basePath) == 0 {
		return false
	}
	// drop preceding "/" if exists
	if strings.HasPrefix(basePath, PATH_SEPARATOR) {
		basePath = basePath[1:]
	}
	// slice base path into substrings seperated by "/"
	// if there are more than one substrings, basePath has path parameters
	// basePath will have path parameters if substrings count is more than 1
	basePathElements := strings.Split(basePath, PATH_SEPARATOR)
	if len(basePathElements) > 1 {
		return false
	}
	return true
}
