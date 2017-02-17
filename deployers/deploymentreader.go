package deployers

import (
	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/parsers"
)

type DeploymentReader struct {
	serviceDeployer      *ServiceDeployer
	DeploymentDescriptor *parsers.DeploymentYAML
}

func NewDeploymentReader(serviceDeployer *ServiceDeployer) *DeploymentReader {
	var dep DeploymentReader
	dep.serviceDeployer = serviceDeployer

	return &dep
}

// Wrapper parser to handle yaml dir
func (reader *DeploymentReader) HandleYaml() error {

	dep := reader.serviceDeployer

	deploymentParser := parsers.NewYAMLParser()
	deployment := deploymentParser.ParseDeployment(dep.DeploymentPath)

	reader.DeploymentDescriptor = deployment

	return nil
}

// Update entities with deployment settings
func (reader *DeploymentReader) BindAssets() error {

	reader.bindPackageInputsAndAnnotations()
	reader.bindActionInputsAndAnnotations()
	reader.bindTriggerInputsAndAnnotations()

	return nil

}

func (reader *DeploymentReader) bindPackageInputsAndAnnotations() {

	packArray := make([]parsers.Package, 0)

	if reader.DeploymentDescriptor.Application.Packages == nil {
		packArray = append(packArray, reader.DeploymentDescriptor.Application.Package)
	} else {
		for _, depPacks := range reader.DeploymentDescriptor.Application.Packages {
			packArray = append(packArray, depPacks)
		}
	}

	for _, pack := range packArray {
		serviceDeployPack := reader.serviceDeployer.Deployment.Packages[pack.Packagename]

		var keyValArr whisk.KeyValueArr
		for name, input := range pack.Inputs {
			var keyVal whisk.KeyValue

			keyVal.Key = name
			keyVal.Value = input

			keyValArr = append(keyValArr, keyVal)
		}

		serviceDeployPack.Package.Parameters = keyValArr

		keyValArr = keyValArr[:0]
		for name, input := range pack.Annotations {
			var keyVal whisk.KeyValue

			keyVal.Key = name
			keyVal.Value = input

			keyValArr = append(keyValArr, keyVal)
		}

		serviceDeployPack.Package.Annotations = keyValArr

	}
}

func (reader *DeploymentReader) bindActionInputsAndAnnotations() {

	packArray := make([]parsers.Package, 1)

	if reader.DeploymentDescriptor.Application.Packages == nil {
		packArray = append(packArray, reader.DeploymentDescriptor.Application.Package)
	} else {
		for _, depPacks := range reader.DeploymentDescriptor.Application.Packages {
			packArray = append(packArray, depPacks)
		}
	}

	for _, pack := range packArray {

		for actionName, action := range pack.Actions {

			serviceDeployPack := reader.serviceDeployer.Deployment.Packages[pack.Packagename]

			var keyValArr whisk.KeyValueArr
			for name, input := range action.Inputs {
				var keyVal whisk.KeyValue

				keyVal.Key = name
				keyVal.Value = input

				keyValArr = append(keyValArr, keyVal)
			}

			if wskAction, exists := serviceDeployPack.Actions[actionName]; exists {
				wskAction.Action.Parameters = keyValArr
			}

			keyValArr = keyValArr[:0]

			for name, input := range action.Annotations {
				var keyVal whisk.KeyValue

				keyVal.Key = name
				keyVal.Value = input

				keyValArr = append(keyValArr, keyVal)
			}

			if wskAction, exists := serviceDeployPack.Actions[actionName]; exists {
				wskAction.Action.Annotations = keyValArr
			}
		}

	}
}

func (reader *DeploymentReader) bindTriggerInputsAndAnnotations() {

	packArray := make([]parsers.Package, 1)

	if reader.DeploymentDescriptor.Application.Packages == nil {
		packArray = append(packArray, reader.DeploymentDescriptor.Application.Package)
	} else {
		for _, depPacks := range reader.DeploymentDescriptor.Application.Packages {
			packArray = append(packArray, depPacks)
		}
	}

	for _, pack := range packArray {

		serviceDeployment := reader.serviceDeployer.Deployment

		for triggerName, trigger := range pack.Triggers {

			var keyValArr whisk.KeyValueArr
			for name, input := range trigger.Inputs {
				var keyVal whisk.KeyValue

				keyVal.Key = name
				keyVal.Value = input

				keyValArr = append(keyValArr, keyVal)
			}

			if wskTrigger, exists := serviceDeployment.Triggers[triggerName]; exists {
				wskTrigger.Parameters = keyValArr
			}

			keyValArr = keyValArr[:0]

			for name, input := range trigger.Annotations {
				var keyVal whisk.KeyValue

				keyVal.Key = name
				keyVal.Value = input

				keyValArr = append(keyValArr, keyVal)
			}

			if wskTrigger, exists := serviceDeployment.Triggers[triggerName]; exists {
				wskTrigger.Annotations = keyValArr
			}
		}

	}
}
