// dependencies.go
package utils

import (
	"strings"

	"github.com/openwhisk/openwhisk-client-go/whisk"
)

type DependencyRecord struct {
	ProjectPath string
	Packagename string
	Location    string
	Version     string
	Parameters  whisk.KeyValueArr
	Annotations whisk.KeyValueArr
	IsBinding   bool
}

func LocationIsBinding(location string) bool {
	if strings.HasPrefix(location, "/whisk.system") || strings.HasPrefix(location, "whisk.system") {
		return true
	}

	return false
}

func LocationIsGithub(location string) bool {
	if strings.HasPrefix(location, "github.com") || strings.HasPrefix(location, "https://github.com") || strings.HasPrefix(location, "http://github.com") {
		return true
	}

	return false
}
