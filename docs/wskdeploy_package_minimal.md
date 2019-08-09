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

# Packages

The wskdeploy utility works primarily with the OpenWhisk **Package** resource as described in the [OpenWhisk Packaging Specification](https://github.com/apache/openwhisk-wskdeploy/tree/master/specification#openwhisk-packaging-specification).

## Creating a minimal OpenWhisk Package

### Start with a Package Manifest (YAML) file
The ```wskdeploy``` utility mainly uses a single YAML file, called a "Package Manifest", to describe all the OpenWhisk components that make up your OpenWhisk Package including Actions, Triggers, Rules, etc.

### Manifest

The minimal manifest file would include only a package declaration, a version number and a license for the package:
```
packages:
  hello_world_package:
    version: 1.0
    license: Apache-2.0
```

Save this into a file called ```"manifest.yaml"``` in a directory of your choice.

### Deploying

#### using the project path
Simply execute the ```wskdeploy``` utility binary against the directory you saved your "manifest.yaml" file in by pointing it to the package location using the ```-p``` flag.

```sh
$ wskdeploy -p <my_directory>
```
wskdeploy will automatically look for any file named ```"manifest.yaml"``` or ```"manifest.yml"``` in the directory it is pointed; however, the _manifest file can be called anything_ as long as it has a .yaml or .yml extension and passed on the command line using the ```-m``` flag.

#### using a named manifest file
If you called your manifest "manifest_helloworld.yaml" (not using the default manifest.yaml name) and placed it in a directory below your project directory, you could simply provide the project-relative path to the manifest file as follows:
```sh
$ wskdeploy -p <my_directory> -m docs/examples/manifest_package_minimal.yaml
```

#### Dry Run mode

If you want to simply verify your manifest file can be read and parsed properly before deploying, you can add the ```--preview``` flag:

```sh
$ ./wskdeploy --preview -m docs/examples/manifest_package_minimal.yaml
```

and the utility will stop, show you all the OpenWhisk package components it will deploy from your manifest.

```sh
Package:
Name: hello_world_package
  bindings:
  annotations:
Triggers:
Rules:
```

### Result
You can use the Whisk CLI to confirm your package was created:
```sh
$ wsk package list

packages
/<default_namespace>/hello_world_package

```

### Discussion

- The package '```hello_world_package```' was created in the user's default namespace at their target OpenWhisk provider.
- Currently, OpenWhisk does not yet support the '```version```' or '```license```' fields, but are planned for future versions.  However, their values will be validated against the specification.

#### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_package_minimal.yaml](https://github.com/apache/openwhisk-wskdeploy/blob/master/docs/examples/manifest_package_minimal.yaml)

### Specification
For convenience, the Packages grammar can be found here:
- **[Packages](https://github.com/apache/openwhisk-wskdeploy/blob/master/specification/html/spec_packages.md#packages)**

---

<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="programming_guide.md">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_action_helloworld.md#actions">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
