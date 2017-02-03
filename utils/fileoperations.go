package utils

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/openwhisk/openwhisk-client-go/whisk"
)

func FileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		err = errors.New("File not found.")
		return false
	} else {
		return true
	}
}

func IsDirectory(filePath string) bool {
	f, err := os.Open(filePath)
	Check(err)

	defer f.Close()

	fi, err := f.Stat()
	Check(err)

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	default:
		return false
	}
}

func CreateActionFromFile(manipath, filePath string) (*whisk.Action, error) {
	ext := path.Ext(filePath)
	baseName := path.Base(filePath)
	//check if the file if from local or from web
	//currently only consider http
	islocal := !strings.HasPrefix(filePath, "http")
	name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	action := new(whisk.Action)
	//better refactor this
	if islocal {
		splitmanipath := strings.Split(manipath, string(os.PathSeparator))
		filePath = strings.TrimRight(manipath, splitmanipath[len(splitmanipath)-1]) + filePath
	}
	// process source code files
	if ext == ".swift" || ext == ".js" || ext == ".py" {

		kind := "nodejs:default"

		switch ext {
		case ".swift":
			kind = "swift:default"
		case ".js":
			kind = "nodejs:default"
		case ".py":
			kind = "python"
		}

		var dat []byte
		var err error

		if islocal {
			dat, err = new(ContentReader).LocalReader.ReadLocal(filePath)
		} else {
			dat, err = new(ContentReader).URLReader.ReadUrl(filePath)
		}

		Check(err)
		action.Exec = new(whisk.Exec)
		action.Exec.Code = string(dat)
		action.Exec = new(whisk.Exec)
		action.Exec.Code = string(dat)
		action.Exec.Kind = kind
		action.Name = name
		action.Publish = false
		return action, nil
		//dat, err := new(ContentReader).URLReader.ReadUrl(filePath)
		//Check(err)

	}
	// If the action is not supported, we better to return an error.
	return nil, errors.New("Unsupported action type.")
}

func ReadProps(path string) (map[string]string, error) {

	props := map[string]string{}

	file, err := os.Open(path)
	if err != nil {
		// If file does not exist, just return props
		log.Printf("Unable to read whisk properties file '%s' (file open error: %s); falling back to default properties\n", path, err)
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

// Load configuration will load properties from a file
func LoadConfiguration(propPath string) ([]string, error) {
	props, err := ReadProps(propPath)
	Check(err)
	Namespace := props["NAMESPACE"]
	Apihost := props["APIHOST"]
	Authtoken := props["AUTH"]
	return []string{Namespace, Apihost, Authtoken}, nil
}
