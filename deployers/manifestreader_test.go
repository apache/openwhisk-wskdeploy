// +build unit

/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package deployers

import (
	"fmt"
	"testing"

	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/apache/openwhisk-wskdeploy/dependencies"
	"github.com/apache/openwhisk-wskdeploy/parsers"
	"github.com/apache/openwhisk-wskdeploy/runtimes"
	"github.com/apache/openwhisk-wskdeploy/utils"
	"github.com/apache/openwhisk-wskdeploy/wskprint"
	"github.com/stretchr/testify/assert"
)

var mr *ManifestReader
var ps *parsers.YAMLParser
var ms *parsers.YAML

const (
	// local error messages
	TEST_ERROR_BUILD_SERVICE_DEPLOYER        = "Manifest [%s]: Failed to build service deployer."
	TEST_ERROR_MANIFEST_PARSE_FAILURE        = "Manifest [%s]: Failed to parse."
	TEST_ERROR_MANIFEST_SET_PACKAGES         = "Manifest [%s]: Failed to set packages."
	TEST_ERROR_FAILED_TO_REPORT_ERROR        = "Manifest [%s]: Failed to report parser error."
	TEST_ERROR_MANIFEST_SET_ANNOTATION       = "[%s]: Failed to set Annotation value."
	TEST_ERROR_MANIFEST_SET_INPUT_PARAMETER  = "[%s]: Failed to set input Parameter value."
	TEST_ERROR_MANIFEST_SET_PUBLISH          = "Package [%s]: Failed to set publish."
	TEST_ERROR_MANIFEST_SET_DEPENDENCIES     = "Package [%s]: Failed to set dependencies."
	TEST_ERROR_MANIFEST_SET_ACTION_CODE      = "Action [%s]: Failed to set action code."
	TEST_ERROR_MANIFEST_SET_ACTION_KIND      = "Action [%s]: Failed to set action kind."
	TEST_ERROR_MANIFEST_SET_ACTION_WEB       = "Action [%s]: Failed to set web action."
	TEST_ERROR_MANIFEST_SET_ACTION_CONDUCTOR = "Action [%s]: Failed to set conductor action."
	TEST_ERROR_MANIFEST_SET_ACTION_IMAGE     = "Action [%s]: Failed to set action image."
)

func init() {
	// Setup "trace" flag for unit tests based upon "go test" -v flag
	utils.Flags.Trace = wskprint.DetectGoTestVerbose()
}

func buildServiceDeployer(manifestFile string) (*ServiceDeployer, error) {
	deploymentFile := ""
	var deployer = NewServiceDeployer()
	deployer.ManifestPath = manifestFile
	deployer.DeploymentPath = deploymentFile
	deployer.Preview = utils.Flags.Preview

	deployer.DependencyMaster = make(map[string]dependencies.DependencyRecord)

	config := whisk.Config{
		Namespace:        "test",
		AuthToken:        "user:pass",
		Host:             "host",
		ApigwAccessToken: "token",
	}
	deployer.ClientConfig = &config

	op, error := runtimes.ParseOpenWhisk(deployer.ClientConfig.Host)
	if error == nil {
		runtimes.SupportedRunTimes = runtimes.ConvertToMap(op)
		runtimes.DefaultRunTimes = runtimes.DefaultRuntimes(op)
		runtimes.FileExtensionRuntimeKindMap = runtimes.FileExtensionRuntimes(op)
		runtimes.FileRuntimeExtensionsMap = runtimes.FileRuntimeExtensions(op)
	}

	return deployer, nil
}

func TestManifestReader_InitPackages(t *testing.T) {
	manifestFile := "../tests/dat/manifest_validate_package_inputs_and_annotations.yaml"
	deployer, err := buildServiceDeployer(manifestFile)
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_BUILD_SERVICE_DEPLOYER, manifestFile))

	var manifestReader = NewManifestReader(deployer)
	manifestReader.IsUndeploy = false
	manifest, manifestParser, err := manifestReader.ParseManifest()
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_PARSE_FAILURE, manifestFile))

	err = manifestReader.InitPackages(manifestParser, manifest, whisk.KeyValue{})
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_SET_PACKAGES, manifestFile))
	assert.Equal(t, 3, len(deployer.Deployment.Packages), fmt.Sprintf(TEST_ERROR_MANIFEST_SET_PACKAGES, manifestFile))

	expectedParametersAndAnnotations := 0

	for packageName, pack := range deployer.Deployment.Packages {
		switch packageName {
		case "helloworld1":
			expectedParametersAndAnnotations = 0
			assert.False(t, *pack.Package.Publish, fmt.Sprintf(TEST_ERROR_MANIFEST_SET_PUBLISH, packageName))
		case "helloworld2":
			expectedParametersAndAnnotations = 1
			for _, param := range pack.Package.Parameters {
				switch param.Key {
				case "helloworld_input1":
					assert.Equal(t, "value1", param.Value,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_INPUT_PARAMETER, packageName))
				}
			}
			for _, annotation := range pack.Package.Annotations {
				switch annotation.Key {
				case "helloworld_annotation1":
					assert.Equal(t, "value1", annotation.Value,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_ANNOTATION, packageName))
				}
			}
			assert.False(t, *pack.Package.Publish, fmt.Sprintf(TEST_ERROR_MANIFEST_SET_PUBLISH, packageName))
		case "helloworld3":
			expectedParametersAndAnnotations = 2
			for _, param := range pack.Package.Parameters {
				switch param.Key {
				case "helloworld_input1":
					assert.Equal(t, "value1", param.Value,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_INPUT_PARAMETER, packageName))
				case "helloworld_input2":
					assert.Equal(t, "value2", param.Value,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_INPUT_PARAMETER, packageName))
				}
			}
			for _, annotation := range pack.Package.Annotations {
				switch annotation.Key {
				case "helloworld_annotation1":
					assert.Equal(t, "value1", annotation.Value,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_ANNOTATION, packageName))
				case "helloworld_annotation2":
					assert.Equal(t, "value2", annotation.Value,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_ANNOTATION, packageName))
				}
			}
			assert.True(t, *pack.Package.Publish, fmt.Sprintf(TEST_ERROR_MANIFEST_SET_PUBLISH, packageName))
		}
		assert.Equal(t, expectedParametersAndAnnotations, len(pack.Package.Parameters),
			fmt.Sprintf(TEST_ERROR_MANIFEST_SET_INPUT_PARAMETER, packageName))
		assert.Equal(t, expectedParametersAndAnnotations, len(pack.Package.Annotations),
			fmt.Sprintf(TEST_ERROR_MANIFEST_SET_ANNOTATION, packageName))
	}
}

func TestManifestReader_SetDependencies(t *testing.T) {
	manifestFile := "../tests/dat/manifest_validate_dependencies.yaml"
	deployer, err := buildServiceDeployer(manifestFile)
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_BUILD_SERVICE_DEPLOYER, manifestFile))

	var manifestReader = NewManifestReader(deployer)
	manifestReader.IsUndeploy = false
	manifest, manifestParser, err := manifestReader.ParseManifest()
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_PARSE_FAILURE, manifestFile))

	err = manifestReader.InitPackages(manifestParser, manifest, whisk.KeyValue{})
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_SET_PACKAGES, manifestFile))

	err = manifestReader.HandleYaml(manifestParser, manifest, whisk.KeyValue{})
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_PARSE_FAILURE, manifestFile))

	expectedLocationHelloWorlds := "https://github.com/apache/incubator-openwhisk-test/packages/helloworlds"
	expectedLocationHelloWhisk := "https://github.com/apache/incubator-openwhisk-test/packages/hellowhisk"
	expectedLocationUtils := "/whisk.system/utils"

	for pkgName, pkg := range deployer.Deployment.Packages {
		switch pkgName {
		case "helloworld1":
			for depName, dep := range pkg.Dependencies {
				switch depName {
				case "dependency1":
				case "helloworlds":
					assert.Equal(t, expectedLocationHelloWorlds, dep.Location,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_DEPENDENCIES, pkgName))
				case "dependency2":
					assert.Equal(t, expectedLocationHelloWhisk, dep.Location,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_DEPENDENCIES, pkgName))
				case "dependency3":
					assert.Equal(t, expectedLocationUtils, dep.Location,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_DEPENDENCIES, pkgName))
				}
			}

		case "helloworld2":
			for depName, dep := range pkg.Dependencies {
				switch depName {
				case "helloworlds":
				case "dependency1":
				case "dependency4":
					assert.Equal(t, expectedLocationHelloWorlds, dep.Location,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_DEPENDENCIES, pkgName))
				case "dependency5":
					assert.Equal(t, expectedLocationHelloWhisk, dep.Location,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_DEPENDENCIES, pkgName))
				case "dependency6":
					assert.Equal(t, expectedLocationUtils, dep.Location,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_DEPENDENCIES, pkgName))
				}
			}
		}
	}
}

func TestManifestReader_SetDependencies_Bogus(t *testing.T) {
	manifestFile := "../tests/dat/manifest_validate_dependencies_bogus.yaml"
	deployer, err := buildServiceDeployer(manifestFile)
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_BUILD_SERVICE_DEPLOYER, manifestFile))

	var manifestReader = NewManifestReader(deployer)
	manifestReader.IsUndeploy = false
	manifest, manifestParser, err := manifestReader.ParseManifest()
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_PARSE_FAILURE, manifestFile))

	err = manifestReader.InitPackages(manifestParser, manifest, whisk.KeyValue{})
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_SET_PACKAGES, manifestFile))

	err = manifestReader.HandleYaml(manifestParser, manifest, whisk.KeyValue{})
	assert.NotNil(t, err, fmt.Sprintf(TEST_ERROR_FAILED_TO_REPORT_ERROR, manifestFile))
}

func TestManifestReader_SetActions(t *testing.T) {
	manifestFile := "../tests/dat/manifest_validate_action_all.yaml"
	deployer, err := buildServiceDeployer(manifestFile)
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_BUILD_SERVICE_DEPLOYER, manifestFile))

	var manifestReader = NewManifestReader(deployer)
	manifestReader.IsUndeploy = false
	manifest, manifestParser, err := manifestReader.ParseManifest()
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_PARSE_FAILURE, manifestFile))

	err = manifestReader.InitPackages(manifestParser, manifest, whisk.KeyValue{})
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_SET_PACKAGES, manifestFile))

	err = manifestReader.HandleYaml(manifestParser, manifest, whisk.KeyValue{})
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_PARSE_FAILURE, manifestFile))

	expectedRuntime := "nodejs:6"
	expectedImage := "openwhisk/skeleton"

	for actionName, action := range deployer.Deployment.Packages["helloworld"].Actions {
		switch actionName {
		case "helloworld1":
		case "helloworld2":
			assert.NotEmpty(t, action.Action.Exec.Code,
				fmt.Sprintf(TEST_ERROR_MANIFEST_SET_ACTION_CODE, actionName))
			assert.Equal(t, expectedRuntime, action.Action.Exec.Kind,
				fmt.Sprintf(TEST_ERROR_MANIFEST_SET_ACTION_KIND, actionName))
		case "helloworld3":
			for _, param := range action.Action.Parameters {
				switch param.Key {
				case "parameter1":
					assert.Equal(t, "value1", param.Value,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_INPUT_PARAMETER, actionName))
				case "parameter2":
					assert.Equal(t, "value2", param.Value,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_INPUT_PARAMETER, actionName))
				}
			}
			for _, annotation := range action.Action.Annotations {
				switch annotation.Key {
				case "annotation1":
					assert.Equal(t, "value1", annotation.Value,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_ANNOTATION, actionName))
				case "annotation2":
					assert.Equal(t, "value2", annotation.Value,
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_ANNOTATION, actionName))
				}
			}
		case "helloworld4":
			assert.True(t, action.Action.WebAction(),
				fmt.Sprintf(TEST_ERROR_MANIFEST_SET_ACTION_WEB, actionName))
		case "helloworld5":
			for _, annotation := range action.Action.Annotations {
				switch annotation.Key {
				case "conductor":
					assert.True(t, annotation.Value.(bool),
						fmt.Sprintf(TEST_ERROR_MANIFEST_SET_ACTION_CONDUCTOR, actionName))
				}
			}
		case "helloworld6":
			assert.Equal(t, expectedImage, action.Action.Exec.Image,
				fmt.Sprintf(TEST_ERROR_MANIFEST_SET_ACTION_IMAGE, actionName))
		}
	}
}

func TestManifestReader_SetSequences_Bogus(t *testing.T) {
	manifestFile := "../tests/dat/manifest_validate_sequences_bogus.yaml"
	deployer, err := buildServiceDeployer(manifestFile)
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_BUILD_SERVICE_DEPLOYER, manifestFile))

	var manifestReader = NewManifestReader(deployer)
	manifestReader.IsUndeploy = false
	manifest, manifestParser, err := manifestReader.ParseManifest()
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_PARSE_FAILURE, manifestFile))

	err = manifestReader.InitPackages(manifestParser, manifest, whisk.KeyValue{})
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_SET_PACKAGES, manifestFile))

	err = manifestReader.HandleYaml(manifestParser, manifest, whisk.KeyValue{})
	assert.NotNil(t, err, fmt.Sprintf(TEST_ERROR_FAILED_TO_REPORT_ERROR, manifestFile))
}

func TestManifestReader_SetTriggers_Bogus(t *testing.T) {
	manifestFile := "../tests/dat/manifest_validate_triggers_bogus.yaml"
	deployer, err := buildServiceDeployer(manifestFile)
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_BUILD_SERVICE_DEPLOYER, manifestFile))

	var manifestReader = NewManifestReader(deployer)
	manifestReader.IsUndeploy = false
	manifest, manifestParser, err := manifestReader.ParseManifest()
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_PARSE_FAILURE, manifestFile))

	err = manifestReader.InitPackages(manifestParser, manifest, whisk.KeyValue{})
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_SET_PACKAGES, manifestFile))

	err = manifestReader.HandleYaml(manifestParser, manifest, whisk.KeyValue{})
	assert.NotNil(t, err, fmt.Sprintf(TEST_ERROR_FAILED_TO_REPORT_ERROR, manifestFile))
}

func TestManifestReader_SetRules_Bogus(t *testing.T) {
	manifestFile := "../tests/dat/manifest_validate_rules_bogus.yaml"
	deployer, err := buildServiceDeployer(manifestFile)
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_BUILD_SERVICE_DEPLOYER, manifestFile))

	var manifestReader = NewManifestReader(deployer)
	manifestReader.IsUndeploy = false
	manifest, manifestParser, err := manifestReader.ParseManifest()
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_PARSE_FAILURE, manifestFile))

	err = manifestReader.InitPackages(manifestParser, manifest, whisk.KeyValue{})
	assert.Nil(t, err, fmt.Sprintf(TEST_ERROR_MANIFEST_SET_PACKAGES, manifestFile))

	err = manifestReader.HandleYaml(manifestParser, manifest, whisk.KeyValue{})
	assert.NotNil(t, err, fmt.Sprintf(TEST_ERROR_FAILED_TO_REPORT_ERROR, manifestFile))
}
