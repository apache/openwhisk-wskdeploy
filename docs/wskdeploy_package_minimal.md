# Packages

The wskdeploy utility works primarily with the OpenWhisk **Package** resource as described in the [OpenWhisk Packaging Specification](https://github.com/apache/incubator-openwhisk-wskdeploy/tree/master/specification#openwhisk-packaging-specification).

## Creating a minimal OpenWhisk Package

### Start with a Package Manifest (YAML) file
The ```wskdeploy``` utility mainly uses a single YAML file, called a "Package Manifest", to describe all the OpenWhisk components that make up your OpenWhisk Package including Actions, Triggers, Rules, etc.

### Manifest

The minimal manifest file would include only a package declaration, a version number and a license for the package:
```
package:
  name: hello_world_package
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
$ wskdeploy -p <my_directory> -m docs/examples/manifest_helloworld.yaml
```

#### Interactive mode

If you want to simply verify your manifest file can be read and parsed properly before deploying, you can add the ```-i``` or ```--allow-interactive``` flag:

```sh
$ ./wskdeploy -i -m docs/examples/manifest_helloworld.yaml
```

and the utility will stop, show you all the OpenWhisk package components it will deploy from your manifest and ask you if you want to deploy them or not.

```sh
Package:
  name: hello_world_package
  bindings:
  annotations:
Triggers:
Rules:

Do you really want to deploy this? (y/N):
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
- [manifest_package_minimal.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/manifest_package_minimal.yaml)

### Specification
For convenience, the Packages grammar can be found here:
- **[Packages](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_packages.md#packages)**

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
    <td><a href="wskdeploy_action_helloworld.md">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
