# Test Cases for wskdeploy.

There are two types of test cases supported (1) unit test and (2) integration test.
You can identify them by the first line of each test file. Unit tests are tagged
with `+build unit` and integration tests are tagged with `+build integration`.
For example, the test file `deploymentreader_test.go` under `deployers/` contains
unit test cases, so it is indicated with the unit tests tag on the top of the file:

```
// +build unit

package tests
...
```

Integration tests are indicated with the integration tests tag on the top of the file:
```
// +build integration

package tests
...
```

## How to run wskdeploy tests

Before running unit or integration tests, you need to install the package
[Testify](https://github.com/stretchr/testify/):

```
cd $GOPATH
go get -u github.com/stretchr/testify

```

### Running Unit Tests

After installing Testify, all the unit tests can be run from the main `incubator-openwhisk-wskdeploy`
repository folder using the following command:

```
$ go test -v ./... -tags unit
```

### Running Integration Tests

After installing Testify, all integration tests can be run rom the main `incubator-openwhisk-wskdeploy`
repository folder using the following command:

```
$ go test -v ./... -tags integration
```

### Run an Individual Test

```
go test -v <path-to-package-dir> -tags unit -run <test-name>
```

For example:

```
go test -v ./parsers/ -tags unit -run TestComposeActionsForWebActions
go test -v ./tests/src/integration/zipaction/ -tags integration
```

## How to Structure Unit Tests

All integration test cases are put under the folder `tests/src/integration`.

Unit tests are co-located in the same directory as the package being tested (by Go convention).
The test files uses the same basename as the file that contains the package,
but with the added '_test' postfix to the base file name.

For example, the package `deployers`, which defines a `DeploymentReader` service,
is declared in the file `deploymentreader.go`; its corresponding unit test file
should be named `deploymentreader_test.go` under `deployers/`.

Also, the manifest and deployment YAML files used by unit tests should go under `tests/dat`.

## How to Structure Integration Tests

Every integration test has to have a test file 
Integration tests are created under `tests/src/integration`. The test file and
manifest and deployment YAML files should go under the same directory.

For example, `helloworld` integration test:

```
ls -1 tests/src/integration/helloworld/
README.md
actions/
deployment.yaml
helloworld_test.go
manifest.yaml
```
