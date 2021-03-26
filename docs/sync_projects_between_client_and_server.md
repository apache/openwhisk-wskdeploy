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

# How to Sync OpenWhisk Projects between Client and Server?

The answer is running Whisk Deploy in `sync` mode.

Whisk Deploy in `sync` mode, deploys all Apache OpenWhisk entities from the manifest file and attaches an annotation called `whisk-managed` to all entities from that manifest file. `whisk-managed` annotation contains following keys:

```
whisk-managed:
    projectName: <project-name>
    projectHash: SHA1("OpenWhisk " + <size_of_manifest_file> + "\0" + <contents_of_manifest_file>)
    projectDeps: <list of dependent packages>
    file: Relative path of manifest file on the file system
```

> Where the text “OpenWhisk” is a constant prefix and “\0” is the NULL character. The <size_of_manifest_file> and <contents_of_manifest_file> vary depending on the file.

Now, subsequent deployments of the same project in `sync` mode, calculates a new `projectHash` on client and compares it with the one on the server for every entity in that project. This comparison could lead us to following two scenarios:

* **Scenario 1:** If `projectHash` on client is same as `projectHash` on the server i.e. there were no changes in the project on the client side, the project on server side is left as is except wskdeploy redeploys all the entities from manifest file to capture any changes in deployment file.

* **Scenario 2:** If `projectHash` on client is different from `projectHash` on the server i.e. there were some changes in the project on the client side, `wskdeploy` redeploys all the entities from the manifest file on the client and updates their `projectHash` on the server. `wskdeploy` also searches all the entities including packages, actions, sequences, rules, and triggers which has the same `projectName` i.e. it belonged to the same project but has different `projectHash` i.e. its been deleted from the manifest file on the client.

Project name in the manifest file is mandatory to sync that project between the client and the server:

```
project:
    name: <project-name>
    packages:
        package1:
            ....
```

If for some reason, you want to avoid specifying project name in the manifest file, it can be specified on command line using `--projectname`.

OpenWhisk entities which are deployed using some other tool or automation and were part of the same project but has been deleted from that project now, are left unmodified. These kind of entities are classified as external entities and are not part of sync.

Undeployment of such project can be driven by `wskdeploy` with:

```
wskdeploy undeploy --projectname <project-name>
```

Lets look at a sample project to understand sync mode. Whisk Deploy GitHub repo has a sample project with many different manifest files to demonstrate sync mode.

#### Step 1: Deploy MyFirstManagedProject using `sync`:

```
~/wskdeploy sync -m manifest.yaml
Success: Deployment completed successfully.
```

Here is a list of entities deployed on OpenWhisk Server:

```
bx wsk list
Entities in namespace: guest
packages
/guest/ManagedPackage-2                                                private
/guest/ManagedPackage-1                                                private
actions
/guest/ManagedPackage-1/ManagedSequence-2                              private sequence
/guest/ManagedPackage-1/ManagedSequence-1                              private sequence
/guest/ManagedPackage-2/ManagedSequence-1                              private sequence
/guest/ManagedPackage-1/HelloWorld-2                                   private nodejs:default
/guest/ManagedPackage-1/HelloWorld-1                                   private nodejs:default
/guest/ManagedPackage-1/HelloWorld-3                                   private nodejs:default
/guest/ManagedPackage-2/HelloWorld-1                                   private nodejs:default
triggers
/guest/ManagedTrigger-2                                                private
/guest/ManagedTrigger-1                                                private
rules
/guest/ManagedRule-2                                                   private              active
/guest/ManagedRule-1                                                   private              active

```

`whisk-managed` annotation:

```
bx wsk package get ManagedPackage-1
...
    "name": "ManagedPackage-1",
    "annotations": [
        {
            "key": "whisk-managed",
            "value": {
                "file": "manifest.yaml",
                "projectDeps": [],
                "projectHash": "d164caed3ab86106495b3963232c4840429df8ea",
                "projectName": "MyFirstManagedProject"
            }
        }
```

#### Step 2: Sync Client and Server — Deleting ManagedPackage-2

```
~/wskdeploy sync -m 00-manifest-minus-second-package.yaml
Deployment completed successfully.
```

Here is a list of entities deployed on OpenWhisk Server:

```
bx wsk list
Entities in namespace: guest
packages
/guest/ManagedPackage-1                                                private
actions
/guest/ManagedPackage-1/ManagedSequence-2                              private sequence
/guest/ManagedPackage-1/ManagedSequence-1                              private sequence
/guest/ManagedPackage-1/HelloWorld-3                                   private nodejs:default
/guest/ManagedPackage-1/HelloWorld-2                                   private nodejs:default
/guest/ManagedPackage-1/HelloWorld-1                                   private nodejs:default
triggers
/guest/ManagedTrigger-1                                                private
rules
/guest/ManagedRule-1                                                   private              active

```

#### Step 3: Sync Client and Server — Deleting Sequence ManagedSequence-2

```
~/wskdeploy sync -m 01-manifest-minus-sequence-2.yaml
Deployment completed successfully.
```

Here is a list of entities deployed on OpenWhisk Server:

```
bx wsk list
Entities in namespace: guest
packages
/guest/ManagedPackage-1                                                private
actions
/guest/ManagedPackage-1/ManagedSequence-1                              private sequence
/guest/ManagedPackage-1/HelloWorld-3                                   private nodejs:default
/guest/ManagedPackage-1/HelloWorld-2                                   private nodejs:default
/guest/ManagedPackage-1/HelloWorld-1                                   private nodejs:default
triggers
/guest/ManagedTrigger-1                                                private
rules
/guest/ManagedRule-1                                                   private              active
```

#### Step 4: Sync Client and Server — Deleting Action Helloworld-3

```
~/wskdeploy sync -m 02-manifest-minus-action-3.yaml --managed
Deployment completed successfully.
```

Here is a list of entities deployed on OpenWhisk Server:

```
bx wsk list
Entities in namespace: guest
packages
/guest/ManagedPackage-1                                                private
actions
/guest/ManagedPackage-1/ManagedSequence-1                              private sequence
/guest/ManagedPackage-1/HelloWorld-2                                   private nodejs:default
/guest/ManagedPackage-1/HelloWorld-1                                   private nodejs:default
triggers
/guest/ManagedTrigger-1                                                private
rules
/guest/ManagedRule-1                                                   private              active
```

#### Step 5: Sync Client and Server — Deleting ManagedPackage-1 (deleting the entire project)

```
~/wskdeploy undeploy --projectname MyFirstManagedProject
```

## Sync Projects with Dependencies

Now, a question is how does `sync` work in case of a project which has dependencies. Whisk Deploy supports two different types of dependencies which are explained [here](../tests/usecases/dependency/README.md). Let's look at how projects are synced with each type of dependency.

### GitHub Dependency

[Here](../tests/src/integration/managed-deployment/07-manifest-with-dependency.yaml) is a sample project definition with single dependency:

```
project:
    name: MyManagedProjectWithSingleDependency
    packages:
        Extension1:
            dependencies:
                helloworlds:
                    location: github.com/apache/openwhisk-test/packages/helloworlds
```

After deploying this project with `wskdeploy sync -m manifest.yaml`, package `Extension2` has following annotation:

```
bx wsk package get Extension1
ok: got package Extension1
{
    "namespace": "guest",
    "name": "Extension1",
    "version": "0.0.2",
    "publish": false,
    "annotations": [
        {
            "key": "whisk-managed",
            "value": {
                "file": "manifest.yaml",
                "projectDeps": [
                    {
                        "key": "/guest/helloworlds",
                        "value": {
                            "file": ".../src/github.com/apache/openwhisk-wskdeploy/Packages/helloworlds-master/packages/helloworlds/manifest.yaml",
                            "projectDeps": [],
                            "projectHash": "0ae33344f7885df01aefd053e94a4eb3675ef72d",
                            "projectName": "HelloWorlds"
                        }
                    }
                ],
                "projectHash": "6636904777a298993127202d509d3d405e10592c",
                "projectName": "MyManagedProjectWithSingleDependency"
            }
        }
    ],
    "binding": {}
}
```

And, dependent package `helloworlds` has following annotation:

```
bx wsk package get helloworlds
...
    "annotations": [
        {
            "key": "whisk-managed",
            "value": {
                "file": ".../src/github.com/apache/openwhisk-wskdeploy/Packages/helloworlds-master/packages/helloworlds/manifest.yaml",
                "projectDeps": [],
                "projectHash": "0ae33344f7885df01aefd053e94a4eb3675ef72d",
                "projectName": "HelloWorlds"
            }
        }
    ],
    "binding": {},
...
```

Here, on server side, package `Extension1` is showing dependency on `helloworlds` using `whisk-managed` annotation.

Now, let's add one more dependency in our project with [manifest.yaml](../tests/src/integration/managed-deployment/06-manifest-with-single-dependency.yaml):

```
project:
    name: MyManagedProjectWithDependency
    packages:
        Extension2:
            dependencies:
                helloworlds:
                    location: github.com/apache/openwhisk-test/packages/helloworlds
                custom-hellowhisk:
                    location: github.com/apache/openwhisk-test/packages/hellowhisk
...
```

After deploying this project with `wskdeploy sync -m manifest.yaml`, package `Extension2` has following annotation:

```
bx wsk package get Extension2
ok: got package Extension2
{
    "namespace": "guest",
    "name": "Extension2",
    "version": "0.0.2",
    "publish": false,
    "annotations": [
        {
            "key": "whisk-managed",
            "value": {
                "file": "manifest.yaml",
                "projectDeps": [
                    {
                        "key": "/guest/helloworlds",
                        "value": {
                            "file": ".../src/github.com/apache/openwhisk-wskdeploy/Packages/helloworlds-master/packages/helloworlds/manifest.yaml",
                            "projectDeps": [],
                            "projectHash": "0ae33344f7885df01aefd053e94a4eb3675ef72d",
                            "projectName": "HelloWorlds"
                        }
                    },
                    {
                        "key": "/guest/hellowhisk",
                        "value": {
                            "file": "../src/github.com/apache/openwhisk-wskdeploy/Packages/custom-hellowhisk-master/packages/hellowhisk/manifest.yaml",
                            "projectDeps": [],
                            "projectHash": "8970b2e820631322ae630e9366b2342ab3b67a57",
                            "projectName": "HelloWhisk"
                        }
                    }
                ],
                "projectHash": "554e0aaef1ed2ff0dcf7f68d0c4a8f825d1a49c7",
                "projectName": "MyManagedProjectWithDependency"
            }
        }
    ],
    "binding": {}
}
```

Here, on server side, package `Extension2` is showing both dependencies (1) `helloworlds` and (2) `hellowhisk`.

Now, when we try to undeploy `Extension2` using `~/wskdeploy undeploy --projectname MyManagedProjectWithDependency`, only one dependency `hellowhisk` is deleted and the other dependency `helloworlds` is not deleted as its referred by `Extension1`.

> **Note:** Here the dependent package belongs to project `HelloWorlds` which has a single package `helloworlds`. Whisk Deploy does not support having dependencies with multiple packages. The reason behind this limitation is when a dependency is specified in manifest file using label (`custom-hellowhisk`) which is different than the dependent package name (`hellowhisk`), `custome-hellowhisk` is created as a binding to `hellowhisk`. OpenWhisk does not support creating a package which is bound to multiple packages.

### Package Binding

When a project has dependency to any already deployed package/utility, for example, `/whisk.system/utils` or any `/<namespace>/<package>`:

```
project:
    name: MyManagedProjectWithWhiskSystemDependency
    packages:
        Extension3:
            dependencies:
                whiskUtility:
                    location: /whisk.system/utils
            triggers:
                triggerInExtension3:
            rules:
                ruleInExtension3:
                    trigger: triggerInExtension3
                    action: whiskUtility/sort
```

For these kind of references, dependent package (`/whisk.system/utils` in our example) is not listed as one of the dependencies since the dependent package is part of a catalog and was pre-installed therefore cannot be managed by the current project:

```
bx wsk package get Extension3
ok: got package Extension3
{
    "namespace": "guest",
    "name": "Extension3",
    "version": "0.0.2",
    "publish": false,
    "annotations": [
        {
            "key": "whisk-managed",
            "value": {
                "file": "manifest.yaml",
                "projectDeps": [],
                "projectHash": "f5ec911a4e3f011343e4bb0af1900dda3279be75",
                "projectName": "MyManagedProjectWithWhiskSystemDependency"
            }
        }
    ],
    "binding": {}
}
```

Here, `whiskUtility` is a binding to package `/whisk.system/utils` and part of the current project. This signifies that `Extension3` has a dependency on a pre-installed package `/whisk.system/utils`:

```
bx wsk package get whiskUtility
ok: got package whiskUtility
{
    "namespace": "guest",
    "name": "whiskUtility",
    "version": "0.0.1",
    "publish": false,
    "annotations": [
        {
            "key": "whisk-managed",
            "value": {
                "file": "tests/src/integration/managed-deployment/08-manifest-with-dependencies-on-whisk-system.yaml",
                "projectDeps": [],
                "projectHash": "f5ec911a4e3f011343e4bb0af1900dda3279be75",
                "projectName": "MyManagedProjectWithWhiskSystemDependency"
            }
        },
        {
            "key": "binding",
            "value": {
                "name": "utils",
                "namespace": "whisk.system"
            }
        }
    ],
...
```

Now, when we undeploy `Extension3` with `~/wskdeploy undeploy --projectname MyManagedProjectWithWhiskSystemDependency`, nothing changes under `/whisk.system`. `utils` package under `/whisk.system` remains undeployed. Similarly, with an undeployment of `Extension3`, `/<namespace>/<package>` is not deleted.
