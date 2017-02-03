package deployers

import (
	"fmt"
	"log"
	"strings"

	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/config"
	"github.com/openwhisk/openwhisk-wskdeploy/parsers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
)

type ManifestReader struct {
	serviceDeployer *ServiceDeployer
}

func NewManfiestReader(serviceDeployer *ServiceDeployer) *ManifestReader {
	var dep ManifestReader
	dep.serviceDeployer = serviceDeployer

	return &dep
}

// Wrapper parser to handle yaml dir
func (deployer *ManifestReader) HandleYaml() error {

	dep := deployer.serviceDeployer

	manifestParser := parsers.NewYAMLParser()
	manifest := manifestParser.ParseManifest(dep.ManifestPath)

	packg, err := manifestParser.ComposePackage(manifest)

	utils.Check(err)

	actions, err := manifestParser.ComposeActions(manifest, dep.ManifestPath)
	utils.Check(err)

	sequences, err := manifestParser.ComposeSequences(dep.ClientConfig.Namespace, packg.Name, manifest)
	utils.Check(err)

	triggers, err := manifestParser.ComposeTriggers(manifest)
	utils.Check(err)

	rules, err := manifestParser.ComposeRules(manifest)
	utils.Check(err)

	if !deployer.SetPackage(packg) {
		log.Panicln("Cannot assign package " + packg.Name)
	}
	if !deployer.SetActions(packg.Name, actions) {
		log.Panicln("duplication founded during deploy actions")
	}
	if !deployer.SetSequences(packg.Name, sequences) {
		log.Panicln("duplication founded during deploy actions")
	}
	if !deployer.SetTriggers(packg.Name, triggers) {
		log.Panicln("duplication founded during deploy triggers")
	}
	if !deployer.SetRules(packg.Name, rules) {
		log.Panicln("duplication founded during deploy rules")
	}

	//deployer.createPackage(packg)

	return nil
}

func (reader *ManifestReader) SetPackage(pkg *whisk.SentPackageNoPublish) bool {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()
	depPkg, exist := dep.Deployment.Packages[pkg.Name]
	if exist {
		if dep.IsDefault == true {
			existPkg := depPkg.Package
			existPkg.Annotations = pkg.Annotations
			existPkg.Namespace = pkg.Namespace
			existPkg.Parameters = pkg.Parameters
			existPkg.Publish = pkg.Publish
			existPkg.Version = pkg.Version

			dep.Deployment.Packages[pkg.Name].Package = existPkg
			return true
		} else {
			return false
		}
	}

	newPack := NewDeploymentPackage()
	newPack.Package = pkg
	dep.Deployment.Packages[pkg.Name] = newPack
	return true
}

func (reader *ManifestReader) SetActions(packageName string, actions []utils.ActionRecord) bool {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, action := range actions {
		//fmt.Println(action.Action.Name)
		existAction, exist := dep.Deployment.Packages[packageName].Actions[action.Action.Name]

		if exist {
			if dep.IsDefault == true {
				// look for actions declared in filesystem default as well as manifest
				// if one exists, merge if they are the same (either same Filepath or manifest doesn't specify a Filepath)
				// if they are not the same log error
				if action.Filepath != "" {
					if existAction.Filepath != action.Filepath {
						log.Printf("Action %s has location %s in manifest but already exists at %s", action.Action.Name, action.Filepath, existAction)
						return false
					} else {
						// merge the two, overwrite existing action with manifest values
						existAction.Action.Annotations = action.Action.Annotations
						existAction.Action.Exec.Kind = action.Action.Exec.Kind
						existAction.Action.Limits = action.Action.Limits
						existAction.Action.Namespace = action.Action.Namespace
						existAction.Action.Parameters = action.Action.Parameters
						existAction.Action.Publish = action.Action.Publish
						existAction.Action.Version = action.Action.Version

						existAction.Packagename = packageName

						dep.Deployment.Packages[packageName].Actions[action.Action.Name] = existAction

						return true
					}
				}
			} else {
				// no defaults, so assume everything is in the incoming ActionRecord
				// return false since it means the action is declared twice in the manifest
				log.Printf("Action %s is declared more than once", action.Action.Name)
				return false
			}
		}
		// doesn't exist so just add to deployer actions
		action.Packagename = packageName

		dep.Deployment.Packages[packageName].Actions[action.Action.Name] = action
	}
	return true
}

func (reader *ManifestReader) SetSequences(packageName string, actions []utils.ActionRecord) bool {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, action := range actions {
		//fmt.Println(action.Action.Name)
		existAction, exist := dep.Deployment.Packages[packageName].Sequences[action.Action.Name]

		if exist {
			if dep.IsDefault == true {
				// look for actions declared in filesystem default as well as manifest
				// if one exists, merge if they are the same (either same Filepath or manifest doesn't specify a Filepath)
				// if they are not the same log error
				if action.Filepath != "" {
					if existAction.Filepath != action.Filepath {
						log.Printf("Action %s has location %s in manifest but already exists at %s", action.Action.Name, action.Filepath, existAction)
						return false
					} else {
						// merge the two, overwrite existing action with manifest values
						existAction.Action.Annotations = action.Action.Annotations
						existAction.Action.Exec.Kind = action.Action.Exec.Kind
						existAction.Action.Limits = action.Action.Limits
						existAction.Action.Namespace = action.Action.Namespace
						existAction.Action.Parameters = action.Action.Parameters
						existAction.Action.Publish = action.Action.Publish
						existAction.Action.Version = action.Action.Version

						existAction.Packagename = packageName

						dep.Deployment.Packages[packageName].Sequences[action.Action.Name] = existAction

						return true
					}
				}
			} else {
				// no defaults, so assume everything is in the incoming ActionRecord
				// return false since it means the action is declared twice in the manifest
				log.Printf("Action %s is declared more than once", action.Action.Name)
				return false
			}
		}
		// doesn't exist so just add to deployer sequences
		action.Packagename = packageName
		dep.Deployment.Packages[packageName].Sequences[action.Action.Name] = action

	}
	return true
}

func (reader *ManifestReader) SetTriggers(packageName string, triggers []*whisk.Trigger) bool {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, trigger := range triggers {
		existTrigger, exist := dep.Deployment.Packages[packageName].Triggers[trigger.Name]
		if exist {
			existTrigger.Name = trigger.Name
			existTrigger.ActivationId = trigger.ActivationId
			existTrigger.Namespace = trigger.Namespace
			existTrigger.Annotations = trigger.Annotations
			existTrigger.Version = trigger.Version
			existTrigger.Parameters = trigger.Parameters
			existTrigger.Publish = trigger.Publish
		} else {
			dep.Deployment.Packages[packageName].Triggers[trigger.Name] = trigger
		}

	}
	return true
}

func (reader *ManifestReader) SetRules(packageName string, rules []*whisk.Rule) bool {
	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, rule := range rules {
		existRule, exist := dep.Deployment.Packages[packageName].Rules[rule.Name]
		if exist {
			existRule.Name = rule.Name
			existRule.Publish = rule.Publish
			existRule.Version = rule.Version
			existRule.Namespace = rule.Namespace
			existRule.Action = rule.Action
			existRule.Trigger = rule.Trigger
			existRule.Status = rule.Status
		} else {
			dep.Deployment.Packages[packageName].Rules[rule.Name] = rule
		}

	}
	return true
}

// from whisk go client
func (deployer *ManifestReader) getQualifiedName(name string, namespace string) string {
	if strings.HasPrefix(name, "/") {
		return name
	} else if strings.HasPrefix(namespace, "/") {
		return fmt.Sprintf("%s/%s", namespace, name)
	} else {
		if len(namespace) == 0 {
			namespace = config.ClientConfig.Namespace
		}
		return fmt.Sprintf("/%s/%s", namespace, name)
	}
}
