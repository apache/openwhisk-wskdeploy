# Show case for dependencies

This is a show case for `wskdeploy`. This package shows how to deploy a package with dependencies to existing system packages or external resources.

Below is the `dependencies` segement in `manifest.yaml`.
```
dependencies:
  hellowhisk:
    location: github.com/apache/incubator-openwhisk-test/packages/hellowhisk
  myhelloworlds:
    location: github.com/apache/incubator-openwhisk-test/packages/helloworlds
  myCloudant:
    location: /whisk.system/cloudant
    inputs:
      dbname: myGreatDB
    annotations:
      myAnnotation: Here it is
```

It defines three dependencies:
- `hellowhisk` is a dependency to an external package whose source code is in a folder in a git repo `https://github.com/apache/incubator-openwhisk-test`. `wskdeploy` will deploy the package according to the manifest and deployment files in the folder `packages/hellowhisk` of this repo.
- `myhelloworlds` is a dependency to an external package whose source code is another folder in the git repo `https://github.com/apache/incubator-openwhisk-test`. `wskdeploy` will deploy the package according to the manifest and deployment files in the folder `packages/helloworlds` in this repo. When `wskdeploy` notices the dependency name `myhelloworlds` is different from the original package name `helloworlds`, it will create a package binding `myhelloworlds` to `helloworlds` so that all references to the assets in package `helloworlds` can be referenced by the new package name `myhelloworlds`, for example `myhelloworlds/hello-js` in the sequence definition:
```
sequences:
  mySequence:
    actions: hellowhisk/greeting, hellowhisk/httpGet, myhelloworlds/hello-js
```
- `myCloudant` is a binding to the existing system package `/whisk.system/cloudant`.  `wskdeploy` will create a package binding `myCloudant` to `/whisk.system/cloudant` with the input parameters and annotations.

To install this show case, try this:
```
$ wskdeploy -m ./tests/usecases/dependency/manifest.yaml
```

To verify all the assets are deployed in OpenWhisk, try this:
```
$ wsk list
Entities in namespace: <your_name_space>
packages
/<your_name_space>/hellowhisk                                    private
/<your_name_space>/myCloudant                                    private
/<your_name_space>/myhelloworlds                                 private
/<your_name_space>/helloworlds                                   private
/<your_name_space>/opentest                                      private
actions
/<your_name_space>/opentest/mySequence                           private sequence
/<your_name_space>/hellowhisk/greeting                           private nodejs:6
/<your_name_space>/hellowhisk/httpGet                            private swift:3
/<your_name_space>/helloworlds/hello-js                       	 private nodejs:6
/<your_name_space>/helloworlds/helloworld-js                     private nodejs:6
triggers
/<your_name_space>/myTrigger                                     private
rules
/<your_name_space>/myCloudantRule                                private
/<your_name_space>/myRule                                        private
```
To verify package `hellowhisk` is deployed, try this:
```
$ wsk package get --summary hellowhisk
package /<your_name_space>/hellowhisk
 action /<your_name_space>/hellowhisk/greeting
 action /<your_name_space>/hellowhisk/httpGet
```
To verify package `myhelloworlds` is deployed, try this:
```
$ wsk package get --summary myhelloworlds
package /<your_name_space>/myhelloworlds
 action /<your_name_space>/myhelloworlds/hello-js
 action /<your_name_space>/myhelloworlds/helloworld-js
```
To verify a package binding `mycloudant` is installed, try this:
```
$ wsk package get --summary myCloudant
```
To uninstall this show case, try this:
```
$ wskdeploy undeploy -m ./tests/usecases/dependency/manifest.yaml
```