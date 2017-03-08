package parsers

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
	"gopkg.in/yaml.v2"
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
	content, err := new(utils.ContentReader).LocalReader.ReadLocal(mani)
	utils.Check(err)
	err = mm.Unmarshal(content, &maniyaml)
	utils.Check(err)
	maniyaml.Filepath = mani
	return &maniyaml
}

// Is we consider multi pacakge in one yaml?
func (dm *YAMLParser) ComposePackage(mani *ManifestYAML) (*whisk.Package, error) {
	//mani := dm.ParseManifest(manipath)
	pag := &whisk.Package{}
	pag.Name = mani.Package.Packagename
	//The namespace for this package is absent, so we use default guest here.
	pag.Namespace = mani.Package.Namespace
	pub := false
	pag.Publish = &pub

	keyValArr := make(whisk.KeyValueArr, 0)
	for name, value := range mani.Package.Inputs {
		var keyVal whisk.KeyValue
		keyVal.Key = name
		keyVal.Value = value

		keyValArr = append(keyValArr, keyVal)
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

			if !strings.HasPrefix(act, mani.Package.Packagename+"/") {
				act = path.Join(mani.Package.Packagename, act)
			}
			components = append(components, path.Join("/"+namespace, act))
		}

		wskaction.Exec.Components = components
		wskaction.Name = key
		pub := false
		wskaction.Publish = &pub
		wskaction.Namespace = namespace

		record := utils.ActionRecord{wskaction, mani.Package.Packagename, key}
		s1 = append(s1, record)
	}
	return s1, nil
}

func (dm *YAMLParser) ComposeActions(mani *ManifestYAML, manipath string) ([]utils.ActionRecord, error) {

	var s1 []utils.ActionRecord = make([]utils.ActionRecord, 0)

	for key, action := range mani.Package.Actions {
		splitmanipath := strings.Split(manipath, string(os.PathSeparator))

		wskaction := new(whisk.Action)
		wskaction.Exec = new(whisk.Exec)

		if action.Location != "" {
			filePath := strings.TrimRight(manipath, splitmanipath[len(splitmanipath)-1]) + action.Location
			action.Location = filePath
			dat, err := new(utils.ContentReader).LocalReader.ReadLocal(filePath)
			utils.Check(err)
			code := string(dat)
			wskaction.Exec.Code = &code

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
		}

		if action.Runtime != "" {
			wskaction.Exec.Kind = action.Runtime
		}

		keyValArr := make(whisk.KeyValueArr, 0)
		for name, value := range action.Inputs {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value = value

			keyValArr = append(keyValArr, keyVal)
		}

		if len(keyValArr) > 0 {
			wskaction.Parameters = keyValArr
		}

		wskaction.Name = key
		pub := false
		wskaction.Publish = &pub

		record := utils.ActionRecord{wskaction, mani.Package.Packagename, action.Location}
		s1 = append(s1, record)
	}

	return s1, nil

}

func (dm *YAMLParser) ComposeTriggers(manifest *ManifestYAML) ([]*whisk.Trigger, error) {

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
		for name, value := range trigger.Inputs {
			var keyVal whisk.KeyValue
			keyVal.Key = name
			keyVal.Value = value

			keyValArr = append(keyValArr, keyVal)
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
