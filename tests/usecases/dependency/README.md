# Dependent Packages with `wskdeploy` 

`wskdeploy` supports dependencies where it allows you to declare other OpenWhisk
packages that your application or project (manifest) is dependent on. With declaring
dependent packages, `wskdeploy` supports automatic deployment of those dependent
packages.

`wskdeploy` classifies `dependencies` into two different categories:

### GitHub Dependency

Any package with `manifest.yaml` and/or `deployment.yaml` can be treated as
a dependent package and can be specified in `manifest.yaml` with the following
`dependencies` section:

```yaml
dependencies:
    hellowhisk:
        location: github.com/apache/incubator-openwhisk-test/packages/hellowhisk
```

where `hellowhisk` is an external package whose source code is located in
GitHub repo under https://github.com/apache/incubator-openwhisk-test/. When we
deploy our application, `hellowhisk` will be deployed based on the manifest and
deployment files located in the folder `packages/hellowhisk`.

In case of having multiple dependencies on the same package or customizing
package name, we can define `manifest` with:

```yaml
dependencies:
    hellowhisk:
        location: github.com/apache/incubator-openwhisk-test/packages/hellowhisk
    myhelloworlds:
        location: github.com/apache/incubator-openwhisk-test/packages/helloworlds
```

### Package Binding

Any package binding with external services like cloudant, message hub, etc
can be specified in `manifest.yaml` with:

```yaml
dependencies:
    cloudant-package-binding:
        location: /whisk.system/cloudant
        inputs:
            username: $CLOUDANT_USERNAME
            password: $CLOUDANT_PASSWORD
            host: ${CLOUDANT_USERNAME}.cloudant.com
```

Now, we can define a trigger based on the package binding from `dependencies`
section with:

```
triggers:
    cloudant-binding-trigger:
        source: cloudant-package-binding/changes
        inputs:
            dbname: $CLOUDANT_DATABASE
```


### Step 1: Deploy

Deploy it using `wskdeploy`:

```
$ wskdeploy -m ./tests/usecases/dependency/manifest.yaml
```

### Step 2: Verify

```
$ wsk package get --summary hellowhisk
package /<your_name_space>/hellowhisk
 action /<your_name_space>/hellowhisk/greeting
 action /<your_name_space>/hellowhisk/httpGet

$ wsk package get --summary myhelloworlds
package /<your_name_space>/myhelloworlds
 action /<your_name_space>/myhelloworlds/hello-js
 action /<your_name_space>/myhelloworlds/helloworld-js

$ wsk package get --summary cloudant-package-binding 
```

### Step 3: Run

Fire the `hello-trigger` and notice `hello-series` is activated:

```
Activation: hello-trigger (de471ec7b1644d87aa74bb2f725bddfc)
[]

Activation: github-rule (70bdf87d92de4db49be93f6b6f1eeab5)
[]

Activation: hello-series (fba84d20b6104acea110cd558ba7bc0d)
[
    "1604d3fe30fb449699e367fff41d890b"
]
```

Also, update any existing document or create a new document in cloudant and notice
how `cloudant-binding-trigger` is activated:

```
Activation: read (a784f0ab092d4c67b357d5f401be14c4)
[
    "2017-09-26T22:09:59.394555644Z stdout: success { _id: 'aa4dbc5b6c72c236f9813c7d2a25eb2a',",
    "2017-09-26T22:09:59.394597168Z stdout: _rev: '1-967a00dff5e02add41819138abb3284d' }"
]
```

### Step 4: Uninstall

```
$ wskdeploy undeploy -m ./tests/usecases/dependency/manifest.yaml
```
