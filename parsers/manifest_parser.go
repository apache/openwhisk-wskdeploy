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
	"errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"encoding/base64"
	"encoding/json"

	"fmt"
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"gopkg.in/yaml.v2"
)

// Read existing manifest file or create new if none exists
func ReadOrCreateManifest() (*ManifestYAML, error) {
	maniyaml := ManifestYAML{}

	if _, err := os.Stat(utils.ManifestFileNameYaml); err == nil {
		dat, _ := ioutil.ReadFile(utils.ManifestFileNameYaml)
		err := NewYAMLParser().Unmarshal(dat, &maniyaml)
		if err != nil {
			return &maniyaml, utils.NewInputYamlFileError(err.Error())
		}
	}
	return &maniyaml, nil
}

// Serialize manifest to local file
func Write(manifest *ManifestYAML, filename string) error {
	output, err := NewYAMLParser().Marshal(manifest)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	f, err := os.Create(filename)
	if err != nil {
		return utils.NewInputYamlFileError(err.Error())
	}
	defer f.Close()

	f.Write(output)
	return nil
}

func (dm *YAMLParser) Unmarshal(input []byte, manifest *ManifestYAML) error {
	err := yaml.UnmarshalStrict(input, manifest)
	if err != nil {
		return err
	}
	return nil
}

func (dm *YAMLParser) Marshal(manifest *ManifestYAML) (output []byte, err error) {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		fmt.Printf("err happened during marshal :%v", err)
		return nil, err
	}
	return data, nil
}

func (dm *YAMLParser) ParseManifest(manifestPath string) (*ManifestYAML, error) {
	mm := NewYAMLParser()
	maniyaml := ManifestYAML{}

	content, err := utils.Read(manifestPath)
	if err != nil {
		return &maniyaml, utils.NewInputYamlFileError(err.Error())
	}

	err = mm.Unmarshal(content, &maniyaml)
	if err != nil {
		lines, msgs := dm.convertErrorToLinesMsgs(err.Error())
		return &maniyaml, utils.NewParserErr(manifestPath, lines, msgs)
	}
	maniyaml.Filepath = manifestPath
	return &maniyaml, nil
}

func (dm *YAMLParser) ComposeDependencies(mani *ManifestYAML, projectPath string, filePath string) (map[string]utils.DependencyRecord, error) {

	var errorParser error
	depMap := make(map[string]utils.DependencyRecord)
	for key, dependency := range mani.Package.Dependencies {
		version := dependency.Version
		if version == "" {
			version = "master"
		}

		location := dependency.Location

		isBinding := false
		if utils.LocationIsBinding(location) {

			if !strings.HasPrefix(location, "/") {
				location = "/" + dependency.Location
			}

			isBinding = true
		} else if utils.LocationIsGithub(location) {

			if !strings.HasPrefix(location, "https://") && !strings.HasPrefix(location, "http://") {
				location = "https://" + dependency.Location
			}

			isBinding = false
		} else {
			return nil, errors.New("Dependency type is unknown.  wskdeploy only supports /whisk.system bindings or github.com packages.")
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
			keyVal.Value = utils.GetEnvVar(value)

			keyValArrAnot = append(keyValArrAnot, keyVal)
		}

		packDir := path.Join(projectPath, "Packages")
		depMap[key] = utils.NewDependencyRecord(packDir, mani.Package.Packagename, location, version, keyValArrParams, keyValArrAnot, isBinding)
	}

	return depMap, nil
}

func (dm *YAMLParser) ComposeAllPackages(manifest *ManifestYAML, filePath string) (map[string]*whisk.Package, error) {
	packages := map[string]*whisk.Package{}
	if manifest.Package.Packagename != "" {
		s, err := dm.ComposePackage(manifest.Package, manifest.Package.Packagename, filePath)
		if err == nil {
			packages[manifest.Package.Packagename] = s
		} else {
			return nil, err
		}
	} else if manifest.Packages != nil {
		for n, p := range manifest.Packages {
			s, err := dm.ComposePackage(p, n, filePath)
			if err == nil {
				packages[n] = s
			} else {
				return nil, err
			}
		}
	}
	return packages, nil
}

func (dm *YAMLParser) ComposePackage(pkg Package, packageName string, filePath string) (*whisk.Package, error) {
	var errorParser error
	pag := &whisk.Package{}
	pag.Name = packageName
	//The namespace for this package is absent, so we use default guest here.
	pag.Namespace = pkg.Namespace
	pub := false
	pag.Publish = &pub

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
	return pag, nil
}

func (dm *YAMLParser) ComposeSequencesFromAllPackages(namespace string, mani *ManifestYAML) (ar []utils.ActionRecord, err error) {
	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)
	if mani.Package.Packagename != "" {
		return dm.ComposeSequences(namespace, mani.Package.Sequences, mani.Package.Packagename)
	} else if mani.Packages != nil {
		for n, p := range mani.Packages {
			s, err := dm.ComposeSequences(namespace, p.Sequences, n)
			if err == nil {
				s1 = append(s1, s...)
			} else {
				return nil, err
			}
		}
	}
	return s1, nil
}

func (dm *YAMLParser) ComposeSequences(namespace string, sequences map[string]Sequence, packageName string) ([]utils.ActionRecord, error) {
	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)
	for key, sequence := range sequences {
		wskaction := new(whisk.Action)
		wskaction.Exec = new(whisk.Exec)
		wskaction.Exec.Kind = "sequence"
		actionList := strings.Split(sequence.Actions, ",")

		var components []string
		for _, a := range actionList {
			act := strings.TrimSpace(a)
			if !strings.ContainsRune(act, '/') && !strings.HasPrefix(act, packageName +"/") {
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
			keyVal.Value = utils.GetEnvVar(value)

			keyValArr = append(keyValArr, keyVal)
		}

		if len(keyValArr) > 0 {
			wskaction.Annotations = keyValArr
		}

		record := utils.ActionRecord{Action: wskaction, Packagename: packageName, Filepath: key}
		s1 = append(s1, record)
	}
	return s1, nil
}

func (dm *YAMLParser) ComposeActionsFromAllPackages(manifest *ManifestYAML, filePath string) ([]utils.ActionRecord, error) {
	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)
	if manifest.Package.Packagename != "" {
		return dm.ComposeActions(filePath, manifest.Package.Actions, manifest.Package.Packagename)
	} else if manifest.Packages != nil {
		for n, p := range manifest.Packages {
			a, err := dm.ComposeActions(filePath, p.Actions, n)
			if err == nil {
				s1 = append(s1, a...)
			} else {
				return nil, err
			}
		}
	}
	return s1, nil
}

func (dm *YAMLParser) ComposeActions(filePath string, actions map[string]Action, packageName string) ([]utils.ActionRecord, error) {

	var errorParser error
	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)

	for key, action := range actions {
		splitFilePath := strings.Split(filePath, string(os.PathSeparator))
		//set action.Function to action.Location
		//because Location is deprecated in Action entity
		if action.Function == "" && action.Location != "" {
			action.Function = action.Location
		}

		wskaction := new(whisk.Action)
		//bind action, and exposed URL

		wskaction.Exec = new(whisk.Exec)
		if action.Function != "" {
			filePath := strings.TrimRight(filePath, splitFilePath[len(splitFilePath)-1]) + action.Function

			if utils.IsDirectory(filePath) {
				zipName := filePath + ".zip"
				err := utils.NewZipWritter(filePath, zipName).Zip()
				if err != nil {
					return nil, err
				}
				defer os.Remove(zipName)
				// To do: support docker and main entry as did by go cli?
				wskaction.Exec, err = utils.GetExec(zipName, action.Runtime, false, "")
			} else {
				ext := path.Ext(filePath)
				var kind string

				switch ext {
				case ".swift":
					kind = "swift:3"
				case ".js":
					kind = "nodejs:6"
				case ".py":
					kind = "python"
				case ".java":
					kind = "java"
				case ".php":
					kind = "php:7.1"
				case ".jar":
					kind = "java"
				default:
					kind = "nodejs:6"
					errStr := wski18n.T("Unsupported runtime type, set to nodejs")
					whisk.Debug(whisk.DbgWarn, errStr)
					//add the user input kind here
				}

				wskaction.Exec.Kind = kind

				action.Function = filePath
				dat, err := utils.Read(filePath)
				if err != nil {
					return s1, err
				}
				code := string(dat)
				if ext == ".zip" || ext == ".jar" {
					code = base64.StdEncoding.EncodeToString([]byte(dat))
				}
				if ext == ".zip" && action.Runtime == "" {
                    utils.PrintOpenWhiskOutputln("need explicit action Runtime value")
				}
				wskaction.Exec.Code = &code
			}

		}

		if action.Runtime != "" {
			if utils.CheckExistRuntime(action.Runtime, utils.Rts) {
				wskaction.Exec.Kind = action.Runtime
			} else {
				errStr := wski18n.T("the runtime is not supported by Openwhisk platform.\n")
				whisk.Debug(whisk.DbgWarn, errStr)
			}
		} else {
			errStr := wski18n.T("wskdeploy has chosen a particular runtime for the action.\n")
			whisk.Debug(whisk.DbgWarn, errStr)
		}

		// we can specify the name of the action entry point using main
		if action.Main != "" {
			wskaction.Exec.Main = action.Main
		}

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
		if len(keyValArr) > 0 {
			wskaction.Parameters = keyValArr
		}

		keyValArr = make(whisk.KeyValueArr, 0)
		for name, value := range action.Annotations {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value = utils.GetEnvVar(value)

			keyValArr = append(keyValArr, keyVal)
		}

		// only set the webaction when the annotations are not empty.
		if action.Webexport == "true" {
			//wskaction.Annotations = keyValArr
			wskaction.Annotations, errorParser = utils.WebAction("yes", keyValArr, action.Name, false)
			if errorParser != nil {
				return s1, errorParser
			}
		}

		wskaction.Name = key
		pub := false
		wskaction.Publish = &pub

		record := utils.ActionRecord{Action: wskaction, Packagename: packageName, Filepath: action.Function}
		s1 = append(s1, record)
	}

	return s1, nil

}

func (dm *YAMLParser) ComposeTriggersFromAllPackages(manifest *ManifestYAML, filePath string) ([]*whisk.Trigger, error) {
	var triggers []*whisk.Trigger = make([]*whisk.Trigger, 0)
	if manifest.Package.Packagename != "" {
		return dm.ComposeTriggers(filePath, manifest.Package)
	} else if manifest.Packages != nil {
		for _, p := range manifest.Packages {
			t, err := dm.ComposeTriggers(filePath, p)
			if err == nil {
				triggers = append(triggers, t...)
			} else {
				return nil, err
			}
		}
	}
	return triggers, nil
}

func (dm *YAMLParser) ComposeTriggers(filePath string, pkg Package) ([]*whisk.Trigger, error) {
	var errorParser error
	var t1 []*whisk.Trigger = make([]*whisk.Trigger, 0)
	for _, trigger := range pkg.GetTriggerList() {
		wsktrigger := new(whisk.Trigger)
		wsktrigger.Name = trigger.Name
		wsktrigger.Namespace = trigger.Namespace
		pub := false
		wsktrigger.Publish = &pub

		//print warning information when .Source is not empty
		if trigger.Source != "" {
			warningString := wski18n.T("WARNING: The 'source' YAML key in trigger entity is deprecated. Please use 'feed' instead as described in specifications.\n")
			whisk.Debug(whisk.DbgWarn, warningString)
		}
		if trigger.Feed == "" {
			trigger.Feed = trigger.Source
		}

		keyValArr := make(whisk.KeyValueArr, 0)
		if trigger.Feed != "" {
			var keyVal whisk.KeyValue

			keyVal.Key = "feed"
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

		t1 = append(t1, wsktrigger)
	}
	return t1, nil
}

func (dm *YAMLParser) ComposeRulesFromAllPackages(manifest *ManifestYAML) ([]*whisk.Rule, error) {
	var rules []*whisk.Rule = make([]*whisk.Rule, 0)
	if manifest.Package.Packagename != "" {
		return dm.ComposeRules(manifest.Package, manifest.Package.Packagename)
	} else if manifest.Packages != nil {
		for n, p := range manifest.Packages {
			r, err := dm.ComposeRules(p, n)
			if err == nil {
				rules = append(rules, r...)
			} else {
				return nil, err
			}
		}
	}
	return rules, nil
}


func (dm *YAMLParser) ComposeRules(pkg Package, packageName string) ([]*whisk.Rule, error) {
	var r1 []*whisk.Rule = make([]*whisk.Rule, 0)
	for _, rule := range pkg.GetRuleList() {
		wskrule := rule.ComposeWskRule()
		act := strings.TrimSpace(wskrule.Action.(string))
		if !strings.ContainsRune(act, '/') && !strings.HasPrefix(act, packageName+"/") {
			act = path.Join(packageName, act)
		}
		wskrule.Action = act
		r1 = append(r1, wskrule)
	}
	return r1, nil
}

func (dm *YAMLParser) ComposeApiRecordsFromAllPackages(manifest *ManifestYAML) ([]*whisk.ApiCreateRequest, error) {
	var requests []*whisk.ApiCreateRequest = make([]*whisk.ApiCreateRequest, 0)
	if manifest.Package.Packagename != "" {
		return dm.ComposeApiRecords(manifest.Package)
	} else if manifest.Packages != nil {
		for _, p := range manifest.Packages {
			r, err := dm.ComposeApiRecords(p)
			if err == nil {
				requests = append(requests, r...)
			} else {
				return nil, err
			}
		}
	}
	return requests, nil
}

func (dm *YAMLParser) ComposeApiRecords(pkg Package) ([]*whisk.ApiCreateRequest, error) {
	var acq []*whisk.ApiCreateRequest = make([]*whisk.ApiCreateRequest, 0)
	apis := pkg.GetApis()
	for _, api := range apis {
		acr := new(whisk.ApiCreateRequest)
		acr.ApiDoc = api
		acq = append(acq, acr)
	}
	return acq, nil
}

// TODO(): Support other valid Package Manifest types
// TODO(): i.e., json (valid), timestamp, version, string256, string64, string16
// TODO(): Support JSON schema validation for type: json
// TODO(): Support OpenAPI schema validation

var validParameterNameMap = map[string]string{
	"string":  "string",
	"int":     "integer",
	"float":   "float",
	"bool":    "boolean",
	"int8":    "integer",
	"int16":   "integer",
	"int32":   "integer",
	"int64":   "integer",
	"float32": "float",
	"float64": "float",
}

var typeDefaultValueMap = map[string]interface{}{
	"string":  "",
	"integer": 0,
	"float":   0.0,
	"boolean": false,
	// TODO() Support these types + their validation
	// timestamp
	// null
	// version
	// string256
	// string64
	// string16
	// json
	// scalar-unit
	// schema
	// object
}

func isValidParameterType(typeName string) bool {
	_, isValid := typeDefaultValueMap[typeName]
	return isValid
}

// TODO(): throw errors
func getTypeDefaultValue(typeName string) interface{} {

	if val, ok := typeDefaultValueMap[typeName]; ok {
		return val
	} else {
		// TODO() throw an error "type not found"
	}
	return nil
}

func ResolveParamTypeFromValue(value interface{}, filePath string) (string, error) {
	// Note: string is the default type if not specified.
	var paramType string = "string"
	var err error = nil

	if value != nil {
		actualType := reflect.TypeOf(value).Kind().String()

		// See if the actual type of the value is valid
		if normalizedTypeName, found := validParameterNameMap[actualType]; found {
			// use the full spec. name
			paramType = normalizedTypeName

		} else {
			// raise an error if param is not a known type
			// TODO(): We have information to display to user on an error or warning here
			// TODO(): specifically, we have the parameter name, its value to show on error/warning
			// TODO(): perhaps this is a different Class of error?  e.g., ErrorParameterMismatchError
			lines := []string{"Line Unknown"}
			msgs := []string{"Parameter value is not a known type. [" + actualType + "]"}
			err = utils.NewParserErr(filePath, lines, msgs)
		}
	} else {

		// TODO: The value may be supplied later, we need to support non-fatal warnings
		// raise an error if param is nil
		//err = utils.NewParserErr("",-1,"Paramter value is nil.")
	}
	return paramType, err
}

// Resolve input parameter (i.e., type, value, default)
// Note: parameter values may set later (overriddNen) by an (optional) Deployment file
func ResolveParameter(paramName string, param *Parameter, filePath string) (interface{}, error) {

	var errorParser error
	var tempType string
	// default parameter value to empty string
	var value interface{} = ""

	// Trace Parameter struct before any resolution
	//dumpParameter(paramName, param, "BEFORE")

	// Parameters can be single OR multi-line declarations which must be processed/validated differently
	if !param.multiline {
		// we have a single-line parameter declaration
		// We need to identify parameter Type here for later validation
		param.Type, errorParser = ResolveParamTypeFromValue(param.Value, filePath)

		// In single-line format, the param's <value> can be a "Type name" and NOT an actual value.
		// if this is the case, we must detect it and set the value to the default for that type name.
		if param.Value!=nil && param.Type == "string" {
			// The value is a <string>; now we must test if is the name of a known Type
			var tempValue = param.Value.(string)
			if isValidParameterType(tempValue) {
				// If the value is indeed the name of a Type, we must change BOTH its
				// Type to be that type and its value to that Type's default value
				// (which happens later by setting it to nil here
				param.Type = param.Value.(string)
				param.Value = nil
			}
		}

	} else {
		// we have a multi-line parameter declaration

		// if we do not have a value, but have a default, use it for the value
		if param.Value == nil && param.Default != nil {
			param.Value = param.Default
		}

		// if we also have a type at this point, verify value (and/or default) matches type, if not error
		// Note: if either the value or default is in conflict with the type then this is an error
		tempType, errorParser = ResolveParamTypeFromValue(param.Value, filePath)

		// if we do not have a value or default, but have a type, find its default and use it for the value
		if param.Type != "" && !isValidParameterType(param.Type) {
			lines := []string{"Line Unknown"}
			msgs := []string{"Invalid Type for parameter. [" + param.Type + "]"}
			return value, utils.NewParserErr(filePath, lines, msgs)
		} else if param.Type == "" {
			param.Type = tempType
		}
	}

	// Make sure the parameter's value is a valid, non-empty string and startsWith '$" (dollar) sign
	value = utils.GetEnvVar(param.Value)

	typ := param.Type

	// TODO(Priti): need to validate type is one of the supported primitive types with unit testing
	// TODO(): with the new logic, when would the following Unmarhsall() call be used?
	// if value is of type 'string' and its not empty <OR> if type is not 'string'
	if str, ok := value.(string); ok && (len(typ) == 0 || typ != "string") {
		var parsed interface{}
		err := json.Unmarshal([]byte(str), &parsed)
		if err == nil {
			return parsed, err
		}
	}

	// Default to an empty string, do NOT error/terminate as Value may be provided later bu a Deployment file.
	if value == nil {
		value = getTypeDefaultValue(param.Type)
		// @TODO(): Need warning message here to warn of default usage, support for warnings (non-fatal)
	}

	// Trace Parameter struct after resolution
	//dumpParameter(paramName, param, "AFTER")
	//fmt.Printf("EXIT: value=[%v]\n", value)

	return value, errorParser
}

// Provide custom Parameter marshalling and unmarshalling

type ParsedParameter Parameter

func (n *Parameter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux ParsedParameter

	// Attempt to unmarshall the multi-line schema
	if err := unmarshal(&aux); err == nil {
		n.multiline = true
		n.Type = aux.Type
		n.Description = aux.Description
		n.Value = aux.Value
		n.Required = aux.Required
		n.Default = aux.Default
		n.Status = aux.Status
		n.Schema = aux.Schema
		return nil
	}

	// If we did not find the multi-line schema, assume in-line (or single-line) schema
	var inline interface{}
	if err := unmarshal(&inline); err != nil {
		return err
	}

	n.Value = inline
	n.multiline = false
	return nil
}

func (n *Parameter) MarshalYAML() (interface{}, error) {
	if _, ok := n.Value.(string); len(n.Type) == 0 && len(n.Description) == 0 && ok {
		if !n.Required && len(n.Status) == 0 && n.Schema == nil {
			return n.Value.(string), nil
		}
	}

	return n, nil
}

// Provides debug/trace support for Parameter type
func dumpParameter(paramName string, param *Parameter, separator string) {

	fmt.Printf("%s:\n", separator)
	fmt.Printf("\t%s: (%T)\n", paramName, param)
	if param != nil {
		fmt.Printf("\t\tParameter.Descrption: [%s]\n", param.Description)
		fmt.Printf("\t\tParameter.Type: [%s]\n", param.Type)
		fmt.Printf("\t\tParameter.Value: [%v]\n", param.Value)
		fmt.Printf("\t\tParameter.Default: [%v]\n", param.Default)
	}
}
