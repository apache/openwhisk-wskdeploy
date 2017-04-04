# Test cases for openwhisk-wskdeploy.

There are two types of test cases located under the folder of tests: unit and integration. You can identify
them by the first line of each test case file. Unit tests are tagged with "+build unit" and integration tests
are tagged with "+build integration". For example, the test file deploymentreader_test.go under
tests/src/deployers contains unit test cases, so it is indicated with the unit tag on the top of the file:
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

# How to run the tests
You need to install the package Testify before running the test cases:
```
cd $GOPATH
go get -u github.com/stretchr/testify

```

Then you are able to run all the unit tests via the following command under the folder of openwhisk-wkdeploy:
```
$ go test -v ./... -tags unit
```

In order to run all the integration tests, run the following command under the folder of openwhisk-wkdeploy:
```
$ go test -v ./... -tags integration
```

# How to structure the test cases
All the test cases are put under the folder tests/src, and all the integration test cases are put under
the folder tests/src/integration.

Different unit tests can be grouped into different packages, and they can be put under a subfolder
named after the package under tests/src. The yaml files of manifest and deployment used by unit tests are
put under tests/data.

Different integration tests can be grouped into different packages as well, and they can be put under a
subfolder named after the package under tests/src/integration. The yaml files of manifest and deployment
used by the integration test are put under the same subfolder as the integration test itself. The source
file used by the manifest file can be put under the folder tests/src/integration/<package>/src.