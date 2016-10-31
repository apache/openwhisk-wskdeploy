package tests

import (
	"testing"
	"github.com/openwhisk/wsktool/cmd"
)

var contentReader = new(cmd.ContentReader{new(cmd.UrlContent), new(cmd.LocalContent)})
var filepath = "../tests/dat/deployment.yaml"

func TestLocalReader_ReadLocal(t *testing.T) {
	b, err := contentReader.ReadLocal(filepath)
	cmd.Check(err)
	if b == nil {
		t.Error("get local centent failed")
	}
}

func TestURLReader_ReadUrl(t *testing.T) {
	var exampleUrl = "example.com"
	b, err := contentReader.ReadUrl(exampleUrl)
	cmd.Check(err)
	if b == nil {
		t.Error("get web content failed")
	}
}
