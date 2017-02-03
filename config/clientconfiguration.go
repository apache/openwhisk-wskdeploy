package config

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/parsers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
)

var ClientConfig *whisk.Config

func NewClient(proppath string, deploymentPath string) (*whisk.Client, *whisk.Config) {
	var clientConfig *whisk.Config
	configs, err := LoadConfiguration(proppath)
	utils.Check(err)
	//we need to get Apihost from property file which currently not defined in sample deployment file.
	baseURL, err := utils.GetURLBase(configs[1])
	utils.Check(err)
	if deploymentPath != "" {
		mm := parsers.NewYAMLParser()
		deployment := mm.ParseDeployment(deploymentPath)
		// We get the first package from the sample deployment file.
		pkg := deployment.Application.GetPackageList()[0]
		clientConfig = &whisk.Config{
			AuthToken: pkg.Credential, //Authtoken
			Namespace: pkg.Namespace,  //Namespace
			BaseURL:   baseURL,
			Version:   "v1",
			Insecure:  true,
		}

	} else {
		clientConfig = &whisk.Config{
			AuthToken: configs[2], //Authtoken
			Namespace: configs[0], //Namespace
			BaseURL:   baseURL,
			Version:   "v1",
			Insecure:  true, // true if you want to ignore certificate signing

		}

	}

	// Setup network client
	client, err := whisk.NewClient(http.DefaultClient, clientConfig)
	utils.Check(err)
	return client, clientConfig
}

// Load configuration will load properties from a file
func LoadConfiguration(propPath string) ([]string, error) {
	props, err := ReadProps(propPath)
	utils.Check(err)
	Namespace := props["NAMESPACE"]
	Apihost := props["APIHOST"]
	Authtoken := props["AUTH"]
	return []string{Namespace, Apihost, Authtoken}, nil
}

func ReadProps(path string) (map[string]string, error) {

	props := map[string]string{}

	file, err := os.Open(path)
	if err != nil {
		// If file does not exist, just return props
		fmt.Printf("Unable to read whisk properties file '%s' (file open error: %s); falling back to default properties\n", path, err)
		return props, nil
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	props = map[string]string{}
	for _, line := range lines {
		kv := strings.Split(line, "=")
		if len(kv) != 2 {
			// Invalid format; skip
			continue
		}
		props[kv[0]] = kv[1]
	}

	return props, nil

}
