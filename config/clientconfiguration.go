package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
)

var ClientConfig *whisk.Config

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
