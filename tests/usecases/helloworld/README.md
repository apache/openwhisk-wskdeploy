# "helloworld" test case for wskdeploy.

### Package description

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

### How to deploy and test

Step 1. Deploy the package.

Deploy it using the  ```wskdeploy``` command as follows using the manifest.yaml file found in the "path" indicated on the '-p' command line option:
```
$ wskdeploy -p tests/usecases/helloworld
```

Step 2. Verify the actions were created.
```
$ wsk action list
```

Step 3. Invoke the action.
```
$ wsk action invoke --blocking --result helloworld/hello --param name Bernie --param place Vermont
```
