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
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskenv"
)

// YAML schema key names
// DO NOT translate
const (
	YAML_KEY_ACTION     = "action"
	YAML_KEY_ANNOTATION = "annotoation"
	YAML_KEY_API        = "api"
	YAML_KEY_FEED       = "feed"
	YAML_KEY_MANIFEST   = "manifest"
	YAML_KEY_NAMESPACE  = "namespace"
	YAML_KEY_PACKAGES   = "packages"
	YAML_KEY_PROJECT    = "project"
	YAML_KEY_RULE       = "rule"
	YAML_KEY_SEQUENCE   = "sequence"
	YAML_KEY_TRIGGER    = "trigger"
	YAML_KEY_PACKAGE    = "package"
	YAML_KEY_SOURCE     = "source" // deprecated
)

// YAML schema key values
const (
	YAML_VALUE_BRANCH_MASTER = "master"
)

// default values
const (
	DEFAULT_PACKAGE_LICENSE = "unlicensed"
	DEFAULT_PACKAGE_VERSION = "0.0.1"
)

// Known Limit values
const (
	// supported
	LIMIT_VALUE_TIMEOUT     = "timeout"
	LIMIT_VALUE_MEMORY_SIZE = "memorySize"
	LIMIT_VALUE_LOG_SIZE    = "logSize"
	// unsupported
	LIMIT_VALUE_CONCURRENT_ACTIVATIONS = "concurrentActivations"
	LIMIT_VALUE_USER_INVOCATION_RATE   = "userInvocationRate"
	LIMIT_VALUE_CODE_SIZE              = "codeSize"
	LIMIT_VALUE_PARAMETER_SIZE         = "parameterSize"
)

var LIMITS_SUPPORTED = [](string){
	LIMIT_VALUE_TIMEOUT,
	LIMIT_VALUE_MEMORY_SIZE,
	LIMIT_VALUE_LOG_SIZE,
}

var LIMITS_UNSUPPORTED = [](string){
	LIMIT_VALUE_CONCURRENT_ACTIVATIONS,
	LIMIT_VALUE_USER_INVOCATION_RATE,
	LIMIT_VALUE_CODE_SIZE,
	LIMIT_VALUE_PARAMETER_SIZE,
}

// structs that denote the sample manifest.yaml, wrapped yaml.v2
func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

type YAMLParser struct {
	manifests []*YAML
	lastID    uint32
}

type Action struct {
	//mapping to wsk.Action.Version
	Version  string `yaml:"version"`           //used in manifest.yaml
	Location string `yaml:"location"`          //deprecated, used in manifest.yaml
	Function string `yaml:"function"`          //used in manifest.yaml
	Runtime  string `yaml:"runtime,omitempty"` //used in manifest.yaml
	//mapping to wsk.Action.Namespace
	Namespace  string               `yaml:"namespace"`  //used in deployment.yaml
	Credential string               `yaml:"credential"` //used in deployment.yaml
	Inputs     map[string]Parameter `yaml:"inputs"`     //used in both manifest.yaml and deployment.yaml
	Outputs    map[string]Parameter `yaml:"outputs"`    //used in manifest.yaml
	//mapping to wsk.Action.Name
	Name        string
	Annotations map[string]interface{} `yaml:"annotations,omitempty"`
	//Parameters  map[string]interface{} `yaml:parameters` // used in manifest.yaml
	ExposedUrl string  `yaml:"exposedUrl"` // used in manifest.yaml
	Webexport  string  `yaml:"web-export"` // used in manifest.yaml
	Main       string  `yaml:"main"`       // used in manifest.yaml
	Limits     *Limits `yaml:"limits"`     // used in manifest.yaml
}

type Limits struct {
	Timeout               *int `yaml:"timeout,omitempty"`               //in ms, [100 ms,300000ms]
	Memory                *int `yaml:"memorySize,omitempty"`            //in MB, [128 MB,512 MB]
	Logsize               *int `yaml:"logSize,omitempty"`               //in MB, [0MB,10MB]
	ConcurrentActivations *int `yaml:"concurrentActivations,omitempty"` //not changeable via APIs
	UserInvocationRate    *int `yaml:"userInvocationRate,omitempty"`    //not changeable via APIs
	CodeSize              *int `yaml:"codeSize,omitempty"`              //not changeable via APIs
	ParameterSize         *int `yaml:"parameterSize,omitempty"`         //not changeable via APIs
}

type Sequence struct {
	Actions     string                 `yaml:"actions"` //used in manifest.yaml
	Annotations map[string]interface{} `yaml:"annotations,omitempty"`
}

type Dependency struct {
	Version     string                 `yaml: "version, omitempty"`
	Location    string                 `yaml: "location, omitempty"`
	Inputs      map[string]Parameter   `yaml:"inputs"`
	Annotations map[string]interface{} `yaml:"annotations"`
}

type Parameter struct {
	Type        string      `yaml:"type,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Value       interface{} `yaml:"value,omitempty"` // JSON Value
	Required    bool        `yaml:"required,omitempty"`
	Default     interface{} `yaml:"default,omitempty"`
	Status      string      `yaml:"status,omitempty"`
	Schema      interface{} `yaml:"schema,omitempty"`
	multiline   bool
}

type Trigger struct {
	//mapping to ????
	Feed string `yaml:"feed"` //used in manifest.yaml
	//mapping to wsk.Trigger.Namespace
	Namespace  string               `yaml:"namespace"`  //used in deployment.yaml
	Credential string               `yaml:"credential"` //used in deployment.yaml
	Inputs     map[string]Parameter `yaml:"inputs"`     //used in deployment.yaml
	//mapping to wsk.Trigger.Name
	Name        string
	Annotations map[string]interface{} `yaml:"annotations,omitempty"`
	Source      string                 `yaml:source` // deprecated, used in manifest.yaml
	//Parameters  map[string]interface{} `yaml:parameters` // used in manifest.yaml
}

type Feed struct {
	Namespace  string            `yaml:"namespace"`  //used in deployment.yaml
	Credential string            `yaml:"credential"` //used in both manifest.yaml and deployment.yaml
	Inputs     map[string]string `yaml:"inputs"`     //used in deployment.yaml
	Location   string            `yaml:"location"`   //used in manifest.yaml
	Action     string            `yaml:"action"`     //used in manifest.yaml
	// TODO(): need to define operation structure
	Operations map[string]interface{} `yaml:"operations"` //used in manifest.yaml
	Name       string
}

type Rule struct {
	//mapping to wsk.Rule.Trigger
	Trigger string `yaml:"trigger"` //used in manifest.yaml
	//mapping to wsk.Rule.Action
	Action string `yaml:"action"` //used in manifest.yaml
	Rule   string `yaml:"rule"`   //used in manifest.yaml
	//mapping to wsk.Rule.Name
	Name        string
	Annotations map[string]interface{} `yaml:"annotations,omitempty"`
}

type Repository struct {
	Url         string `yaml:"url"`
	Description string `yaml:"description,omitempty"`
	Credential  string `yaml:"credential,omitempty"`
}

type Package struct {
	//mapping to wsk.SentPackageNoPublish.Name
	Packagename string `yaml:"name"` //used in manifest.yaml
	//mapping to wsk.SentPackageNoPublish.Version
	Version      string                `yaml:"version"` //used in manifest.yaml, mandatory
	License      string                `yaml:"license"` //used in manifest.yaml, mandatory
	Repositories []Repository          `yaml:"repositories,omitempty"`
	Dependencies map[string]Dependency `yaml: dependencies` //used in manifest.yaml
	//mapping to wsk.SentPackageNoPublish.Namespace
	Namespace        string                 `yaml:"namespace"`  //used in both manifest.yaml and deployment.yaml
	Credential       string                 `yaml:"credential"` //used in both manifest.yaml and deployment.yaml
	ApiHost          string                 `yaml:"apiHost"`    //used in both manifest.yaml and deployment.yaml
	ApigwAccessToken string                 `yaml:"apigwAccessToken"`
	Actions          map[string]Action      `yaml:"actions"`  //used in both manifest.yaml and deployment.yaml
	Triggers         map[string]Trigger     `yaml:"triggers"` //used in both manifest.yaml and deployment.yaml
	Feeds            map[string]Feed        `yaml:"feeds"`    //used in both manifest.yaml and deployment.yaml
	Rules            map[string]Rule        `yaml:"rules"`    //used in both manifest.yaml and deployment.yaml
	Inputs           map[string]Parameter   `yaml:"inputs"`   //deprecated, used in deployment.yaml
	Sequences        map[string]Sequence    `yaml:"sequences"`
	Annotations      map[string]interface{} `yaml:"annotations,omitempty"`

	// TODO() this is a convenience we want for package-shared vars that would be
	// propagated to every action within the package.
	//Parameters  map[string]interface{} `yaml: parameters` // used in manifest.yaml
	Apis map[string]map[string]map[string]map[string]string `yaml:"apis"` //used in manifest.yaml
}

type Project struct {
	Name             string             `yaml:"name"`      //used in deployment.yaml
	Namespace        string             `yaml:"namespace"` //used in deployment.yaml
	Credential       string             `yaml:"credential"`
	ApiHost          string             `yaml:"apiHost"`
	ApigwAccessToken string             `yaml:"apigwAccessToken"`
	Version          string             `yaml:"version"`
	Packages         map[string]Package `yaml:"packages"` //used in deployment.yaml
}

type YAML struct {
	Project  Project            `yaml:"project"`  //used in deployment.yaml
	Packages map[string]Package `yaml:"packages"` //used in deployment.yaml
	//Package     Package            `yaml:"package"`   // DEPRECATED.  Should we add warning if found?
	Filepath string //file path of the yaml file
}

// function to return Project or Application depending on what is specified in
// manifest and deployment files
func (yaml *YAML) GetProject() Project {
	return yaml.Project
}

func convertPackageName(packageMap map[string]Package) map[string]Package {
	packages := make(map[string]Package)
	for packName, depPacks := range packageMap {
		name := packName
		packageName := wskenv.InterpolateStringWithEnvVar(packName)
		if str, ok := packageName.(string); ok {
			name = str
		}
		depPacks.Packagename = wskenv.ConvertSingleName(depPacks.Packagename)
		packages[name] = depPacks
	}
	return packages
}

func ReadEnvVariable(yaml *YAML) *YAML {
	yaml.Project.Packages = convertPackageName(yaml.Project.Packages)
	yaml.Packages = convertPackageName(yaml.Packages)
	return yaml
}

//********************Trigger functions*************************//
//add the key/value array as the annotations of the trigger.
func (trigger *Trigger) ComposeWskTrigger(kvarr []whisk.KeyValue) *whisk.Trigger {
	wsktrigger := new(whisk.Trigger)
	wsktrigger.Name = trigger.Name
	wsktrigger.Namespace = trigger.Namespace
	pub := false
	wsktrigger.Publish = &pub
	wsktrigger.Annotations = kvarr
	return wsktrigger
}

//********************Rule functions*************************//
func (rule *Rule) ComposeWskRule() *whisk.Rule {
	wskrule := new(whisk.Rule)
	wskrule.Name = wskenv.ConvertSingleName(rule.Name)
	//wskrule.Namespace = rule.Namespace // TODO() ?
	pub := false
	wskrule.Publish = &pub
	wskrule.Trigger = wskenv.ConvertSingleName(rule.Trigger)
	wskrule.Action = wskenv.ConvertSingleName(rule.Action)
	return wskrule
}

//********************Package functions*************************//
func (pkg *Package) ComposeWskPackage() *whisk.Package {
	wskpag := new(whisk.Package)
	wskpag.Name = pkg.Packagename
	wskpag.Namespace = pkg.Namespace
	pub := false
	wskpag.Publish = &pub
	wskpag.Version = pkg.Version
	return wskpag
}

func (pkg *Package) GetActionList() []Action {
	var s1 []Action = make([]Action, 0)
	for action_name, action := range pkg.Actions {
		action.Name = action_name
		s1 = append(s1, action)
	}
	return s1
}

func (pkg *Package) GetTriggerList() []Trigger {
	var s1 []Trigger = make([]Trigger, 0)
	for trigger_name, trigger := range pkg.Triggers {
		trigger.Name = trigger_name
		s1 = append(s1, trigger)
	}
	return s1
}

func (pkg *Package) GetRuleList() []Rule {
	var s1 []Rule = make([]Rule, 0)
	for rule_name, rule := range pkg.Rules {
		rule.Name = rule_name
		s1 = append(s1, rule)
	}
	return s1
}

//This is for parse the deployment yaml file.
func (pkg *Package) GetFeedList() []Feed {
	var s1 []Feed = make([]Feed, 0)
	for feed_name, feed := range pkg.Feeds {
		feed.Name = feed_name
		s1 = append(s1, feed)
	}
	return s1
}

// This is for parse the manifest yaml file.
func (pkg *Package) GetApis() []*whisk.Api {
	var apis = make([]*whisk.Api, 0)
	for k, v := range pkg.Apis {
		var apiName string = k
		for k, v := range v {
			var gatewayBasePath string = k
			for k, v := range v {
				var gatewayRelPath string = k
				for k, v := range v {
					api := &whisk.Api{}
					api.ApiName = apiName
					api.GatewayBasePath = gatewayBasePath
					api.GatewayRelPath = gatewayRelPath
					action := &whisk.ApiAction{}
					action.Name = k
					action.BackendMethod = v
					api.Action = action
					apis = append(apis, api)
				}
			}
		}
	}
	return apis
}

//********************YAML functions*************************//
func (yaml *YAML) ComposeParsersPackage(wskpag whisk.Package) *Package {
	pkg := new(Package)
	pkg.Packagename = wskpag.Name
	pkg.Namespace = wskpag.Namespace
	pkg.Version = wskpag.Version

	for _, keyval := range wskpag.Parameters {
		param := new(Parameter)
		param.Value = keyval.Value
		pkg.Inputs[keyval.Key] = *param
	}

	return pkg
}

func (yaml *YAML) ComposeParsersAction(wskact whisk.Action) *Action {
	action := new(Action)
	action.Name = wskact.Name
	action.Namespace = wskact.Namespace
	action.Version = wskact.Version
	action.Runtime = wskact.Exec.Kind

	action.Inputs = make(map[string]Parameter)
	for _, keyval := range wskact.Parameters {
		param := new(Parameter)
		param.Value = keyval.Value
		action.Inputs[keyval.Key] = *param
	}

	return action
}

func (yaml *YAML) ComposeParsersTrigger(wsktrg whisk.Trigger) *Trigger {
	trigger := new(Trigger)
	trigger.Name = wsktrg.Name
	trigger.Namespace = wsktrg.Namespace

	for _, keyval := range wsktrg.Parameters {
		param := new(Parameter)
		param.Value = keyval.Value
		trigger.Inputs[keyval.Key] = *param
	}

	return trigger
}

func (yaml *YAML) ComposeParsersRule(wskrule whisk.Rule) *Rule {
	rule := new(Rule)
	rule.Name = wskrule.Name

	rule.Action = wskrule.Action.(map[string]interface{})["name"].(string)
	rule.Trigger = wskrule.Trigger.(map[string]interface{})["name"].(string)

	return rule
}
