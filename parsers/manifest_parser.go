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
	"strings"

	"encoding/base64"

	"fmt"
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"gopkg.in/yaml.v2"
)

// Read existing manifest file or create new if none exists
func ReadOrCreateManifest() (*YAML, error) {
	maniyaml := YAML{}

	if _, err := os.Stat(utils.ManifestFileNameYaml); err == nil {
		dat, _ := ioutil.ReadFile(utils.ManifestFileNameYaml)
		err := NewYAMLParser().Unmarshal(dat, &maniyaml)
		if err != nil {
			return &maniyaml, utils.NewYAMLFileReadError(err.Error())
		}
	}
	return &maniyaml, nil
}

// Serialize manifest to local file
func Write(manifest *YAML, filename string) error {
	output, err := NewYAMLParser().Marshal(manifest)
	if err != nil {
		return utils.NewYAMLFormatError(err.Error())
	}

	f, err := os.Create(filename)
	if err != nil {
		return utils.NewYAMLFileReadError(err.Error())
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

func (dm *YAMLParser) Marshal(manifest *YAML) (output []byte, err error) {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		fmt.Printf("err happened during marshal :%v", err)
		return nil, err
	}
	return data, nil
}

func (dm *YAMLParser) ParseManifest(manifestPath string) (*YAML, error) {
	mm := NewYAMLParser()
	maniyaml := YAML{}

	content, err := utils.Read(manifestPath)
	if err != nil {
		return &maniyaml, utils.NewYAMLFileReadError(err.Error())
	}

	err = mm.Unmarshal(content, &maniyaml)
	if err != nil {
		lines, msgs := dm.convertErrorToLinesMsgs(err.Error())
		return &maniyaml, utils.NewYAMLParserErr(manifestPath, lines, msgs)
	}
	maniyaml.Filepath = manifestPath
	manifest := ReadEnvVariable(&maniyaml)

	return manifest, nil
}

func (dm *YAMLParser) ComposeDependenciesFromAllPackages(manifest *YAML, projectPath string, filePath string) (map[string]utils.DependencyRecord, error) {
	dependencies := make(map[string]utils.DependencyRecord)
	packages := make(map[string]Package)
	if manifest.Package.Packagename != "" {
		return dm.ComposeDependencies(manifest.Package, projectPath, filePath, manifest.Package.Packagename)
	} else {
		if len(manifest.Packages) != 0 {
			packages = manifest.Packages
		} else {
			packages = manifest.GetProject().Packages
		}
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
		depName := packageName + ":" + key
		depMap[depName] = utils.NewDependencyRecord(packDir, packageName, location, version, keyValArrParams, keyValArrAnot, isBinding)
	}

	return depMap, nil
}

func (dm *YAMLParser) ComposeAllPackages(manifest *YAML, filePath string, ma whisk.KeyValue) (map[string]*whisk.Package, error) {
	packages := map[string]*whisk.Package{}
	manifestPackages := make(map[string]Package)

	if manifest.Package.Packagename != "" {
		fmt.Println("WARNING: using package inside of manifest file will soon be deprecated, please use packages instead.")
		s, err := dm.ComposePackage(manifest.Package, manifest.Package.Packagename, filePath, ma)
		if err == nil {
			packages[manifest.Package.Packagename] = s
		} else {
			return nil, err
		}
	} else {
		if len(manifest.Packages) != 0 {
			manifestPackages = manifest.Packages
		} else {
			manifestPackages = manifest.GetProject().Packages
		}
	}

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
	if pkg.Version == "" {
		warningString := wski18n.T("WARNING: Mandatory field Package Version must be set.\n")
		whisk.Debug(whisk.DbgWarn, warningString)
		warningString = wski18n.T("WARNING: Package Version is not saved in the current wskdeploy version.\n")
		whisk.Debug(whisk.DbgWarn, warningString)
		pkg.Version = "0.0.1"
	}

	//License is a mandatory value
	//set license to unknown if it is an empty string
	//And print an warning message
	if pkg.License == "" {
		warningString := wski18n.T("WARNING: Mandatory field Package License must be set.\n")
		whisk.Debug(whisk.DbgWarn, warningString)
		warningString = wski18n.T("WARNING: Package License is not saved in the current wskdeploy version.\n")
		whisk.Debug(whisk.DbgWarn, warningString)
		pkg.License = "unlicensed"
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
		keyVal.Value = utils.GetEnvVar(value)
		listOfAnnotations = append(listOfAnnotations, keyVal)
	}
	if len(listOfAnnotations) > 0 {
		pag.Annotations = append(pag.Annotations, listOfAnnotations...)
	}

	// add Managed Annotations if this is Managed Deployment
	if utils.Flags.Managed {
		pag.Annotations = append(pag.Annotations, ma)
	}

	return pag, nil
}

func (dm *YAMLParser) ComposeSequencesFromAllPackages(namespace string, mani *YAML, ma whisk.KeyValue) ([]utils.ActionRecord, error) {
	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)

	manifestPackages := make(map[string]Package)

	if mani.Package.Packagename != "" {
		return dm.ComposeSequences(namespace, mani.Package.Sequences, mani.Package.Packagename, ma)
	} else {
		if len(mani.Packages) != 0 {
			manifestPackages = mani.Packages
		} else {
			manifestPackages = mani.GetProject().Packages
		}
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

	if manifest.Package.Packagename != "" {
		return dm.ComposeActions(filePath, manifest.Package.Actions, manifest.Package.Packagename, ma)
	} else {
		if len(manifest.Packages) != 0 {
			manifestPackages = manifest.Packages
		} else {
			manifestPackages = manifest.GetProject().Packages
		}
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
				// TODO() add the user input kind here if interactive
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

		/*
 		 *  Action.Runtime
 		 */
		if action.Runtime != "" {
			if utils.CheckExistRuntime(action.Runtime, utils.Rts) {
				wskaction.Exec.Kind = action.Runtime

			} else if utils.Flags.Strict {
				wskaction.Exec.Kind = action.Runtime
			} else {
				errStr := wski18n.T("wskdeploy has chosen a particular runtime for the action.\n")
				whisk.Debug(whisk.DbgWarn, errStr)
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
			//wskaction.Annotations  // TBD
		}

		/*
 		 *  Action.Annotations
 		 */
		listOfAnnotations := make(whisk.KeyValueArr, 0)
		for name, value := range action.Annotations {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value = utils.GetEnvVar(value)
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
		if action.Webexport == "true" {
			wskaction.Annotations, errorParser = utils.WebAction("yes", listOfAnnotations, false)
			if errorParser != nil {
				return s1, errorParser
			}
		}

		/*
 		 *  Action.Limits
 		 */
		if action.Limits!=nil {
			wsklimits :=  new(whisk.Limits)
			if utils.LimitsTimeoutValidation(action.Limits.Timeout) {
				wsklimits.Timeout = action.Limits.Timeout
			} else {
				warningString := wski18n.T("WARNING: Invalid limitation 'timeout' of action in manifest is ignored. Please check errors.\n")
				whisk.Debug(whisk.DbgWarn, warningString)
			}
			if utils.LimitsMemoryValidation(action.Limits.Memory) {
				wsklimits.Memory = action.Limits.Memory
			} else {
				warningString := wski18n.T("WARNING: Invalid limitation 'memorySize' of action in manifest is ignored. Please check errors.\n")
				whisk.Debug(whisk.DbgWarn, warningString)
			}
			if utils.LimitsLogsizeValidation(action.Limits.Logsize) {
				wsklimits.Logsize = action.Limits.Logsize
			} else {
				warningString := wski18n.T("WARNING: Invalid limitation 'logSize' of action in manifest is ignored. Please check errors.\n")
				whisk.Debug(whisk.DbgWarn, warningString)
			}
			if wsklimits.Timeout!=nil || wsklimits.Memory!=nil || wsklimits.Logsize!=nil {
				wskaction.Limits = wsklimits
			}

			//emit warning errors if these limits are not nil
			utils.NotSupportLimits(action.Limits.ConcurrentActivations,"concurrentActivations")
			utils.NotSupportLimits(action.Limits.UserInvocationRate,"userInvocationRate")
			utils.NotSupportLimits(action.Limits.CodeSize,"codeSize")
			utils.NotSupportLimits(action.Limits.ParameterSize,"parameterSize")
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

	if manifest.Package.Packagename != "" {
		return dm.ComposeTriggers(filePath, manifest.Package, ma)
	} else {
		if len(manifest.Packages) != 0 {
			manifestPackages = manifest.Packages
		} else {
			manifestPackages = manifest.GetProject().Packages
		}
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

		listOfAnnotations := make(whisk.KeyValueArr, 0)
		for name, value := range trigger.Annotations {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value = utils.GetEnvVar(value)
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

func (dm *YAMLParser) ComposeRulesFromAllPackages(manifest *YAML) ([]*whisk.Rule, error) {
	var rules []*whisk.Rule = make([]*whisk.Rule, 0)
	manifestPackages := make(map[string]Package)

	if manifest.Package.Packagename != "" {
		return dm.ComposeRules(manifest.Package, manifest.Package.Packagename)
	} else {
		if len(manifest.Packages) != 0 {
			manifestPackages = manifest.Packages
		} else {
			manifestPackages = manifest.GetProject().Packages
		}
	}

	for n, p := range manifestPackages {
		r, err := dm.ComposeRules(p, n)
		if err == nil {
			rules = append(rules, r...)
		} else {
			return nil, err
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

func (dm *YAMLParser) ComposeApiRecordsFromAllPackages(manifest *YAML) ([]*whisk.ApiCreateRequest, error) {
	var requests []*whisk.ApiCreateRequest = make([]*whisk.ApiCreateRequest, 0)
	manifestPackages := make(map[string]Package)

	if manifest.Package.Packagename != "" {
		return dm.ComposeApiRecords(manifest.Package)
	} else {
		if len(manifest.Packages) != 0 {
			manifestPackages = manifest.Packages
		} else {
			manifestPackages = manifest.GetProject().Packages
		}
	}

	for _, p := range manifestPackages {
		r, err := dm.ComposeApiRecords(p)
		if err == nil {
			requests = append(requests, r...)
		} else {
			return nil, err
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
