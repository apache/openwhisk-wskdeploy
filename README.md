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

# Whisk Deploy `wskdeploy`

[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![Build Status](https://travis-ci.com/apache/openwhisk-wskdeploy.svg?branch=master)](https://travis-ci.com/apache/openwhisk-wskdeploy)

`wskdeploy` is a utility to help you describe and deploy any part of the OpenWhisk programming model using a Manifest file written in YAML. Use it to deploy all your OpenWhisk [Packages](https://github.com/apache/openwhisk/blob/master/docs/packages.md), [Actions](https://github.com/apache/openwhisk/blob/master/docs/actions.md), [Triggers, and Rules](https://github.com/apache/openwhisk/blob/master/docs/triggers_rules.md) using a single command!

`wskdeploy export --projectname managed_project_name` allows to "export" a specified managed project into a local file system. Namely, a `managed_project_name.yml` Manifest file will be created automatically. This Manifest file can be used with `wskdeploy` to redeploy the managed project at a different OpenWhisk instance. If the managed project contains dependencies on other managed projects, then these projects will be exported automatically into their respective manifests.

You can use this in addition to the OpenWhisk CLI.  In fact, this utility uses the [OpenWhisk "Go" Client](https://github.com/apache/openwhisk-client-go) to create its HTTP REST calls for deploying and undeploying your packages.

## Here are some quick links for:

- [Downloading wskdeploy](#downloading-released-binaries) - released binaries for Linux, Mac OS and Windows
- [Running wskdeploy](#running-wskdeploy) - run wskdeploy as a binary or Go program
- :eight_spoked_asterisk: [Writing Package Manifests](docs/programming_guide.md#wskdeploy-utility-by-example) - a step-by-step guide on writing Package Manifest files for ```wskdeploy```
- :eight_spoked_asterisk: [Exporting OpenWhisk assets](docs/export.md) - how to use `export` feature
- [Building the project](#building-the-project) - download and build the GoLang source code
- [Contributing to the project](#contributing-to-the-project) - join us!
- [Debugging wskdeploy](docs/wskdeploy_debugging.md) - helpful tips for debugging the code and your manifest files
- [Troubleshooting](#troubleshooting) - known issues (e.g., Git)

---

## Contributing to the project

### GoLang setup

The wskdeploy utility is a GoLang program so you will first need to [Download and install GoLang](https://golang.org/dl/) onto your local machine.

> **Note** Go version 1.15 or higher is recommended

Make sure your `$GOPATH` is defined correctly in your environment. For detailed setup of your GoLang development environment, please read [How to Write Go Code](https://golang.org/doc/code.html).

### Git repository setup

1. [Fork](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo) the Apache repository

    If you intend to contribute code, you will want to fork the `apache/openwhisk-wskdeploy` repository into your github account and use that as the source for your clone.

1. Clone the repository from your fork:

    ```sh
    git clone git@github.com:${GITHUB_ACCOUNT_USERNAME}/openwhisk-wskdeploy.git
    ```

1. Add the Apache repository as a remote with the `upstream` alias:

    ```sh
    git remote add upstream git@github.com:apache/openwhisk-wskdeploy
    ```

    You can now use `git push` to push local `commit` changes to your `origin` repository and submit pull requests to the `upstream` project repository.

 1. Optionally, prevent accidental pushes to `upstream` using this command:

    ```sh
    git remote set-url --push upstream no_push
    ```

### Sync your fork before starting commits

Be sure to [Sync your fork](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/syncing-a-fork) before starting any contributions to keep it up-to-date with the upstream repository.

---

## Building the `wskdeploy` binary

Use the Go utility to build the ```wskdeploy``` binary

Change into the cloned project directory and use `go build` with the target output name for the binary:

```sh
$ go build -o wskdeploy
```

### Building for other Operating Systems (GOOS) and Architectures (GOARCH)

If you would like to build the binary for a specific operating system, you may add the arguments GOOS and GOARCH into the Go build command. You may set

- ```GOOS``` to "linux", "darwin" or "windows"
- ```GOARCH``` to "amd64" or "386"

For example, run the following command to build the binary for 64-bit Linux:

```sh
$ GOOS=linux GOARCH=amd64 go build -o wskdeploy
```

### How to Cross Compile Binary with Gradle/Docker

If you don't want to bother with go installation, build, git clone etc, and you can do it with Gradle/Docker.

After compiling, a suitable wskdeploy binary that works for your OS platform will be available under /bin directory.

1. First you need a docker daemon running locally on your machine.

2. Make sure you have Java 8 or above installed.

3. Clone the wskdeploy repo with command ```git clone https://github.com/apache/openwhisk-wskdeploy.git```

4. If you use Windows OS, type ```gradlew.bat -version ```. For Unix/Linux/Mac, please type ```./gradlew -version```.

5. Make sure you can see the correct Gradle version info on your console. Currently the expected Gradle
version is 3.3.

6. For Windows type ```gradlew.bat distDocker```. For Linux/Unix/Mac, please type ```./gradlew distDocker```. These
commands will start the wskdeploy cross compile for your specific OS platform inside a Docker container.

7. After build success, you should find a correct binary under current /bin dir of you openwhisk-deploy clone dir.

If you would like to build the binaries available for all the operating systems and architectures, run the following command:

```sh
$ ./gradlew distDocker -PcrossCompileWSKDEPLOY=true
```

Then, you will find the binaries and their compressed packages generated under the folder ```bin/<os>/<arch>/``` for each supported Operating System and CPU Architecture pair.

### Building for Internationalization

Please follow this process for building any changes to translatable strings:
- [How to generate the file i18n_resources.go for internationalization](https://github.com/apache/openwhisk-wskdeploy/blob/master/wski18n/README.md)

---

## Running ```wskdeploy```

After building the wskdeploy binary, you can run it as follows:

### Running the Binary file

Start by verifying the utility can display the command line help:
```sh
$ ./wskdeploy --help
```

then try deploying an OpenWhisk Manifest and Deployment file:
```sh
$ ./wskdeploy -m tests/usecases/triggerrule/manifest.yml -d tests/usecases/triggerrule/deployment.yml
```

### Running as a Go program

Since ```wskdeploy``` is a GoLang program, you may choose to run it using the Go utility:
```sh
$ go run main.go --help
```

and deploying using the Go utility would look like:
```sh
$ go run main.go -m tests/usecases/triggerrule/manifest.yml -d tests/usecases/triggerrule/deployment.yml
```

---

## Downloading released binaries

Binaries of `wskdeploy` are available for download on the project's GitHub release page:
- [https://github.com/apache/openwhisk-wskdeploy/releases](https://github.com/apache/openwhisk-wskdeploy/releases).

For each release, we typically provide binaries built for Linux, Mac OS (Darwin) and Windows on the AMD64 architecture. However, we provide instructions on how to build your own binaries as well from source code with the Go tool.  See [Building the project](#building-the-project).

_If you are a Developer or Contributor, **we recommend building from the latest source code** from the project's master branch._

for end-users, please use versioned releases of binaries.
- [https://github.com/apache/openwhisk-wskdeploy/releases](https://github.com/apache/openwhisk-wskdeploy/releases)

<!-- ----------------------------------------------------------------------------- -->

## Troubleshooting

### Known issues

#### Git commands using HTTPS, not SSH

The "go get" command uses HTTPS with GitHub and when you attempt to "commit" code you might be prompted with your GitHub credentials.  If you wish to use your SSH credentials, you may need to issue the following command to set the appropriate URL for your "origin" fork:

```
git remote set-url origin git@github.com:<username>/openwhisk-wskdeploy.git
```

<or> you can manually change the remote (origin) url within your .git/config file:
```
[remote "origin"]
    url = git@github.com:<username>/openwhisk-wskdeploy
```

while there, you can verify that your upstream repository is set correctly:
```
[remote "upstream"]
    url = git@github.com:apache/openwhisk-wskdeploy
```

#### Git clone RPC failed: HTTP 301

This sometimes occurs using "go get" the wskdeploy code (which indirectly invokes "git clone").

<b>Note: Using "go get" for development is unsupported; instead, please use "go deps" for dependency management.</b>

You might get this error when downloading `openwhisk-wskdeploy`:

     Cloning into ''$GOAPTH/src/gopkg.in/yaml.v2'...
     error: RPC failed; HTTP 301 curl 22 The requested URL returned error: 301
     fatal: The remote end hung up unexpectedly

This is caused by newer `git` versions not forwarding requests anymore. One solution is to allow forwarding for `gopkg.in`

```
$ git config --global http.https://gopkg.in.followRedirects true
```

## Creating Tagged Releases

Committers can find instructions on how to create tagged releases here:
- [creating_tagged_releases.md](https://github.com/apache/openwhisk-wskdeploy/tree/master/docs/creating_tagged_releases.md)
