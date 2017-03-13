package deployers

import (
	"errors"
	"fmt"
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

func (deployer *ManifestReader) ParseManifest() (*parsers.ManifestYAML, *parsers.YAMLParser, error) {
	dep := deployer.serviceDeployer
	manifestParser := parsers.NewYAMLParser()
	manifest := manifestParser.ParseManifest(dep.ManifestPath)

	return manifest, manifestParser, nil
}

func (reader *ManifestReader) InitRootPackage(manifestParser *parsers.YAMLParser, manifest *parsers.ManifestYAML) error {
	packg, err := manifestParser.ComposePackage(manifest)

	utils.Check(err)
	reader.SetPackage(packg)

	return nil
}

func (read *ManifestReader) InitRootBindingPackage(manifestParser *parsers.YAMLParser, manifest *parsers.ManifestYAML) error {
	binpackg, err := manifestParser.ComposeBindingPackage(manifest)
	utils.Check(err)
	read.SetBindingPackage(binpackg)
	return nil
}

// Wrapper parser to handle yaml dir
func (deployer *ManifestReader) HandleYaml(manifestParser *parsers.YAMLParser, manifest *parsers.ManifestYAML) error {

	actions, err := manifestParser.ComposeActions(manifest, deployer.serviceDeployer.ManifestPath)
	utils.Check(err)

	sequences, err := manifestParser.ComposeSequences(deployer.serviceDeployer.ClientConfig.Namespace, manifest)
	utils.Check(err)

	triggers, err := manifestParser.ComposeTriggers(manifest)
	utils.Check(err)

	rules, err := manifestParser.ComposeRules(manifest)
	utils.Check(err)

	err = deployer.SetActions(actions)
	utils.Check(err)

	err = deployer.SetSequences(sequences)
	utils.Check(err)

	err = deployer.SetTriggers(triggers)
	utils.Check(err)

	err = deployer.SetRules(rules)
	utils.Check(err)

	return nil
}

func (reader *ManifestReader) SetPackage(pkg *whisk.Package) error {

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
			return nil
		} else {
			return errors.New("Package " + pkg.Name + "exists twice")
		}
	}

	newPack := NewDeploymentPackage()
	newPack.Package = pkg
	dep.Deployment.Packages[pkg.Name] = newPack
	return nil
}

func (reader *ManifestReader) SetBindingPackage(bpkg *whisk.BindingPackage) error {
	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()
	depbPkg, exist := dep.Deployment.BindingPackages[bpkg.Name]
	if exist {
		if dep.IsDefault == true {
			existPkg := depbPkg.BindingPackage
			existPkg.Annotations = bpkg.Annotations
			existPkg.Namespace = bpkg.Namespace
			existPkg.Parameters = bpkg.Parameters
			existPkg.Publish = bpkg.Publish
			existPkg.Version = bpkg.Version

			dep.Deployment.BindingPackages[bpkg.Name].BindingPackage = existPkg
			return nil
		} else {
			return errors.New("BindingPackage " + bpkg.Name + "exists twice")
		}
	}

	newPack := NewDeploymentBindingPackage()
	newPack.BindingPackage = bpkg
	dep.Deployment.BindingPackages[bpkg.Name] = newPack
	return nil
}

func (reader *ManifestReader) SetActions(actions []utils.ActionRecord) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, manifestAction := range actions {
		existAction, exists := reader.serviceDeployer.Deployment.Packages[manifestAction.Packagename].Actions[manifestAction.Action.Name]

		if exists == true {
			if existAction.Filepath == manifestAction.Filepath || manifestAction.Filepath == "" {
				// we're adding a filesystem detected action so just updated code and filepath if needed
				if manifestAction.Action.Exec.Kind != "" {
					existAction.Action.Exec.Kind = manifestAction.Action.Exec.Kind
				}

				if manifestAction.Action.Exec.Code != nil {
					code := *manifestAction.Action.Exec.Code
					if code != "" {
						existAction.Action.Exec.Code = manifestAction.Action.Exec.Code
					}
				}

				existAction.Action.Annotations = manifestAction.Action.Annotations
				existAction.Action.Limits = manifestAction.Action.Limits
				existAction.Action.Parameters = manifestAction.Action.Parameters
				existAction.Action.Version = manifestAction.Action.Version

				if manifestAction.Filepath != "" {
					existAction.Filepath = manifestAction.Filepath
				}

				err := reader.checkAction(existAction)
				utils.Check(err)

			} else {
				// Action exists, but references two different sources
				return errors.New("manifestReader. Error: Conflict detected for action named " + existAction.Action.Name + ". Found two locations for source file: " + existAction.Filepath + " and " + manifestAction.Filepath)
			}
		} else {
			// not a new action so to actions in package

			err := reader.checkAction(manifestAction)
			utils.Check(err)
			reader.serviceDeployer.Deployment.Packages[manifestAction.Packagename].Actions[manifestAction.Action.Name] = manifestAction
		}
	}

	return nil
}

func (reader *ManifestReader) checkAction(action utils.ActionRecord) error {
	if action.Filepath == "" {
		return errors.New("Error: Action " + action.Action.Name + " has no source code location set.")
	}

	if action.Action.Exec.Kind == "" {
		return errors.New("Error: Action " + action.Action.Name + " has no kind set")
	}

	if action.Action.Exec.Code != nil {
		code := *action.Action.Exec.Code
		if code == "" && action.Action.Exec.Kind != "sequence" {
			return errors.New("Error: Action " + action.Action.Name + " has no source code")
		}
	}

	return nil
}

func (reader *ManifestReader) SetSequences(actions []utils.ActionRecord) error {
	return reader.SetActions(actions)
}

func (reader *ManifestReader) SetTriggers(triggers []*whisk.Trigger) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, trigger := range triggers {
		existTrigger, exist := dep.Deployment.Triggers[trigger.Name]
		if exist {
			existTrigger.Name = trigger.Name
			existTrigger.ActivationId = trigger.ActivationId
			existTrigger.Namespace = trigger.Namespace
			existTrigger.Annotations = trigger.Annotations
			existTrigger.Version = trigger.Version
			existTrigger.Parameters = trigger.Parameters
			existTrigger.Publish = trigger.Publish
		} else {
			dep.Deployment.Triggers[trigger.Name] = trigger
		}

	}
	return nil
}

func (reader *ManifestReader) SetRules(rules []*whisk.Rule) error {
	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, rule := range rules {
		existRule, exist := dep.Deployment.Rules[rule.Name]
		if exist {
			existRule.Name = rule.Name
			existRule.Publish = rule.Publish
			existRule.Version = rule.Version
			existRule.Namespace = rule.Namespace
			existRule.Action = rule.Action
			existRule.Trigger = rule.Trigger
			existRule.Status = rule.Status
		} else {
			dep.Deployment.Rules[rule.Name] = rule
		}

	}
	return nil
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
