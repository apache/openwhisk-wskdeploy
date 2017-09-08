# Show case for dependencies

This is a show case for `wskdeploy`. This package shows how to deploy a package with dependencies to existing system packages or external resources.

Below is the `dependencies` segement in `manifest.yaml`.
```
dependencies:
  hellowhisk:
    location: github.com/paulcastro/hellowhisk
  newpkg:
    location: github.com/daisy-ycguo/wskdeploy-test/helloworld
  myCloudant:
    location: /whisk.system/cloudant
    inputs:
      dbname: myGreatDB
    annotations:
      myAnnotation: Here it is
```

It defines three dependencies:
- `hellowhisk` is a dependency to an external package whose source code is in a git repo `https://github.com/paulcastro/hellowhisk`. `wskdeploy` will deploy the package according to the manifest and deployment files in the root folder of `https://github.com/paulcastro/hellowhisk`.
- `newpkg` is a dependency to an external package whose source code is in a sub folder `helloworld` in git repo `https://github.com/daisy-ycguo/wskdeploy-test`. `wskdeploy` will deploy the package according to the manifest and deployment files in the folder `helloworld` in `https://github.com/daisy-ycguo/wskdeploy-test`. When `wskdeploy` notices the dependency name `newpkg` is different from the original package name `myhelloworld`, it will create a package binding `newpkg` to `myhelloworld` so that all references to the assets in package `myhelloworld` can be referenced by the new package name `newpkg`, for example `newpkg/helloworld` in the sequence definition:
```
sequences:
  mySequence:
    actions: hellowhisk/greeting, hellowhisk/httpGet, newpkg/helloworld
```
- `myCloudant` is a binding to the existing system package `/whisk.system/cloudant`.  `wskdeploy` will create a package binding `myCloudant` to `/whisk.system/cloudant` with the input parameters and annotations.

To install this show case, try this:
```
$ wskdeploy -m ./tests/usecases/dependency/manifest.yaml
Deploying package opentest ... Done!
Deploying dependency newpkg ... 
Deploying package myhelloworld ... Done!
Deploying action myhelloworld/helloworld ... Done!
Done!
Deploying package binding newpkg ... Done!
Deploying dependency myCloudant ... 
Deploying package binding myCloudant ... Done!
Deploying dependency hellowhisk ... 
Deploying package hellowhisk ... Done!
Deploying action hellowhisk/httpGet ... Done!
Deploying action hellowhisk/greeting ... Done!
Done!
Deploying action opentest/mySequence ... Done!
Done!
Deploying rule myRule ... Done!
Deploying rule myCloudantRule ... Done!
```

To verify all the assets are deployed in OpenWhisk, try this:
```
$ wsk list
Entities in namespace: <your_name_space>
packages
/<your_name_space>/hellowhisk                                    private
/<your_name_space>/myCloudant                                    private
/<your_name_space>/newpkg                                        private
/<your_name_space>/myhelloworld                                  private
/<your_name_space>/opentest                                      private
actions
/<your_name_space>/opentest/mySequence                           private sequence
/<your_name_space>/hellowhisk/greeting                           private nodejs:6
/<your_name_space>/hellowhisk/httpGet                            private swift:3
/<your_name_space>/myhelloworld/helloworld                       private nodejs:6
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
To verify package `myhelloworld` is deployed, try this:
```
$ wsk package get --summary myhelloworld
package /<your_name_space>/myhelloworld
 action /<your_name_space>/myhelloworld/helloworld
```
To verify a package binding `newpkg` is installed, try this:
```
$ wsk package get --summary newpkg
package /<your_name_space>/newpkg
 action /<your_name_space>/newpkg/helloworld
```
To verify a package binding `mycloudant` is installed, try this:
```
$ wsk package get --summary myCloudant
```
To uninstall this show case, try this:
```
$ wskdeploy undeploy -m ./tests/usecases/dependency/manifest.yaml
```
