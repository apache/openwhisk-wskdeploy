package utils

import (
	"github.com/openwhisk/openwhisk-client-go/whisk"
	"gopkg.in/yaml.v2"
	"log"
)

// structs that denotes the sample manifest.yaml, wrapped yaml.v2
func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

type ParseDeployerYaml interface {
	Unmarshal(input []byte, deploy *ManifestYAML) error
	Marshal(manifest *ManifestYAML) (output []byte, err error)

	//Compose Package entity according to yaml content
	ComposePackage(manifestpath string) ([]*whisk.Package, error)

	// Compose Action entities according to yaml content
	ComposeActions(manifestpath string) ([]*whisk.Action, error)

	// Compose Trigger entities according to yaml content
	ComposeTriggers(manifestpath string) ([]*whisk.Trigger, error)

	// Compose Rule entities according to yaml content
	ComposeRule(manifestpath string) ([]*whisk.Rule, error)
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
	Name string
}

type Trigger struct {
	//mapping to ????
	Feed string `yaml:"feed"` //used in manifest.yaml
	//mapping to wsk.Trigger.Namespace
	Namespace  string                 `yaml:"namespace"`  //used in deployment.yaml
	Credential string                 `yaml:"credential"` //used in deployment.yaml
	Inputs     map[string]interface{} `yaml:"inputs"`     //used in deployment.yaml
	//mapping to wsk.Trigger.Name
	Name string
}

type Feed struct {
	Namespace  string            `yaml:"namespace"`  //used in deployment.yaml
	Credential string            `yaml:"credential"` //used in both manifest.yaml and deployment.yaml
	Inputs     map[string]string `yaml:"inputs"`     //used in deployment.yaml
	Location string `yaml:"location"`    //used in manifest.yaml
	Action string `yaml:"action"`    //used in manifest.yaml
	//TODO: need to define operation structure
	Operations map[string]interface{} `yaml:"operations"`    //used in manifest.yaml
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

type Package struct {
	//mapping to wsk.SentPackageNoPublish.Name
	Packagename string `yaml:"name"` //used in manifest.yaml
	//mapping to wsk.SentPackageNoPublish.Version
	Version           string `yaml:"version"`            //used in manifest.yaml
	License           string `yaml:"license"`            //used in manifest.yaml
	Function          string `yaml:"function"`           //used in deployment.yaml
	PackageCredential string `yaml:"package_credential"` //used in deployment.yaml
	//mapping to wsk.SentPackageNoPublish.Namespace
	Namespace  string             `yaml:"namespace"`  //used in deployment.yaml
	Credential string             `yaml:"credential"` //used in deployment.yaml
	Actions    map[string]Action  `yaml:"actions"`    //used in both manifest.yaml and deployment.yaml
	Triggers   map[string]Trigger `yaml:"triggers"`   //used in both manifest.yaml and deployment.yaml
	Feeds	map[string]Feed `yaml:"feeds"`   //used in both manifest.yaml and deployment.yaml
	Rules      map[string]Rule    `yaml:"rules"`      //used in both manifest.yaml and deployment.yaml
	Inputs     map[string]string  `yaml:"inputs"`     //used in deployment.yaml
}

type Application struct {
	Name      string             `yaml:"name"`      //used in deployment.yaml
	Namespace string             `yaml:"namespace"` //used in deployment.yaml
	Packages  map[string]Package `yaml:"packages"`  //used in deployment.yaml
}

type DeploymentYAML struct {
	Application Application `yaml:"application"` //used in deployment.yaml
	Filepath    string      //file path of the yaml file
}

type ManifestYAML struct {
	Package  Package `yaml:"package"` //used in both manifest.yaml and deployment.yaml
	Filepath string  //file path of the yaml file
}

//********************Action functions*************************//
func (action *Action) ComposeWskAction(manipath string) (*whisk.Action, error) {
	wskaction, err := CreateActionFromFile(manipath, action.Location)
	Check(err)
	wskaction.Name = action.Name
	wskaction.Version = action.Version
	wskaction.Namespace = action.Namespace
	return wskaction, err
}

//********************Trigger functions*************************//
func (trigger *Trigger) ComposeWskTrigger() *whisk.Trigger {
	wsktrigger := new(whisk.Trigger)
	wsktrigger.Name = trigger.Name
	wsktrigger.Namespace = trigger.Namespace
	wsktrigger.Publish = false
	return wsktrigger
}

//********************Rule functions*************************//
func (rule *Rule) ComposeWskRule() *whisk.Rule {
	wskrule := new(whisk.Rule)
	wskrule.Name = rule.Name
	//wskrule.Namespace = rule.Namespace
	wskrule.Publish = false
	wskrule.Trigger = rule.Trigger
	wskrule.Action = rule.Action
	return wskrule
}

//********************Package functions*************************//
func (pkg *Package) ComposeWskPackage() *whisk.SentPackageNoPublish {
	wskpag := new(whisk.SentPackageNoPublish)
	wskpag.Name = pkg.Packagename
	wskpag.Namespace = pkg.Namespace
	wskpag.Publish = false
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

func (pkg *Package) GetFeedList() []Feed {
	var s1 []Feed = make([]Feed, 0)
	for feed_name, feed := range pkg.Feeds {
		feed.Name = feed_name
		s1 = append(s1, feed)
	}
	return s1
}

//********************Application functions*************************//
func (app *Application) GetPackageList() []Package {
	var s1 []Package = make([]Package, 0)
	for pkg_name, pkg := range app.Packages {
		pkg.Packagename = pkg_name
		s1 = append(s1, pkg)
	}
	return s1
}

//********************YAMLParser functions*************************//
func (dm *YAMLParser) ParseManifest(mani string) *ManifestYAML {
	mm := NewYAMLParser()
	maniyaml := ManifestYAML{}
	content, err := new(ContentReader).LocalReader.ReadLocal(mani)
	Check(err)
	err = mm.Unmarshal(content, &maniyaml)
	maniyaml.Filepath = mani
	return &maniyaml
}

func (dm *YAMLParser) ParseDeployment(dply string) *DeploymentYAML {
	dplyyaml := DeploymentYAML{}
	content, err := new(ContentReader).LocalReader.ReadLocal(dply)
	Check(err)
	err = dm.UnmarshalDeployment(content, &dplyyaml)
	dplyyaml.Filepath = dply
	return &dplyyaml
}

func (dm *YAMLParser) ComposeActions(manipath string) ([]ActionRecord, error) {
	mani := dm.ParseManifest(manipath)
	var s1 []ActionRecord = make([]ActionRecord, 0)
	for _, action := range mani.Package.Actions {
		wskaction, err := CreateActionFromFile(manipath, action.Location)
		Check(err)

		record := ActionRecord{wskaction, action.Location}
		s1 = append(s1, record)
	}
	return s1, nil
}

func (dm *YAMLParser) ComposeTriggers(mani string) ([]*whisk.Trigger, error) {
	mm := NewYAMLParser()
	manifest := mm.ParseManifest(mani)
	var t1 []*whisk.Trigger = make([]*whisk.Trigger, 0)
	pkg := manifest.Package
	for _, trigger := range pkg.GetTriggerList() {
		wsktrigger := trigger.ComposeWskTrigger()
		t1 = append(t1, wsktrigger)
	}
	return t1, nil
}

func (dm *YAMLParser) ComposeRules(mani string) ([]*whisk.Rule, error) {
	mm := NewYAMLParser()
	manifest := mm.ParseManifest(mani)
	var r1 []*whisk.Rule = make([]*whisk.Rule, 0)
	pkg := manifest.Package
	for _, rule := range pkg.GetRuleList() {
		wskrule := rule.ComposeWskRule()
		r1 = append(r1, wskrule)
	}

	return r1, nil
}

// Is we consider multi pacakge in one yaml?
func (dm *YAMLParser) ComposePackage(manipath string) (*whisk.SentPackageNoPublish, error) {
	mani := dm.ParseManifest(manipath)
	pag := &whisk.SentPackageNoPublish{}
	pag.Name = mani.Package.Packagename
	//The namespace for this package is absent, so we use default guest here.
	pag.Namespace = mani.Package.Namespace
	pag.Publish = false
	return pag, nil
}

func (dm *YAMLParser) Unmarshal(input []byte, manifest *ManifestYAML) error {
	err := yaml.Unmarshal(input, manifest)
	if err != nil {
		log.Fatalf("error happened during unmarshal :%v", err)
		return err
	}
	return nil
}

func (dm *YAMLParser) UnmarshalDeployment(input []byte, deploy *DeploymentYAML) error {
	err := yaml.Unmarshal(input, deploy)
	if err != nil {
		log.Fatalf("error happened during unmarshal :%v", err)
		return err
	}
	return nil
}

func (dm *YAMLParser) Marshal(manifest *ManifestYAML) (output []byte, err error) {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		log.Fatalf("err happened during marshal :%v", err)
		return nil, err
	}
	return data, nil
}

func (dm *YAMLParser) MarshalDeployment(deployment *DeploymentYAML) (output []byte, err error) {
	data, err := yaml.Marshal(deployment)
	if err != nil {
		log.Fatalf("err happened during marshal :%v", err)
		return nil, err
	}
	return data, nil
}
