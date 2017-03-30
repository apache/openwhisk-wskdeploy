package cmdImp

import (
	"errors"
	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-client-go/wski18n"
	"github.com/openwhisk/openwhisk-wskdeploy/deployers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
	"log"
	"path"
	"path/filepath"
	"regexp"
)

type DeployParams struct {
	Verbose        bool
	ProjectPath    string
	ManifestPath   string
	DeploymentPath string
	UseDefaults    bool
	UseInteractive bool
}

func Deploy(params DeployParams) error {

	whisk.SetVerbose(params.Verbose)

	projectPath, err := filepath.Abs(params.ProjectPath)
	utils.Check(err)

	if params.ManifestPath == "" {
		if ok, _ := regexp.Match(ManifestFileNameYml, []byte(params.ManifestPath)); ok {
			params.ManifestPath = path.Join(projectPath, ManifestFileNameYml)
		} else {
			params.ManifestPath = path.Join(projectPath, ManifestFileNameYaml)
		}

	}

	if params.DeploymentPath == "" {
		if ok, _ := regexp.Match(DeploymentFileNameYml, []byte(params.ManifestPath)); ok {
			params.DeploymentPath = path.Join(projectPath, DeploymentFileNameYml)
		} else {
			params.DeploymentPath = path.Join(projectPath, DeploymentFileNameYaml)
		}

	}

	if utils.MayExists(params.ManifestPath) {

		var deployer = deployers.NewServiceDeployer()
		deployer.ProjectPath = projectPath
		deployer.ManifestPath = params.ManifestPath
		deployer.DeploymentPath = params.DeploymentPath
		// perform some quick check here.
		go func() {
			deployer.Check()
		}()
		deployer.IsDefault = params.UseDefaults

		deployer.IsInteractive = params.UseInteractive

		propPath := ""
		if !utils.Flags.WithinOpenWhisk {
			userHome := utils.GetHomeDirectory()
			propPath = path.Join(userHome, ".wskprops")
		}
		whiskClient, clientConfig := deployers.NewWhiskClient(propPath, params.DeploymentPath, deployer.IsInteractive)
		deployer.Client = whiskClient
		deployer.ClientConfig = clientConfig

		err := deployer.ConstructDeploymentPlan()

		if err != nil {
			utils.Check(err)
			return err
		}

		err = deployer.Deploy()
		if err != nil {
			utils.Check(err)
			return err
		} else {
			return nil
		}

	} else {
		if utils.Flags.WithinOpenWhisk {
			utils.PrintOpenWhiskError(wski18n.T("missing manifest.yaml file"))
			return errors.New("missing manifest.yaml file")
		} else {
			log.Println("missing manifest.yaml file")
			return errors.New("missing manifest.yaml file")
		}
	}

}
