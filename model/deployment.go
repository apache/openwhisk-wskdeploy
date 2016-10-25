package model

import (
	"gopkg.in/yaml.v2"
	"log"
)

var Deployer *DeploymentManager

func init() {
	Deployer = &DeploymentManager{}
}

type ParseDeploymentYaml interface {
	unmarshal(input []byte, deploy *DeploymentYAML) error
	marshal(deployment *DeploymentYAML) (output []byte, err error)
}

type DeploymentManager struct {
	deployments []*DeploymentYAML
	lastID      uint32
}

type Action struct {
	Name     string                   `yaml:"name"`
	Version  string                   `yaml:"version"`
	Function string                   `yaml:"function"`
	Runtime  string                   `yaml:"runtime"`
	Input    []map[string]interface{} `yaml:"inputs"`
	Output   []map[string]interface{} `yaml:"outputs"`
}

type Package struct {
	Packagename string   `yaml:"packagename"`
	Version     string   `yaml:"version"`
	License     string   `yaml:"license"`
	Actions     []Action `yaml:"actions"`
}

type DeploymentYAML struct{
	Package  Package `yaml:"package"`
}

func (dm *DeploymentManager) Unmarshal(input []byte, deploy *DeploymentYAML) error {
	err := yaml.Unmarshal(input, deploy)
	if err != nil {
		log.Fatalf("error happened during unmarshal :%v", err)
		return err
	}
	return nil
}

func (dm *DeploymentManager) Marshal(deployment *DeploymentYAML) (output []byte, err error) {
	data, err := yaml.Marshal(deployment)
	if err != nil {
		log.Fatalf("err happened during marshal :%v", err)
		return nil, err
	}
	return data, nil
}
