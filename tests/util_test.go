package tests

import (
	"github.com/openwhisk/wskdeploy/utils"
	"testing"
)

var contentReader = new(utils.ContentReader)
var filepath = "../tests/dat/deployment.yaml"

func TestLocalReader_ReadLocal(t *testing.T) {
	b, err := contentReader.ReadLocal(filepath)
	utils.Check(err)
	if b == nil {
		t.Error("get local centent failed")
	}
}

func TestURLReader_ReadUrl(t *testing.T) {
	var exampleUrl = "http://www.baidu.com"
	b, err := contentReader.ReadUrl(exampleUrl)
	utils.Check(err)
	if b == nil {
		t.Error("get web content failed")
	}
}
