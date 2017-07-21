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
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"encoding/base64"
	"encoding/json"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"gopkg.in/yaml.v2"
	"fmt"
)

// Read existing manifest file or create new if none exists
func ReadOrCreateManifest() *ManifestYAML {
	maniyaml := ManifestYAML{}

	if _, err := os.Stat("manifest.yaml"); err == nil {
		dat, _ := ioutil.ReadFile("manifest.yaml")
		err := NewYAMLParser().Unmarshal(dat, &maniyaml)
		utils.Check(err)
	}
	return &maniyaml
}

// Serialize manifest to local file
func Write(manifest *ManifestYAML, filename string) {
	output, err := NewYAMLParser().Marshal(manifest)
	utils.Check(err)

	f, err := os.Create(filename)
	utils.Check(err)
	defer f.Close()

	f.Write(output)
}

func (dm *YAMLParser) Unmarshal(input []byte, manifest *ManifestYAML) error {
	err := yaml.Unmarshal(input, manifest)
	if err != nil {
		log.Printf("error happened during unmarshal :%v", err)
		return err
	}
	return nil
}

func (dm *YAMLParser) Marshal(manifest *ManifestYAML) (output []byte, err error) {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		log.Printf("err happened during marshal :%v", err)
		return nil, err
	}
	return data, nil
}

func (dm *YAMLParser) ParseManifest(mani string) *ManifestYAML {
	mm := NewYAMLParser()
	maniyaml := ManifestYAML{}

	content, err := utils.Read(mani)
	utils.Check(err)

	err = mm.Unmarshal(content, &maniyaml)
	utils.Check(err)
	maniyaml.Filepath = mani
	return &maniyaml
}

func (dm *YAMLParser) ComposeDependencies(mani *ManifestYAML, projectPath string) (map[string]utils.DependencyRecord, error) {

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

			keyVal.Value = ResolveParameter(keyVal.Value, &param)

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
		depMap[key] = utils.DependencyRecord{packDir, mani.Package.Packagename, location, version, keyValArrParams, keyValArrAnot, isBinding}
	}

	return depMap, nil
}

// Is we consider multi pacakge in one yaml?
func (dm *YAMLParser) ComposePackage(mani *ManifestYAML) (*whisk.Package, error) {
	var errorParser error

	//mani := dm.ParseManifest(manipath)
	pag := &whisk.Package{}
	pag.Name = mani.Package.Packagename
	//The namespace for this package is absent, so we use default guest here.
	pag.Namespace = mani.Package.Namespace
	pub := false
	pag.Publish = &pub

	keyValArr := make(whisk.KeyValueArr, 0)
	for name, param := range mani.Package.Inputs {
		var keyVal whisk.KeyValue
		keyVal.Key = name

		keyVal.Value = ResolveParameter(keyVal.Value, &param)

		if keyVal.Value != nil {
			keyValArr = append(keyValArr, keyVal)
		}
	}

	if len(keyValArr) > 0 {
		pag.Parameters = keyValArr
	}
	return pag, nil
}

func (dm *YAMLParser) ComposeSequences(namespace string, mani *ManifestYAML) ([]utils.ActionRecord, error) {
	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)
	for key, sequence := range mani.Package.Sequences {
		wskaction := new(whisk.Action)
		wskaction.Exec = new(whisk.Exec)
		wskaction.Exec.Kind = "sequence"
		actionList := strings.Split(sequence.Actions, ",")

		var components []string
		for _, a := range actionList {

			act := strings.TrimSpace(a)

			if !strings.ContainsRune(act, '/') && !strings.HasPrefix(act, mani.Package.Packagename+"/") {
				act = path.Join(mani.Package.Packagename, act)
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

		record := utils.ActionRecord{wskaction, mani.Package.Packagename, key}
		s1 = append(s1, record)
	}
	return s1, nil
}

func (dm *YAMLParser) ComposeActions(mani *ManifestYAML, manipath string) (ar []utils.ActionRecord, aub []*utils.ActionExposedURLBinding, err error) {

	var errorParser error
	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)
	var au []*utils.ActionExposedURLBinding = make([]*utils.ActionExposedURLBinding, 0)

	for key, action := range mani.Package.Actions {
		splitmanipath := strings.Split(manipath, string(os.PathSeparator))

		wskaction := new(whisk.Action)
		//bind action, and exposed URL
		aubinding := new(utils.ActionExposedURLBinding)
		aubinding.ActionName = key
		aubinding.ExposedUrl = action.ExposedUrl

		wskaction.Exec = new(whisk.Exec)
		if action.Location != "" {
			filePath := strings.TrimRight(manipath, splitmanipath[len(splitmanipath)-1]) + action.Location

			if utils.IsDirectory(filePath) {
				zipName := filePath + ".zip"
				err = utils.NewZipWritter(filePath, zipName).Zip()
				defer os.Remove(zipName)
				utils.Check(err)
				// To do: support docker and main entry as did by go cli?
				wskaction.Exec, err = utils.GetExec(zipName, action.Runtime, false, "")
			} else {
				ext := path.Ext(filePath)
				kind := "nodejs:default"

				switch ext {
				case ".swift":
					kind = "swift:default"
				case ".js":
					kind = "nodejs:default"
				case ".py":
					kind = "python"
				}

				wskaction.Exec.Kind = kind

				action.Location = filePath
				dat, err := utils.Read(filePath)
				utils.Check(err)
				code := string(dat)
				if ext == ".zip" || ext == ".jar" {
					code = base64.StdEncoding.EncodeToString([]byte(dat))
				}
				wskaction.Exec.Code = &code
			}

		}

		if action.Runtime != "" {
			wskaction.Exec.Kind = action.Runtime
		}

		// we can specify the name of the action entry point using main
		if action.Main != "" {
			wskaction.Exec.Main = action.Main
		}

		keyValArr := make(whisk.KeyValueArr, 0)
		for name, param := range action.Inputs {
			var keyVal whisk.KeyValue
			keyVal.Key = name

			keyVal.Value = ResolveParameter(keyVal.Value, &param)

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
		if len(keyValArr) > 0 && action.Webexport == "true" {
			//wskaction.Annotations = keyValArr
			wskaction.Annotations, err = utils.WebAction("yes", keyValArr, action.Name, false)
			utils.Check(err)
		}

		wskaction.Name = key
		pub := false
		wskaction.Publish = &pub

		record := utils.ActionRecord{wskaction, mani.Package.Packagename, action.Location}
		s1 = append(s1, record)

		//only append when the fields are exists
		if aubinding.ActionName != "" && aubinding.ExposedUrl != "" {
			au = append(au, aubinding)
		}

	}

	return s1, au, nil

}

func (dm *YAMLParser) ComposeTriggers(manifest *ManifestYAML) ([]*whisk.Trigger, error) {
	var errorParser error
	var t1 []*whisk.Trigger = make([]*whisk.Trigger, 0)
	pkg := manifest.Package
	for _, trigger := range pkg.GetTriggerList() {
		wsktrigger := new(whisk.Trigger)
		wsktrigger.Name = trigger.Name
		wsktrigger.Namespace = trigger.Namespace
		pub := false
		wsktrigger.Publish = &pub

		keyValArr := make(whisk.KeyValueArr, 0)
		if trigger.Source != "" {
			var keyVal whisk.KeyValue

			keyVal.Key = "feed"
			keyVal.Value = trigger.Source

			keyValArr = append(keyValArr, keyVal)

			wsktrigger.Annotations = keyValArr
		}

		keyValArr = make(whisk.KeyValueArr, 0)
		for name, param := range trigger.Inputs {
			var keyVal whisk.KeyValue
			keyVal.Key = name

			keyVal.Value = ResolveParameter(keyVal.Value, &param)

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

func (dm *YAMLParser) ComposeRules(manifest *ManifestYAML) ([]*whisk.Rule, error) {

	var r1 []*whisk.Rule = make([]*whisk.Rule, 0)
	pkg := manifest.Package
	for _, rule := range pkg.GetRuleList() {
		wskrule := rule.ComposeWskRule()

		act := strings.TrimSpace(wskrule.Action.(string))

		if !strings.ContainsRune(act, '/') && !strings.HasPrefix(act, pkg.Packagename+"/") {
			act = path.Join(pkg.Packagename, act)
		}

		wskrule.Action = act

		r1 = append(r1, wskrule)
	}

	return r1, nil
}

func (action *Action) ComposeWskAction(manipath string) (*whisk.Action, error) {
	wskaction, err := utils.CreateActionFromFile(manipath, action.Location)
	utils.Check(err)
	wskaction.Name = action.Name
	wskaction.Version = action.Version
	wskaction.Namespace = action.Namespace
	return wskaction, err
}


// TODO() Support other valid Package Manifest types
// TODO() i.e., json (valid), timestamp, version, string256, string64, string16
// TODO() Support JSON schema validation for type: json
// TODO(): Support OpenAPI schema validation

var validParameterTypeDefaultMap = map[string]string{
	"string": "",
	"integer" : "0",
	"float" : "0.0",
	"boolean" : "false",
}

func ResolveParamType(param *Parameter) (string, string) {

	var paramType string = "string"
	var defaultValue string = ""

	// TODO() raise an error if param is nil
	if param.Type != "" {

		paramType = param.Type
		if defaultValue, ok := validParameterTypeDefaultMap[paramType]; !ok {
			// TODO() emit an error, cease processing, inform user where error occurred in YAML
			println(ok, defaultValue)
		}
	}
	return paramType, defaultValue
}

// Resolve input parameter (i.e., type, value, default)
// Note: parameter values may set later (overridden) by an (optional) Deployment file
func ResolveParameter(simpleValue interface{}, param *Parameter) interface{} {
	// default parameter value to empty string
        var value interface{} = ""

	//// Resolve type of parameter's value
	//paramType, typeDefault := ResolveParamType(param)
	//println("type=["+paramType+"] defaultvalue=["+typeDefault+"]")
	//
	//// If param.Value is nil, that implies a "multi-line" (complex) parameter schema
	//if param.Value == nil {
	//
	//	// TODO()
	//	return value
	//}

	// Make sure the parameter's value is a valid, non-empty string and startsWith '$" (dollar) sign
	if strValue, ok := value.(string); ok && strValue != "" && len(strValue) > 0 {
		value = utils.GetEnvVar(param.Value)
	}

	typ := param.Type
	// if value is of type 'string' and its not empty <OR> if type is not 'string'
	// TODO(): need to validate type is one of the supported primitive types with unit testing
	if str, ok := value.(string); ok && (len(typ) == 0 || typ != "string") {
		var parsed interface{}
		err := json.Unmarshal([]byte(str), &parsed)
		if err == nil {
			return parsed, err
		}
	}

	// Default to an empty string, do NOT error/terminate as Value may be provided later bu a Deployment file.
	if (value == nil) {
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
	if(param!= nil) {
		fmt.Printf("\t\tParameter.Descrption: [%s]\n", param.Description)
		fmt.Printf("\t\tParameter.Type: [%s]\n", param.Type)
		fmt.Printf("\t\tParameter.Value: [%v]\n", param.Value)
		fmt.Printf("\t\tParameter.Default: [%v]\n", param.Default)
	}
}
