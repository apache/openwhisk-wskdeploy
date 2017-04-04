package cmdImp

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/deployers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
)

func Undeploy(params DeployParams) error {
	// TODO: Work your own magic here
	whisk.SetVerbose(params.Verbose)

	projectPath, err := filepath.Abs(params.ProjectPath)
	utils.Check(err)

	if params.ManifestPath != "" {
		if ok, _ := regexp.Match(deployers.ManifestFileNameYml, []byte(params.ManifestPath)); ok {
			params.ManifestPath = path.Join(params.ProjectPath, deployers.ManifestFileNameYml)
		} else {
			params.ManifestPath = path.Join(params.ProjectPath, deployers.ManifestFileNameYaml)
		}
	} else {
		if utils.FileExists(path.Join(projectPath, "manifest.yaml")) {
			params.ManifestPath = path.Join(projectPath, deployers.ManifestFileNameYaml)
		} else {
			params.ManifestPath = path.Join(projectPath, deployers.ManifestFileNameYml)
		}
	}

	if params.DeploymentPath == "" {
		if ok, _ := regexp.Match(deployers.DeploymentFileNameYml, []byte(params.ManifestPath)); ok {
			params.DeploymentPath = path.Join(params.ProjectPath, deployers.DeploymentFileNameYml)
		} else {
			params.DeploymentPath = path.Join(params.ProjectPath, deployers.DeploymentFileNameYaml)
		}

	} else {
		if _, err := os.Stat(path.Join(projectPath, "deployment.yaml")); err == nil {
			params.DeploymentPath = path.Join(projectPath, deployers.DeploymentFileNameYaml)
		} else if _, err := os.Stat(path.Join(projectPath, "deployment.yml")); err == nil {
			params.DeploymentPath = path.Join(projectPath, deployers.DeploymentFileNameYml)
		}
	}

	if utils.FileExists(params.ManifestPath) {

		var deployer = deployers.NewServiceDeployer()
		deployer.ProjectPath = params.ProjectPath
		deployer.ManifestPath = params.ManifestPath
		deployer.DeploymentPath = params.DeploymentPath

		deployer.IsInteractive = params.UseInteractive
		deployer.IsDefault = params.UseDefaults

		userHome := utils.GetHomeDirectory()
		propPath := path.Join(userHome, ".wskprops")

		whiskClient, clientConfig := deployers.NewWhiskClient(propPath, params.DeploymentPath, deployer.IsInteractive)
		deployer.Client = whiskClient
		deployer.ClientConfig = clientConfig

		verifiedPlan, err := deployer.ConstructUnDeploymentPlan()
		err = deployer.UnDeploy(verifiedPlan)
		if err != nil {
			utils.Check(err)
			return err
		} else {
			return nil
		}

	} else {
		log.Println("missing manifest file")
		return errors.New("missing manifest file")
	}
}
