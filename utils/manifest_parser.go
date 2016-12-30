package utils

import (
	"github.com/openwhisk/openwhisk-client-go/whisk"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

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

func (dm *YAMLParser) Marshal(manifest *ManifestYAML) (output []byte, err error) {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		log.Fatalf("err happened during marshal :%v", err)
		return nil, err
	}
	return data, nil
}

func (dm *YAMLParser) ParseManifest(mani string) *ManifestYAML {
	mm := NewYAMLParser()
	maniyaml := ManifestYAML{}
	content, err := new(ContentReader).LocalReader.ReadLocal(mani)
	Check(err)
	err = mm.Unmarshal(content, &maniyaml)
	maniyaml.Filepath = mani
	return &maniyaml
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

func (dm *YAMLParser) ComposeTriggers(mani string, deploy string) ([]*whisk.Trigger, error) {
	mm := NewYAMLParser()
	manifest := mm.ParseManifest(mani)
	//parse the deployment yaml to get feed info as annotations for triggers
	deployment := mm.ParseDeployment(deploy)
	//This is to get the pkgs in deployment yaml
	var kva []whisk.KeyValue = make([]whisk.KeyValue, 0)
	for _, dep_pkg := range deployment.Application.GetPackageList() {
		for _, feed := range dep_pkg.GetFeedList() {
			log.Printf("feed name is %v", feed.Name)
			feedFullName := strings.Join([]string{dep_pkg.Packagename, feed.Name}, "/")
			kv := whisk.KeyValue{"feed", feedFullName}
			kva = append(kva, kv)
			for k, v := range feed.Inputs {
				kv := whisk.KeyValue{k, v}
				kva = append(kva, kv)
			}
		}
	}

	//parse the manifest yaml
	var t1 []*whisk.Trigger = make([]*whisk.Trigger, 0)
	pkg := manifest.Package
	for _, trigger := range pkg.GetTriggerList() {
		wsktrigger := trigger.ComposeWskTrigger(kva)
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

func (action *Action) ComposeWskAction(manipath string) (*whisk.Action, error) {
	wskaction, err := CreateActionFromFile(manipath, action.Location)
	Check(err)
	wskaction.Name = action.Name
	wskaction.Version = action.Version
	wskaction.Namespace = action.Namespace
	return wskaction, err
}
