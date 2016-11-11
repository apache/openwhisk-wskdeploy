package utils

import (
	"github.com/openwhisk/openwhisk-client-go/whisk"
	"gopkg.in/yaml.v2"
	"log"
)

// structs that denotes the sample manifest.yaml, wrapped yaml.v2
func NewManifestManager() *ManifestManager {
	return &ManifestManager{}
}

type ParseManifestYaml interface {
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

type ManifestManager struct {
	manifests []*ManifestYAML
	lastID    uint32
}

type Action struct {
	Version  string                 `yaml:"version"`
	Location string                 `yaml:"location"`
	Runtime  string                 `yaml:"runtime"`
	Inputs   map[string]interface{} `yaml:"inputs"`
	Outputs  map[string]interface{} `yaml:"outputs"`
}

type Trigger struct {
	Name    string                 `yaml:"name"`
	Inputs  map[string]interface{} `yaml:"inputs"`
	Outputs map[string]interface{} `yaml:"outputs"`
}

type Rule struct {
	Name        string `yaml:"name"`
	Triggername string `yaml:"name"`
	Actionname  string `yaml:"name"`
}

type Package struct {
	Packagename string             `yaml:"name"`
	Version     string             `yaml:"version"`
	License     string             `yaml:"license"`
	Actions     map[string]Action  `yaml:"actions"`
	Triggers    map[string]Trigger `yaml:"triggers"`
	Rules       map[string]Rule    `yaml:"rules"`
}

type ManifestYAML struct {
	Package Package `yaml:"package"`
}

func readAndParse(mani string) *ManifestYAML {
	mm := NewManifestManager()
	maniyaml := ManifestYAML{}
	content, err := new(ContentReader).LocalReader.ReadLocal(mani)
	Check(err)
	err = mm.Unmarshal(content, &maniyaml)
	return &maniyaml
}

func (dm *ManifestManager) ComposeActions(manipath string) ([]ActionRecord, error) {
	mani := readAndParse(manipath)
	var s1 []ActionRecord = make([]ActionRecord, 0)
	for _, action := range mani.Package.Actions {
		wskaction, err := CreateActionFromFile(manipath, action.Location)
		Check(err)

		record := ActionRecord{wskaction, action.Location}
		s1 = append(s1, record)
	}
	return s1, nil
}

func (dm *ManifestManager) ComposeTriggers(mani string) ([]*whisk.Trigger, error) {
	return nil, nil
}

// Is we consider multi pacakge in one yaml?
func (dm *ManifestManager) ComposePackage(manipath string) (*whisk.SentPackageNoPublish, error) {
	mani := readAndParse(manipath)
	pag := &whisk.SentPackageNoPublish{}
	pag.Name = mani.Package.Packagename
	//The namespace for this package is absent, so we use default guest here.
	pag.Namespace = "guest"
	pag.Publish = false
	return pag, nil
}

func (dm *ManifestManager) Unmarshal(input []byte, deploy *ManifestYAML) error {
	err := yaml.Unmarshal(input, deploy)
	if err != nil {
		log.Fatalf("error happened during unmarshal :%v", err)
		return err
	}
	return nil
}

func (dm *ManifestManager) Marshal(manifest *ManifestYAML) (output []byte, err error) {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		log.Fatalf("err happened during marshal :%v", err)
		return nil, err
	}
	return data, nil
}
