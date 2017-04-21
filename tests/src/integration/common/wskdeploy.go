package common

import (
	"os"
	"os/exec"
)

const cmd = "wskdeploy"

type Wskdeploy struct {
	Path string
	Dir  string
}

func NewWskdeploy() *Wskdeploy {
	return NewWskWithPath(os.Getenv("GOPATH") + "/src/github.com/openwhisk/openwhisk-wskdeploy/")
}

func NewWskWithPath(path string) *Wskdeploy {
	var dep Wskdeploy
	dep.Path = cmd
	dep.Dir = path
	return &dep
}

func (wskdeploy *Wskdeploy) RunCommand(s ...string) ([]byte, error) {
	command := exec.Command(wskdeploy.Path, s...)
	command.Dir = wskdeploy.Dir
	return command.CombinedOutput()
}

func (wskdeploy *Wskdeploy) Deploy(manifestPath string, deploymentPath string) ([]byte, error) {
	return wskdeploy.RunCommand("-m", manifestPath, "-d", deploymentPath)
}

func (wskdeploy *Wskdeploy) Undeploy(manifestPath string, deploymentPath string) ([]byte, error) {
	return wskdeploy.RunCommand("undeploy", "-m", manifestPath, "-d", deploymentPath)
}
