<!--
#
# Licensed to the Apache Software Foundation (ASF) under one or more contributor
# license agreements.  See the NOTICE file distributed with this work for additional
# information regarding copyright ownership.  The ASF licenses this file to you
# under the Apache License, Version 2.0 (the # "License"); you may not use this
# file except in compliance with the License.  You may obtain a copy of the License
# at:
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed
# under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
# CONDITIONS OF ANY KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations under the License.
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

Now, subsequent deployments of the same project in `sync` mode, calculates a new `projectHash` on client and compares it with the one on the server for every entity in that project. This comparision could lead us to following two scenarios:

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
/guest/ManagedPackage-1/HelloWorld-2                                   private nodejs:6
/guest/ManagedPackage-1/HelloWorld-1                                   private nodejs:6
/guest/ManagedPackage-1/HelloWorld-3                                   private nodejs:6
/guest/ManagedPackage-2/HelloWorld-1                                   private nodejs:6
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
/guest/ManagedPackage-1/HelloWorld-3                                   private nodejs:6
/guest/ManagedPackage-1/HelloWorld-2                                   private nodejs:6
/guest/ManagedPackage-1/HelloWorld-1                                   private nodejs:6
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
/guest/ManagedPackage-1/HelloWorld-3                                   private nodejs:6
/guest/ManagedPackage-1/HelloWorld-2                                   private nodejs:6
/guest/ManagedPackage-1/HelloWorld-1                                   private nodejs:6
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
/guest/ManagedPackage-1/HelloWorld-2                                   private nodejs:6
/guest/ManagedPackage-1/HelloWorld-1                                   private nodejs:6
triggers
/guest/ManagedTrigger-1                                                private
rules
/guest/ManagedRule-1                                                   private              active
```

#### Step 5: Sync Client and Server — Deleting ManagedPackage-1 (deleting the entire project)

```
~/wskdeploy undeploy --projectname MyFirstManagedProject
```



