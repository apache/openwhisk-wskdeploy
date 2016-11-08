# Whisk Deploy `wskdeploy`

`wskdeploy` is a utility to help you create and deploy OpenWhisk projects. Deploy all your actions, triggers, and rules using a single command! You can use this in addition to the OpenWhisk CLI.

`wskdeploy` is currenty under development and in its very early stages.  Check back often to see how its progressing.

# How to Build
`wskdeploy` is written in Go.

1. Setup your [Go development environment](https://golang.org/doc/code.html).

2.  `wskdeploy` depends on the `github.com/openwhisk/openwhisk-client-go/whisk` . To install:

``` go get github.com/openwhisk/openwhisk-client-go/whisk ```

3. Clone this repo into `$GOPATH/src/github.com/openwhisk`, which should have been created by Step #2.

```
$ cd $GOPATH/src/github.com/openwhisk
$ git clone http://github.com/openwhisk/wskdeploy
```

4. Tagged releases are in master. The latest build is always in the development branch. Inside `$GOPATH/src/github.com/openwhisk/wskdeply`:

```
$ git checkout development   ## or skip this step and just build master
$ go build
```

# whiskdeploy
Utility for managing OpenWhisk Projects
>>>>>>> Update README.md
