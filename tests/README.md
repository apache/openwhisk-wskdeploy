<!--
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
-->

# `wskdeploy` Test Cases and Real World Applications

## Test Cases

There are two types of test cases supported (1) unit test and (2) integration test.
You can identify them by the first line of each test file `_test.go`.

### Unit Tests

Unit tests are tagged with `+build unit` tag. For example, the test file
[deploymentreader_test.go](https://github.com/apache/openwhisk-wskdeploy/blob/master/deployers/deploymentreader_test.go)
under `deployers/` contains unit test cases which is indicated with the unit tests
tag on the top of the file:

```
// +build unit

package tests
...
```

#### How do I run unit tests?

In order to run any unit tests, you need to install the package [Testify](https://github.com/stretchr/testify/).
After installing Testify, all the unit tests can be run from the main `openwhisk-wskdeploy`
repository folder using the following command:


```
cd $GOPATH
go get -u github.com/stretchr/testify
cd openwhisk-wskdeploy/
$ go test -v ./... -tags unit
```

#### How do I run an individual test?

Above command will run all the unit tests from `openwhisk-wskdeploy`, in
order to run a specific test, use:

```
go test -v <path-to-package-dir> -tags unit -run <test-name>
```

For example:

```
go test -v ./parsers/ -tags unit -run TestComposeActionsForWebActions
```

#### Where should I write unit tests?

Unit tests are co-located in the same directory as the package being tested (by Go convention).
The test files uses the same basename as the file that contains the package,
but with the added `_test` postfix to the base file name.

For example, the package `deployers`, which defines a `DeploymentReader` service,
is declared in the file `deploymentreader.go`; its corresponding unit test file
should be named `deploymentreader_test.go` under `deployers/`.

Also, the manifest and deployment YAML files used by unit tests should go under `tests/dat`.


### Integration Tests

Integration tests are tagged with `+build integration` tag. For example, the test
file [zipaction_test.go](https://github.com/apache/openwhisk-wskdeploy/blob/master/tests/src/integration/zipaction/zipaction_test.go)
contains integration test which is indicated with the integration tests tag on
the top of the file:

```
// +build integration

package tests
...
```

#### How do I run integration tests?

In order to run any integration tests, you need to install the package [Testify](https://github.com/stretchr/testify/).
After installing Testify, all the integration tests can be run from the main `openwhisk-wskdeploy`
repository folder using the following command:


```
cd $GOPATH
go get -u github.com/stretchr/testify
cd openwhisk-wskdeploy/
$ go test -v ./... -tags integration
```

`wskdeploy` tests are located under [tests/](https://github.com/apache/openwhisk-wskdeploy/tree/master/tests)
folder.

#### How do I run an individual test?

Above command will run all the integration tests from `openwhisk-wskdeploy`, in
order to run a specific test, use:

```
go test -v <path-to-package-dir> -tags integration -run <test-name>
```

For example:

```
go test -v ./tests/src/integration/zipaction/ -tags integration -run TestZipAction
```

#### Where should I write integration tests?

All integration test cases are created under the folder `tests/src/integration`.
Every integration test has to have a test file `_test.go` along with manifest and/or
deployment YAML file under the same directory.

For example, `helloworld` integration test:

```
ls -1 tests/src/integration/helloworld/
README.md
actions/
deployment.yaml
helloworld_test.go
manifest.yaml
```

## Real World Applications

[apps](https://github.com/apache/openwhisk-wskdeploy/tree/master/tests/apps)
holds various real world applications which are being deployed using `wskdeploy`.
This space gives an opportunity to `wskdeploy` consumers to integrate with `wskdeploy`
and verify deployment/undeployment of their applications against a clean OpenWhisk
instance. With this shared platform, application developer can work with `wskdeploy`
developers to implement their requirements and usecases as they come in.

As an application developer, you can follow [Contributing to Project](https://github.com/apache/openwhisk-wskdeploy#contributing-to-the-project)\
guide to add your application under [apps](https://github.com/apache/openwhisk-wskdeploy/tree/master/tests/apps)
or if you want to skip cloning the whole [openwhisk-wskdeploy](https://github.com/apache/openwhisk-wskdeploy)
GitHub repo. There is a one time settings possible if you just want to
clone your own application and submit pull requests:

```
# create a new directory where you want to clone your application
$ mkdir <my-wskdeploy-application>
$ cd <my-wskdeploy-application>
# initialize empty local repo
$ git init
# add the remote named origin using your fork
$ git remote add origin -f https://github.com/<application-developer>/openwhisk-wskdeploy.git
# the following git command is very important where we tell git we are checking out specifics
$ git config core.sparsecheckout true
$ echo "apps/*" >> .git/info/sparse-checkout
$ git pull origin master
```
