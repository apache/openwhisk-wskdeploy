# Packages

The wskdeploy utility works primarily with the OpenWhisk **Package** resource as described in the [OpenWhisk Packaging Specification](https://github.com/apache/incubator-openwhisk-wskdeploy/tree/master/specification#openwhisk-packaging-specification).

For convenience, the schema and grammar for declaring a **Package** can be found here:
[https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_packages.md#packages](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_packages.md#packages)

## Creating an OpenWhisk Package

### Start with a Package Manifest (YAML) file
The wskdeploy utility mainly uses a single file, that uses a YAML syntax, called a "Package Manifest", to describe all the OpenWhisk components that make up your OpenWhisk Package including Actions, Triggers, Rules, etc.

The minimal manifest file would include only a package declaration, a version number and a license for the package:
```
package:
  name: hello_world_package
  version: 1.0
  license: Apache-2.0
```

Save this into a file called "manifest.yaml" in a directory of your choice.

### Executing the wskdeploy utility
Simply execute the wskdeploy binary against the directory you saved your "manifest.yaml" file in by pointing it to the package location using the ```-p``` flag.

```sh
$ wskdeploy -p <my_directory>
```
wskdeploy will automatically look for any file named "manifest.yaml" or "manifest.yml" in the directory it is pointed; however, the manifest file can be called anything as long as it has a .yaml or .yml extension and passed on the command line using the ```-m``` flag.

For example, if you called your manifest "my_pkg_manifest.yml" you could simply provide the manifest file name as follows:
```sh
$ wskdeploy -p <my_directory> -m my_pkg_manifest.yaml
```
### Interactive mode

if you want to simply verify your manifest is able to read and parse your manifest file, you can add the ```-i``` or ```--allow-interactive``` flag:

```sh
$ ./wskdeploy -i -m docs/examples/manifest_hello_world.yaml
```

and the utility will stop, show you all the OpenWhisk package components it will deploy and ask you if you want to deploy them or not.

```sh
Package:
  name: hello_world_package
  bindings:

  * action: hello_world
    bindings:
    annotations:

  Triggers:
  Rules:

Do you really want to deploy this? (y/N):
```

Now you can see what will be deployed

### What if ```wskdeploy``` finds an error in my manifest?

- The ```wskdeploy``` utility will not attempt to deploy a package if an error in the manifest is detected, but will report as much information as it can to help you locate the error in the YAML file.

### What if ```wskdeploy``` encounters an error during deployment?

-  The ```wskdeploy``` utility will cease deploying as soon as it receives an error from the target platform and display what error information it receives to you.
- It will then attempt to undeploy any entities that it attempted to deploy.
  - If "interactive mode" was used to deploy then you will be prompted to confirm you wish to undeploy.

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="programming_guide.md#guided-examples">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Index</a></td>
    <td><a href="wskdeploy_hello_world.md#creating-a-hello-world-package">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
