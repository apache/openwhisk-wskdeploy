package deployers

import (
	"net/http"
	"net/url"

	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/parsers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
)

func NewWhiskClient(proppath string, deploymentPath string) (*whisk.Client, *whisk.Config) {
	var clientConfig *whisk.Config
	configs, err := utils.LoadConfiguration(proppath)
	utils.Check(err)

	credential := configs[2]
	namespace := configs[0]
	//we need to get Apihost from property file which currently not defined in sample deployment file.
	baseURL, err := utils.GetURLBase(configs[1])
	utils.Check(err)

	if deploymentPath != "" {
		mm := parsers.NewYAMLParser()
		deployment := mm.ParseDeployment(deploymentPath)
		// We get the first package from the sample deployment file.
		credentialDep := deployment.Application.Credential
		namespaceDep := deployment.Application.Namespace
		baseUrlDep := deployment.Application.BaseUrl

		if credentialDep != "" {
			credential = credentialDep
		}

		if namespaceDep != "" {
			namespace = namespaceDep
		}

		if baseUrlDep != "" {
			u, err := url.Parse(baseUrlDep)
			utils.Check(err)

			baseURL = u
		}
	}

	clientConfig = &whisk.Config{
		AuthToken: credential, //Authtoken
		Namespace: namespace,  //Namespace
		BaseURL:   baseURL,
		Version:   "v1",
		Insecure:  true, // true if you want to ignore certificate signing

	}

	// Setup network client
	client, err := whisk.NewClient(http.DefaultClient, clientConfig)
	utils.Check(err)
	return client, clientConfig

}
