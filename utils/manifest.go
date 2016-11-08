package utils

import (
	"gopkg.in/yaml.v2"
	"log"
)

var Deployer *ManifestManager

// structs that denotes the sample manifest.yaml, wrapped yaml.v2

func init() {
	Deployer = &ManifestManager{}
}

type ParseManifestYaml interface {
	unmarshal(input []byte, deploy *ManifestYAML) error
	marshal(manifest *ManifestYAML) (output []byte, err error)
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
