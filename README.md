# Whisk Deploy `wskdeploy`


DISCLAIMER - wskdeploy is an experimental tool.
-----------------------------------------------


`wskdeploy` is a utility to help you create and deploy OpenWhisk projects. Deploy all your actions, triggers, and rules using a single command! You can use this in addition to the OpenWhisk CLI.

`wskdeploy` is currenty under development and in its very early stages.  Check back often to see how its progressing.

# How to use
`wskdeploy` is written in Go. You can invoke it as a Go program, or run its binary file `wskdeploy` directly after building.

You can get the detail usage of this tool by using following commands:

```
$ go run main.go --help
```

or

```
$ ./wskdeploy --help
```

For example,

```
$ go run main.go -m tests/usecases/triggerrule/manifest.yml -d tests/usecases/triggerrule/deployment.yml
```

or

```
$ ./wskdeploy -m tests/usecases/triggerrule/manifest.yml -d tests/usecases/triggerrule/deployment.yml
```

will deploy the `triggerrule` test case.

# Where to download the binary wskdeploy
`wskdeploy` is available on the release page of openwhisk-wskdeploy project: [click here to download](https://github.com/apache/incubator-openwhisk-wskdeploy/releases).
We currently have binaries available for Linux, Mac OS and windows under amd64 architecture. You can find the binary, which fits your local environment.


# How to build on local host
There is another approach to get the binary `wskdeploy`, which is to build it from the source code with Go tool.

Make sure `$GOPATH` is defined. If not, setup your [Go development environment](https://golang.org/doc/code.html).

Then download `wskdeploy` and dependencies by typing:

```sh
$ cd $GOPATH
$ go get github.com/apache/incubator-openwhisk-wskdeploy  # see known issues below if you get an error
```

And finally build `wskdeploy`

```sh
$ cd src/github.com/apache/incubator-openwhisk-wskdeploy/
$ go build -o wskdeploy
```

If you want to build with the godep tool, please execute the following commands.

```sh
$ go get github.com/tools/godep # Install the godep tool.
$ godep get                     # Download and install packages with specified dependencies.
$ godep go build -o wskdeploy   # build the wskdeploy tool.
```

You can verify your build by running:

```sh
./wskdeploy --help
```

Note: we have no releases yet so you should build the `development` branch.

# Contributing

Start by creating a fork of `openwhisk-wskdeploy` and then change the git `origin` to point to
your forked repository, as follows:

```sh
$ cd $GOPATH/src/github.com/apache/incubator-openwhisk-wskdeploy
$ git remote rename origin upstream
$ git remote add origin https://github.com/<your fork>/incubator-openwhisk-wskdeploy
$ git branch --set-upstream-to origin/master  # track master from origin now
```

You can now use `git push` to push changes to your repository and submit pull requests.

# How to Cross Compile Binary with Gradle/Docker
If you don't want to bother with go installation, build, git clone etc, and you can do it with Gradle/Docker.
After compiling, a suitable wskdeploy binary that works for your OS platform will be available under /bin directory.

1. First you need a docker daemon running locally on your machine.

2. Make sure you have Java 1.7 or above installed.

3. Clone the wskdeploy repo with command ```git clone https://github.com/apache/incubator-openwhisk-wskdeploy.git```

4. If you use Windows OS, type ```gradlew.bat -version ```. For Unix/Linux/Mac, please type ```./gradlew -version```.

5. Make sure you can see the correct Gradle version info on your console. Currently the expected Gradle
version is 3.3.

6. For Windows type ```gradlew.bat distDocker```. For Linux/Unix/Mac, please type ```./gradlew distDocker```. These
commands will start the wskdeploy cross compile for your specific OS platform inside a Docker container.

7. After build success, you should find a correct binary under current /bin dir of you openwhisk-deploy clone dir.


# Known issues

You might get this error when downloading `openwhisk-wskdeploy`

     Cloning into ''$GOAPTH/src/gopkg.in/yaml.v2'...
     error: RPC failed; HTTP 301 curl 22 The requested URL returned error: 301
     fatal: The remote end hung up unexpectedly

This is caused by newer `git` not forwarding request anymore. One solution is to allow forwarding for `gopkg.in`

```sh
$ git config --global http.https://gopkg.in.followRedirects true
```

DISCLAIMER - wskdeploy is an experimental tool.
-----------------------------------------------
