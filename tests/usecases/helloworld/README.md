# Test case for Whisk Deploy

This is "helloworld" test case for wskdeploy.

The [manifest.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/tests/usecases/helloworld/manifest.yaml) file defines:

- a Package named "helloworld" which contains two Actions:
    - an Action named "hello"
    - an Action named "helloworld"
    - Two optional triggers "trigger1" and "trigger2"

Both the "hello" and "helloworld" Actions accept the same two parameters:
    - "name" (string)
    - "place" (string)

The "hello" Action has the following defaults for these properties:
    - name: Paul
    - place Boston

The "helloworld" Action has the following defaults for these properties:
    - name: Bernie
    - place Vermont

and will return a greeting message:

```
"Hello, <name> from <place>!"
```

Deploy it using the  ```wskdeploy``` command as follows using the manifest.yaml file found in the "path" indicated on the '-p' command line option:
```
$ wskdeploy -p tests/usecases/helloworld
```
and can be verified and invoked as follows:
```
$ wsk action list
$ wsk action invoke --blocking --result helloworld/hello --param name Bernie --param place Vermont
```
