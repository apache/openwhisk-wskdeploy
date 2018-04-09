# Using `wskdeploy` for exporting `OpenWhisk` assets

`wskdeploy export` can be used to export `OpenWhisk` assets previously deployed as a *managed project* via `wskdeploy -m manifest.yaml`. 
`wskdeploy export` will create a manifest for the managed project assets and separate manifests for each managed project that this managed 
project depends upon, if such dependencies exist and have been described in `manifest.yml` when the managed project has been initially deployed.
The manifest(s) resulting from executing `wskdeploy export` can be later used for deploying at a different `OpenWhisk` instance. The code of actions, which are defined in the packages of the exported project will be saved into folders with the names being the names of the package, the actions belong to.

## Use Cases

### Copy `OpenWhisk` Assets

One common scenario, in which the export feature is useful, is populating a newly installed `OpenWhisk` instance with assets from a 
another `OpenWhisk` instance. One might consider a scenario, in which an `OpenWhisk` instance in installed on premises with another `OpenWhisk`
instance residing in the cloud. Consider, for example using one `OpenWhisk` instance on premises and another one in the cloud (e.g., the 
second `OpenWhisk` instance can be [IBM Cloud Functions](https://console.bluemix.net/openwhisk/)). A fairly common scenario is that a developer
will need to deploy assets from the cloud `OpenWhisk` instance on the on premises one and vice-versa. 

### `OpenWhisk` at the Edge

In a variety of IoT and other edge computing scenarios, such as running Virtual Network Functions (VNF) as `OpenWhisk` actions in "edge Data Centers" 
embedded with 5G-MEDIA infrastructure (as pioneered in [5G-MEDIA EU H2020 project](http://www.5gmedia.eu/), there is a need to distribute 
`OpenWhisk` assets developed centrally in the cloud (e.g., [IBM Cloud Functions](https://console.bluemix.net/openwhisk/)) to multiple 
`OpenWhisk` instances running at the edge Data Centers. Again, `wskdeploy export` is handy as a basic tool that allows to automate this. 
management task.

## Basic Usage by Example

Consider a simple manifest file [manifest_lib1.yml](https://github.com/davidbreitgand/incubator-openwhisk-wskdeploy/blob/add-export-doc2readme/tests/src/integration/export/manifest_lib1.yaml) for a sample project `lib1`.
The project contains a single `lib1_package` package that comprise three actions (the code of the action in this simple example is the
same for all three, but the action names are different).


### Step 1: deploy `lib1` as a managed project 

```sh
$ ./wskdeploy sync -m manifest_lib1.yaml
```

### Step 2: validate `lib1` deployment

```sh
$ ./wsk package get lib1_package
```

<details><summary>You should see an output similar to this one (clickable):</summary>
<p>

```json
ok: got package lib1_package
{
    "namespace": "your_namespace",
    "name": "lib1_package",
    "version": "0.0.2",
    "publish": false,
    "annotations": [
        {
            "key": "whisk-managed",
            "value": {
                "file": "/root/go_projects/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/export/manifest_lib1.yaml",
                "projectDeps": [],
                "projectHash": "80eec5f8e3ee874e22bdacb76aa4cc69aad459c1",
                "projectName": "lib1"
            }
        }
    ],
    "binding": {},
    "actions": [
        {
            "name": "lib1_greeting3",
            "version": "0.0.1",
            "annotations": [
                {
                    "key": "whisk-managed",
                    "value": {
                        "file": "/root/go_projects/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/export/manifest_lib1.yaml",
                        "projectDeps": [],
                        "projectHash": "80eec5f8e3ee874e22bdacb76aa4cc69aad459c1",
                        "projectName": "lib1"
                    }
                },
                {
                    "key": "exec",
                    "value": "nodejs:6"
                }
            ]
        },
        {
            "name": "lib1_greeting2",
            "version": "0.0.1",
            "annotations": [
                {
                    "key": "whisk-managed",
                    "value": {
                        "file": "/root/go_projects/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/export/manifest_lib1.yaml",
                        "projectDeps": [],
                        "projectHash": "80eec5f8e3ee874e22bdacb76aa4cc69aad459c1",
                        "projectName": "lib1"
                    }
                },
                {
                    "key": "exec",
                    "value": "nodejs:6"
                }
            ]
        },
        {
            "name": "lib1_greeting1",
            "version": "0.0.1",
            "annotations": [
                {
                    "key": "whisk-managed",
                    "value": {
                        "file": "/root/go_projects/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/export/manifest_lib1.yaml",
                        "projectDeps": [],
                        "projectHash": "80eec5f8e3ee874e22bdacb76aa4cc69aad459c1",
                        "projectName": "lib1"
                    }
                },
                {
                    "key": "exec",
                    "value": "nodejs:6"
                }
            ]
        }
    ]
}
```
</p>
</details>

### Step 3: Export the newly deployed `lib1` project

Exporting into the current directory.

```sh
$ ./wskdeploy export --projectname lib1 -m my_new_lib1_manifest.yaml
```


### Step 4: Inspect the newly exported manifest (`my_new_lib1_manifest.yaml`) 

<details><summary>You should see an output similar to this one (clickable):</summary>

```yaml
project:
  name: lib1
  namespace: ""
  credential: ""
  apiHost: ""
  apigwAccessToken: ""
  version: ""
  packages: {}
packages:
  lib1_package:
    name: lib1_package
    version: 0.0.2
    license: ""
    dependencies: {}
    namespace: your_namespace
    credential: ""
    apiHost: ""
    apigwAccessToken: ""
    actions:
      lib1_greeting1:
        name: lib1_greeting1
        location: ""
        version: 0.0.1
        function: lib1_package/lib1_greeting1.js
        code: ""
        runtime: nodejs:6
        namespace: your_namespace/lib1_package
        credential: ""
        exposedUrl: ""
        web-export: ""
        main: ""
        limits: null
        inputs: {}
        outputs: {}
        annotations:
          exec: nodejs:6
      lib1_greeting2:
        name: lib1_greeting2
        location: ""
        version: 0.0.1
        function: lib1_package/lib1_greeting2.js
        code: ""
        runtime: nodejs:6
        namespace: your_namespace/lib1_package
        credential: ""
        exposedUrl: ""
        web-export: ""
        main: ""
        limits: null
        inputs: {}
        outputs: {}
        annotations:
          exec: nodejs:6
      lib1_greeting3:
        name: lib1_greeting3
        location: ""
        version: 0.0.1
        function: lib1_package/lib1_greeting3.js
        code: ""
        runtime: nodejs:6
        namespace: your_namespace/lib1_package
        credential: ""
        exposedUrl: ""
        web-export: ""
        main: ""
        limits: null
        inputs: {}
        outputs: {}
        annotations:
          exec: nodejs:6
    triggers: {}
    feeds: {}
    rules: {}
    inputs: {}
    sequences: {}
    apis: {}
filepath: ""

```
</details>

### Step 5: Inspect the newly exported package 

The code of the actions defined in the packages comprising the exported project will be saved into the folders named after
the respective packages. The packages' folders will be created in the same directory into which the manifest file of the
project is exported. Let's examine the current directory into which `my_new_lib1_manifest.yaml` was exported in Step 3 above.

```sh
$ ls -al lib1_package
```

<details><summary>You should see an output similar to this one (clickable):</summary>
<pre>
drwxr-xr-x  2 root root 4096 Apr  8 22:52 .
drwxr-xr-x 26 root root 4096 Apr  8 23:38 ..
-rw-r--r--  1 root root  331 Apr  8 22:59 lib1_greeting1.js
-rw-r--r--  1 root root  331 Apr  8 22:58 lib1_greeting2.js
-rw-r--r--  1 root root  331 Apr  8 22:58 lib1_greeting3.js
</pre>
</details>

## Advanced Usage 

The dependencies mechanism allows to express a project structure, in which one project uses another project as a library. Also dependencies can be defined for multiple projects. Consider a project `lib2` with the manifest [manifest_lib2.yml](https://github.com/davidbreitgand/incubator-openwhisk-wskdeploy/blob/add-export-doc2readme/tests/src/integration/export/manifest_lib2.yaml) and a project `EXT_PROJECT` with the manifest [manifest_ext.yml](https://github.com/davidbreitgand/incubator-openwhisk-wskdeploy/blob/add-export-doc2readme/tests/src/integration/export/manifest_ext.yaml). `EXT_PROJECT` (stands for _extending project_) uses actions from both package `lib1_package` (defined in the `lib1` project) and `lib2_package` (defined in the `lib2` project) in order to define rules specific to `EXT_PROJECT`.

`wskdeploy export` will automatically export both `lib1` and `lib2` along with `EXT_PROJECT`. Each exported project will have a manifest and package folder structure as explained [above](#-Basic-Usage-by-Example). 

## Notes

+ Currently, dependencies are not treated recursively when exporting a project. Hence, only the upper level projects defined as a dependency will be exported.
+ To redeploy a project with dependencies, a user should first deploy dependency projects projects (`lib1` and `lib2` in our example) and only after that, `EXT_PROJECT` can be deployed successfully. 
+ `wskdeploy export` does not check for circular dependencies. In case of circular dependencies specified by the user, `wskdeploy`'s behavior is undefined.


