package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os/user"

	"bufio"

	"github.com/openwhisk/openwhisk-client-go/whisk"
	"os"
	"reflect"
	"strings"
)

// ActionRecord is a container to keep track of
// a whisk action struct and a location filepath we use to
// map files and manifest declared actions
type ActionRecord struct {
	Action      *whisk.Action
	Packagename string
	Filepath    string
}

type TriggerRecord struct {
	Trigger     *whisk.Trigger
	Packagename string
}

type RuleRecord struct {
	Rule        *whisk.Rule
	Packagename string
}

// Utility to convert hostname to URL object
func GetURLBase(host string) (*url.URL, error) {

	urlBase := fmt.Sprintf("%s/api", host)
	url, err := url.Parse(urlBase)

	if len(url.Scheme) == 0 || len(url.Host) == 0 {
		urlBase = fmt.Sprintf("https://%s/api", host)
		url, err = url.Parse(urlBase)
	}

	return url, err
}

func GetHomeDirectory() string {
	usr, err := user.Current()
	Check(err)

	return usr.HomeDir
}

// Potentially complex structures(such as DeploymentApplication, DeploymentPackage)
// could implement those interface which is convenient for put, get subtract in
// containers etc.
type Comparable interface {
	HashCode() uint32
	Equals() bool
}

func IsFeedAction(trigger *whisk.Trigger) (string, bool) {
	for _, annotation := range trigger.Annotations {
		if annotation.Key == "feed" {
			return annotation.Value.(string), true
		}
	}

	return "", false
}

func IsJSON(s string) (interface{}, bool) {
	var js interface{}
	if json.Unmarshal([]byte(s), &js) == nil {
		return js, true
	}
	return nil, false

}

// Common utilities

// Prompt for user input
func Ask(reader *bufio.Reader, question string, def string) string {
	fmt.Print(question + " (" + def + "): ")
	answer, _ := reader.ReadString('\n')
	len := len(answer)
	if len == 1 {
		return def
	}
	return answer[:len-1]
}

// Get the env variable value by key.
// Get the env variable if the key is start by $
func GetEnvVar(key interface{}) interface{} {
	if reflect.TypeOf(key).String() == "string" {
		if strings.HasPrefix(key.(string), "$") {
			envkey := strings.Split(key.(string), "$")[1]
			value := os.Getenv(envkey)
			if value != "" {
				return value
			}
			return envkey
		}
		return key.(string)
	}
	return key
}
