# Test cases for openwhisk-wskdeploy.

There are two types of test cases located under the 'tests': unit and integration. You can identify them by the first line of each test case file. Unit tests are tagged with "+build unit" and integration tests are tagged with "+build integration". For example, the test file deploymentreader_test.go under tests/src/deployers contains unit test cases, so it is indicated with the unit tag on the top of the file:

```
// +build unit

package tests
...
```

Integration tests are indicated with the integration tag on the top of the file:
```
// +build integration

package tests
...
```

# How to run wskdeploy tests

Before running unit or integration tests, you need to install the package [Testify](https://github.com/stretchr/testify/) before running the test cases:
```
cd $GOPATH
go get -u github.com/stretchr/testify

```

## Running unit tests

After installing Testify, all the unit tests can be run from the main 'incubator-openwhisk-wskdeploy' repository folder using the following command:

```
$ go test -v ./... -tags unit
```

## Running integration tests

After installing Testify, all integration tests can be run rom the main 'incubator-openwhisk-wskdeploy' repository folder using the following command:

```
$ go test -v ./... -tags integration
```

# How to structure the test cases

All integration test cases are put under the folder 'tests/src/integration'.

Unit tests are co-located in the same directory as the package they are testing (by Go convention). The test files would use the same basename as the file that contains the package they are providing tests for, but with the added '_test' postfix to the base file name.

For example, the package 'deployers', which defines a 'DeploymentReader' service, is declared in the file 'deploymentreader.go'; its corresponding unit test is in the file named 'deploymentreader_test.go' in the same directory.

## Grouping tests

Additional unit tests can be grouped into different packages, and they can be put under a subfolder
named after the package under 'tests/src'. The yaml files of manifest and deployment used by unit tests are put under 'tests/data'.

Different integration tests can be grouped into different packages as well, and they can be put under a subfolder named after the package under 'tests/src/integration'. The yaml files of manifest and deployment used by the integration test are put under the same subfolder as the integration test itself. The source file used by the manifest file can be put under the folder 'tests/src/integration/<package>/src'.

# Unit test listing

| Test File | Manifest | Deployment | Description |
| ------| ------ | ------ | ------ |
| [deployers / deploymentreader_test.go](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/deployers/deploymentreader_test.go) | usecases/badyaml/manifest.yaml | usecases/badyaml/deployment.yaml| Tests DeploymentReader service. |
| [deployers / manifestreader_test.go](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/deployers/manifestreader_test.go) | dat/manifest6.yaml | N/A | Tests ManifestReader service |
| [utils / utils_test.go](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/utils/util_test.go) | N/A | dat/deployment.yaml| Tests ContentReader, ReadUrl |
| [cmd / root_test.go](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/cmd/root_test.go) | N/A | N/A | Tests Cobra frameworks "Root" command (i.e., "wskdeploy") and its child commands|
| [parsers / yamlparser_test.go](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/parsers/yamlparser_test.go) | dat/manifest1.yaml, dat/manifest2.yaml, dat/manifest3.yaml, dat/manifest4.yaml, dat/manifest5.yaml, dat/manifest6.yaml | dat/deploy1.yaml, dat/deploy2.yaml, dat/deploy3.yaml, dat/deploy4.yaml | Tests YAML parser against various Manifest and Deployment files. |
<!-- | []() | <manifest> | <depl> | <desc> | -->
