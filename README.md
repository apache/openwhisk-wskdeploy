# Whisk Deploy `wskdeploy`

`wskdeploy` is a utility to help you create and deploy OpenWhisk projects. Deploy all your actions, triggers, and rules using a single command! You can use this in addition to the OpenWhisk CLI.

`wskdeploy` is currenty under development and in its very early stages.  Check back often to see how its progressing.

# How to use
`wskdeploy` is written in Go. You can invoke it as a Go program, or run its binary file `wskdeploy` directly after building.

Using command
```
$ go run main.go --help
```
or
```
$ ./wskdeploy --help
```
, you can get the detail usage of this tool.


For example,
```
$ go run main.go -m tests/testcases/triggerrule/manifest.yml -d tests/testcases/triggerrule/deployment.yml
```
or
```
$ ./wskdeploy -m tests/testcases/triggerrule/manifest.yml -d tests/testcases/triggerrule/deployment.yml
```
will deploy the `triggerrule` test case.

# How to build on local host
`wskdeploy` can be built with Go tool.

Make sure `$GOPATH` is defined. If not, setup your [Go development environment](https://golang.org/doc/code.html).

Then download `openwhisk-wskdeploy` and dependencies by typing:

```sh
$ cd $GOPATH
$ go get github.com/openwhisk/openwhisk-wskdeploy  # see known issues below if you get an error
```

And finally build `wskdeploy`

```sh
$ cd src/github.com/openwhisk/openwhisk-wskdeploy/
$ go build
```

If you want to build with the godep tool, please execute the following commands.

```
$ go get github.com/tools/godep # Install the godep tool.
$ godep get                     # Download and install packages with specified dependencies.
$ godep go build                # build the wskdeploy tool.
```

Note: we have no releases yet so you should build the `development` branch.

# Contributing

Start by creating a fork of `openwhisk-wskdeploy` and then change the git `origin` to point to
your forked repository, as follows:

```sh
$ cd $GOPATH/src/github.com/openwhisk/openwhisk-wskdeploy
$ git remote rename origin upstream
$ git remote add origin https://github.com/<your fork>/openwhisk-wskdeploy
$ git branch --set-upstream-to origin/master  # track master from origin now
```

You can now use `git push` to push changes to your repository and submit pull requests.

# How to Build with Docker
If you don't want to bother with go installation, build, git clone etc, you can do it with Docker, then
you can run wskdeploy tool in your container.

1. First you need a docker daemon running locally on your machine or whatever in a VM etc.

2. Get the Docker file.
 ```
 wget -O Dockerfile https://raw.githubusercontent.com/openwhisk/wskdeploy/master/Dockerfile
 ```

3. Build and tag a docker image.
```
docker build -f Dockerfile .  -t openwhisk/wskdeploy
```

4. Bring up the docker container.
```
docker run -ti openwhisk/wskdeploy
```
5. Inside the container, run `wskdeploy` and have fun.

Note: Based on user role, you may need add sudo before your command to run as root.

# Known issues

You might get this error when downloading `openwhisk-wskdeploy`

     Cloning into ''$GOAPTH/src/gopkg.in/yaml.v2'...
     error: RPC failed; HTTP 301 curl 22 The requested URL returned error: 301
     fatal: The remote end hung up unexpectedly

This is caused by newer `git` not forwarding request anymore. One solution is to allow forwarding for `gopkg.in`

```sh
$ git config --global http.https://gopkg.in.followRedirects true
```
