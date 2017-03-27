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
	"github.com/openwhisk/openwhisk-client-go/whisk"
)

// structs that denotes the sample manifest.yaml, wrapped yaml.v2
func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

type ParseYaml interface {
	Unmarshal(input []byte, deploy *ManifestYAML) error
	Marshal(manifest *ManifestYAML) (output []byte, err error)

	//Compose Package entity according to yaml content
	ComposePackages(manifestpath string) ([]*whisk.Package, error)

	// Compose Action entities according to yaml content
	ComposeActions(manifestpath string) ([]*whisk.Action, error)

	// Compose Trigger entities according to deployment and manifest yaml content
	ComposeTriggers(manifestpath string, deploymentpath string) ([]*whisk.Trigger, error)

	// Compose Rule entities according to yaml content
	ComposeRules(manifestpath string) ([]*whisk.Rule, error)
}

type YAMLParser struct {
	manifests []*ManifestYAML
	lastID    uint32
}

type Action struct {
	//mapping to wsk.Action.Version
	Version  string `yaml:"version"`  //used in manifest.yaml
	Location string `yaml:"location"` //used in manifest.yaml
	Runtime  string `yaml:"runtime"`  //used in manifest.yaml
	//mapping to wsk.Action.Namespace
	Namespace  string                 `yaml:"namespace"`  //used in deployment.yaml
	Credential string                 `yaml:"credential"` //used in deployment.yaml
	Inputs     map[string]interface{} `yaml:"inputs"`     //used in both manifest.yaml and deployment.yaml
	Outputs    map[string]interface{} `yaml:"outputs"`    //used in manifest.yaml
	//mapping to wsk.Action.Name
	Name        string
	Annotations map[string]interface{} `yaml:"annotations,omitempty"`
	//Parameters  map[string]interface{} `yaml:parameters` // used in manifest.yaml
}

type Sequence struct {
	Actions     string                 `yaml:"actions"` //used in manifest.yaml
	Annotations map[string]interface{} `yaml:"annotations,omitempty"`
}

type Dependency struct {
	Name    string
	Url     string
	Version string
}

type Trigger struct {
	//mapping to ????
	Feed string `yaml:"feed"` //used in manifest.yaml
	//mapping to wsk.Trigger.Namespace
	Namespace  string                 `yaml:"namespace"`  //used in deployment.yaml
	Credential string                 `yaml:"credential"` //used in deployment.yaml
	Inputs     map[string]interface{} `yaml:"inputs"`     //used in deployment.yaml
	//mapping to wsk.Trigger.Name
	Name        string
	Annotations map[string]interface{} `yaml:"annotations,omitempty"`
	Source      string                 `yaml:source` // used in manifest.yaml
	//Parameters  map[string]interface{} `yaml:parameters` // used in manifest.yaml
}

type Feed struct {
	Namespace  string            `yaml:"namespace"`  //used in deployment.yaml
	Credential string            `yaml:"credential"` //used in both manifest.yaml and deployment.yaml
	Inputs     map[string]string `yaml:"inputs"`     //used in deployment.yaml
	Location   string            `yaml:"location"`   //used in manifest.yaml
	Action     string            `yaml:"action"`     //used in manifest.yaml
	//TODO: need to define operation structure
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
	Name string
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
	Version           string                `yaml:"version"` //used in manifest.yaml
	License           string                `yaml:"license"` //used in manifest.yaml
	Repositories      []Repository          `yaml:"repositories,omitempty"`
	Dependencies      map[string]Dependency `yaml: dependencies`        // used in manifest.yaml
	Function          string                `yaml:"function"`           //used in deployment.yaml
	PackageCredential string                `yaml:"package_credential"` //used in deployment.yaml
	//mapping to wsk.SentPackageNoPublish.Namespace
	Namespace   string                 `yaml:"namespace"`  //used in deployment.yaml
	Credential  string                 `yaml:"credential"` //used in deployment.yaml
	Actions     map[string]Action      `yaml:"actions"`    //used in both manifest.yaml and deployment.yaml
	Triggers    map[string]Trigger     `yaml:"triggers"`   //used in both manifest.yaml and deployment.yaml
	Feeds       map[string]Feed        `yaml:"feeds"`      //used in both manifest.yaml and deployment.yaml
	Rules       map[string]Rule        `yaml:"rules"`      //used in both manifest.yaml and deployment.yaml
	Inputs      map[string]interface{} `yaml:"inputs"`     //used in deployment.yaml
	Sequences   map[string]Sequence    `yaml:"sequences"`
	Annotations map[string]interface{} `yaml:"annotations,omitempty"`
	//Parameters  map[string]interface{} `yaml: parameters` // used in manifest.yaml
}

type Application struct {
	Name       string             `yaml:"name"`      //used in deployment.yaml
	Namespace  string             `yaml:"namespace"` //used in deployment.yaml
	Credential string             `yaml:"credential"`
	BaseUrl    string             `yaml: baseUrl`
	Version    string             `yaml:"version"`
	Packages   map[string]Package `yaml:"packages"` //used in deployment.yaml
	Package    Package            `yaml:"package"`
}

type DeploymentYAML struct {
	Application Application `yaml:"application"` //used in deployment.yaml
	Filepath    string      //file path of the yaml file
}

type ManifestYAML struct {
	Package  Package `yaml:"package"` //used in both manifest.yaml and deployment.yaml
	Filepath string  //file path of the yaml file
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
	wskrule.Name = rule.Name
	//wskrule.Namespace = rule.Namespace
	pub := false
	wskrule.Publish = &pub
	wskrule.Trigger = rule.Trigger
	wskrule.Action = rule.Action
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
