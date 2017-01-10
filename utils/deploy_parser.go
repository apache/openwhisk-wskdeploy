package utils

import (
	"gopkg.in/yaml.v2"
	"log"
)

func (dm *YAMLParser) UnmarshalDeployment(input []byte, deploy *DeploymentYAML) error {
	err := yaml.Unmarshal(input, deploy)
	if err != nil {
		log.Fatalf("error happened during unmarshal :%v", err)
		return err
	}
	return nil
}

func (dm *YAMLParser) MarshalDeployment(deployment *DeploymentYAML) (output []byte, err error) {
	data, err := yaml.Marshal(deployment)
	if err != nil {
		log.Fatalf("err happened during marshal :%v", err)
		return nil, err
	}
	return data, nil
}

func (dm *YAMLParser) ParseDeployment(dply string) *DeploymentYAML {
	dplyyaml := DeploymentYAML{}
	content, err := new(ContentReader).LocalReader.ReadLocal(dply)
	Check(err)
	err = dm.UnmarshalDeployment(content, &dplyyaml)
	dplyyaml.Filepath = dply
	return &dplyyaml
}

//********************Application functions*************************//
//This is for parse the deployment yaml file.
func (app *Application) GetPackageList() []Package {
	var s1 []Package = make([]Package, 0)
	for pkg_name, pkg := range app.Packages {
		pkg.Packagename = pkg_name
		s1 = append(s1, pkg)
	}
	return s1
}
