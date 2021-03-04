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

`wskdeploy` is a utility to help you describe and deploy any part of the OpenWhisk programming model using a YAML manifest file. Use it to deploy all of your OpenWhisk project's [Packages](https://github.com/apache/openwhisk/blob/master/docs/packages.md), [Actions](https://github.com/apache/openwhisk/blob/master/docs/actions.md), [Triggers, and Rules](https://github.com/apache/openwhisk/blob/master/docs/triggers_rules.md), together, using a single command!

#### Running `wskdeploy` standalone

You can use this utility separately from the OpenWhisk CLI as it uses the same [OpenWhisk "Go" Client](https://github.com/apache/openwhisk-client-go) as the Openwhisk CLI does to create its HTTP REST calls for deploying and undeploying your Openwhisk packages and entities.

#### Running `wskdeploy` as part of the `wsk` CLI

Alternatively, you can use the `wskdeploy` functionality within the OpenWhisk CLI as it is now embedded as the `deploy` command. That is, you can invoke it as `wsk deploy` using all the same parameters documented for the standalone utility.

#### Using `wskdeploy` to manage OpenWhisk entities as projects

In addition to simple deployment, `wskdeploy` also has the powerful `export` command to manage sets of OpenWhisk entities that work together as a named project. The command:

```sh
wskdeploy export --projectname <managed_project_name>`
```

allows you to "export" a specified project into a local file system and manage it as a single entity.

In the above example, a `<managed_project_name>.yml` Manifest file would be created automatically which can be used with `wskdeploy` to redeploy the managed project on a different OpenWhisk instance. If the managed project contains dependencies on other managed projects, then these projects will be exported automatically into their respective manifests.

## Getting started

Here are some quick links to help you get started:

- [Downloading released binaries](#downloading-released-binaries) - released binaries for Linux, Mac OS and Windows
- [Running wskdeploy](#running-wskdeploy) - run `wskdeploy` as a binary or Go program
- :eight_spoked_asterisk: [Writing Package Manifests](docs/programming_guide.md#wskdeploy-utility-by-example) - a step-by-step guide on writing Package Manifest files for ```wskdeploy```
- :eight_spoked_asterisk: [Exporting OpenWhisk assets](docs/export.md) - how to use `export` feature
- [Building the project](#building-the-project) - download and build the GoLang source code
- [Contributing to the project](#contributing-to-the-project) - join us!
- [Debugging wskdeploy](docs/wskdeploy_debugging.md) - helpful tips for debugging the code and your manifest files
- [Troubleshooting](#troubleshooting) - known issues (e.g., Git)

---

## Downloading released binaries

Executable binaries of `wskdeploy` are available for download on the project's GitHub [releases](https://github.com/apache/openwhisk-wskdeploy/releases) page:

- [https://github.com/apache/openwhisk-wskdeploy/releases](https://github.com/apache/openwhisk-wskdeploy/releases).

We currently provide binaries for the following Operating Systems (OS) and architecture combinations:

Operating System | Architectures
--- | ---
Linux | 386, AMD64, ARM, ARM64, PPC64 (Power), S/390 and IBM Z
Mac OS (Darwin) | 386<sup>[1](#1)</sup>, AMD64
Windows | 386, AMD64

1. Mac OS, 32-bit (386) released versions are not available for builds using Go lang version 1.15 and greater.

We also provide instructions on how to build your own binaries from source code. See [Building the project](#building-the-project).

---

## Running ```wskdeploy```

Start by verifying the utility can display the command line help:
```sh
$ ./wskdeploy --help
```

then try deploying an OpenWhisk Manifest and Deployment file:
```sh
$ ./wskdeploy -m tests/usecases/triggerrule/manifest.yml -d tests/usecases/triggerrule/deployment.yml
```

---

## Building the project

### GoLang setup

The wskdeploy utility is a GoLang program so you will first need to [Download and install GoLang](https://golang.org/dl/) onto your local machine.

> **Note** Go version 1.15 or higher is recommended

Make sure your `$GOPATH` is defined correctly in your environment. For detailed setup of your GoLang development environment, please read [How to Write Go Code](https://golang.org/doc/code.html).

### Download the source code from GitHub

As the code is managed using GitHub, it is easiest to retrieve the code using the `git clone` command.

if you just want to build the code and do not intend to be a Contributor, you can clone the latest code from the Apache repository:

```sh
git clone git@github.com:apache/openwhisk-wskdeploy
```

or you can specify a release (tag) if you do not want the latest code by using the `--branch <tag>` flag. For example, you can clone the source code for the tagged 1.1.0 [release](https://github.com/apache/openwhisk-wskdeploy/releases/tag/1.1.0)

```sh
git clone --branch 1.1.0 git@github.com:apache/openwhisk-wskdeploy
```

You can also pull the code from a fork of the repository. If you intend to become a Contributor to the project, read the section [Contributing to the project](#contributing-to-the-project) below on how to setup a fork.

### Build using `go build`

Use the Go utility to build the ```wskdeploy``` binary.

Change into the cloned project directory and use `go build` with the target output name for the binary:

```sh
$ go build -o wskdeploy
```

If successful, an executable named `wskdeploy` will be created in the project directory compatible with your current operating system and architecture.

#### Building for other Operating Systems (GOOS) and Architectures (GOARCH)

If you would like to build the binary for a specific operating system and processor architecture, you may add the arguments `GOOS` and `GOARCH` into the Go build command (as inline environment variables).

For example, run the following command to build the binary for 64-bit Linux:

```sh
$ GOOS=linux GOARCH=amd64 go build -o wskdeploy
```

Supported value combinations include:

`GOOS` | `GOARCH`
--- | ---
linux | 386 (32-bit), amd64 (64-bit), s390x (S/390, Z), ppc64le (Power), arm (32-bit), arm64 (64-bit)
darwin (Mac OS) | amd64
windows | 386 (32-bit), amd64 (64-bit)

### Build using Gradle

The project includes its own packaged version of Gradle called Gradle Wrapper which is invoked using the `gradlew` command on Linux/Unix/Mac or `gradlew.bat` on Windows.

1. Gradle requires requires you to [install Java JDK version 8](https://gradle.org/install/) or higher

1. Clone the `openwhisk-wskdeploy` repo:

    ```sh
    git clone https://github.com/apache/openwhisk-wskdeploy
    ```

    and change into the project directory.

1. Cross-compile binaries for all supported Operating Systems and Architectures:

    ```sh
    ./gradlew goBuild
    ```

1. Upon a successful build, the `wskdeploy` binaries can be found under the corresponding `build/<os>-<architecture>/` folder of your project:

    ```sh
    $ ls build
    darwin-amd64  linux-amd64   linux-arm64   linux-s390x   windows-amd64
    linux-386     linux-arm     linux-ppc64le windows-386
    ```

#### Compiling for a single OS/ARCH

1. View gradle build tasks for supported Operating Systems and Architectures:

    ```sh
    ./gradlew tasks
    ```

    you will see build tasks for supported OS/ARCH combinations:

    ```sh
    Gogradle tasks
    --------------
    buildDarwinAmd64 - Custom go task.
    buildLinux386 - Custom go task.
    buildLinuxAmd64 - Custom go task.
    buildLinuxArm - Custom go task.
    buildLinuxArm64 - Custom go task.
    buildLinuxPpc64le - Custom go task.
    buildLinuxS390x - Custom go task.
    buildWindows386 - Custom go task.
    buildWindowsAmd64 - Custom go task.
    ```

    > **Note**: The `buildWindows386` option is only supported on Golang versions less than 1.15.

1. Build using one of these tasks, for example:

    ```sh
    $ ./gradlew buildDarwinAmd64
    ```

#### Using your own local Gradle to build

Alternatively, you can choose to [Install Gradle](https://gradle.org/install/) and use it instead of the project's Gradle Wrapper.  If so, you would use the `gradle` command instead of `gradlew`. If you do elect to use your own Gradle, verify its version is `6.8.1` or higher:

```sh
gradle -version
```

> **Note** If using your own local Gradle installation, use the `gradle` command instead of the `./gradlew` command in the build instructions below.

#### Building for internationalization

Please follow this process for building any changes to translatable strings:
- [How to generate the file i18n_resources.go for internationalization](https://github.com/apache/openwhisk-wskdeploy/blob/master/wski18n/README.md)

---

### Running as a Go program

Since ```wskdeploy``` is a GoLang program, you may choose to run it using the Go utility. After building the wskdeploy binary, you can run it as follows:

```sh
$ go run main.go --help
```

and deploying using the Go utility would look like:
```sh
$ go run main.go -m tests/usecases/triggerrule/manifest.yml -d tests/usecases/triggerrule/deployment.yml
```

---

## Contributing to the project

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

> Be sure to [Sync your fork](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/syncing-a-fork) before starting any contributions to keep it up-to-date with the upstream repository.

### Running unit tests

You may use `go test` to test all unit tests within a package, for example:

```sh
go test ./deployers -tags=unit -v
go test ./parsers -tags=unit -v
```

or to run individual function tests, for example:

```sh
go test ./parsers -tags=unit -v -run TestParseManifestForSingleLineParams
```

### Running integration tests

Integration tests are best left to the Travis CI build as they depend on a fully functional OpenWhisk environment to be deployed.

### Adding new dependencies

Please use `go get` to add new dependencies to the `go.mod` file:

```sh
go get github.com/project/libname@v1.2.0
```

> Please avoid using commit hashes for referencing non-OpenWhisk libraries.

### Removing unused dependencies

Please us `go tidy` to remove any unused dependencies after any significant code changes:

```sh
go mod tidy
```

### Updating dependency versions

Although you might be tempted to edit the go.mod file directly, please use the recommended method of using the `go get` command:

```sh
go get -u github.com/project/libname  # Using "latest" version
go get -u github.com/project/libname@v1.1.0 # Using tagged version
go get -u github.com/project/libname@aee5cab1c  # Using a commit hash
```

### Updating Go version

Although you could edit the version directly in the go.mod file, it is better to use the `go edit` command:

```sh
go mod edit -go=1.15
```

### Creating Tagged Releases

Committers can find instructions on how to create tagged releases here:
- [creating_tagged_releases.md](https://github.com/apache/openwhisk-wskdeploy/tree/master/docs/creating_tagged_releases.md)

---

## Troubleshooting

### Known issues

#### Git commands using HTTPS, not SSH

The "go get" command uses HTTPS with GitHub and when you attempt to "commit" code you might be prompted with your GitHub credentials.  If you wish to use your SSH credentials, you may need to issue the following command to set the appropriate URL for your "origin" fork:

```sh
git remote set-url origin git@github.com:<username>/openwhisk-wskdeploy.git
```

or you can manually change the remote (origin) url within your .git/config file:

```sh
[remote "origin"]
    url = git@github.com:<username>/openwhisk-wskdeploy
```

while there, you can verify that your upstream repository is set correctly:

```sh
[remote "upstream"]
    url = git@github.com:apache/openwhisk-wskdeploy
```

#### Git clone RPC failed: HTTP 301

This sometimes occurs using "go get" the wskdeploy code (which indirectly invokes "git clone").

You might get this error when downloading `openwhisk-wskdeploy`:

```sh
Cloning into ''$GOAPTH/src/gopkg.in/yaml.v2'...
error: RPC failed; HTTP 301 curl 22 The requested URL returned error: 301
fatal: The remote end hung up unexpectedly
```

This is caused by newer `git` versions not forwarding requests anymore. One solution is to allow forwarding for `gopkg.in`

```sh
$ git config --global http.https://gopkg.in.followRedirects true
```
