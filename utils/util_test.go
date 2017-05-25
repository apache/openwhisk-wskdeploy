// +build unit

package utils

import (
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
)

var contentReader = new(ContentReader)
var testfile = "../tests/dat/deployment.yaml"

func TestLocalReader_ReadLocal(t *testing.T) {
	b, err := contentReader.ReadLocal(testfile)
	Check(err)
	if b == nil {
		t.Error("get local centent failed")
	}
}

func TestURLReader_ReadUrl(t *testing.T) {
	var exampleUrl = "http://www.baidu.com"
	b, err := contentReader.ReadUrl(exampleUrl)
	Check(err)
	if b == nil {
		t.Error("get web content failed")
	}
}

// The dollar sign test cases.
func TestGetEnvVar(t *testing.T) {
	os.Setenv("NoDollar", "NO dollar")
	os.Setenv("WithDollar", "oh, dollars!")
	os.Setenv("5000", "5000")
	fmt.Println(GetEnvVar("NoDollar"))
	fmt.Println(GetEnvVar("$WithDollar"))
	fmt.Println(GetEnvVar("$5000"))
	assert.Equal(t, "NoDollar", GetEnvVar("NoDollar"), "NoDollar should be no change.")
	assert.Equal(t, "oh, dollars!", GetEnvVar("$WithDollar"), "dollar sign should be handled")
	assert.Equal(t, "5000", GetEnvVar("5000"), "Should be no difference between integer and string")
	assert.Equal(t, "WithDollarAgain", GetEnvVar("$WithDollarAgain"), "if not found, just return the env")
}
