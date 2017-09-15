# Working with packages

## Creating an OpenWhisk Package

### Start with a Package Manifest (YAML) file
The wskdeploy utility maninly uses a single file, that uses a YAML syntax, called a "Package Manifest", to describe all the OpenWhisk components that make up your OpenWhisk Package including Actions, Triggers, Rules, etc.

The minimal manifest file would include only a package declaration, a version number and a license for the package:
```
package:
  name: helloworld
  version: 1.0
  license: Apache-2.0
```

Save this into a file called "manifest.yaml" in a directory of your choice.

### Executing the wskdeploy utility
Simply execute the wskdeploy binary against the directory you saved your "manifest.yaml" file in by pointing it to the package location using the ```-p``` flag.

```
$ wskdeploy -p <my_directory>
```
wskdeploy will automatically look for any file named "manifest.yaml" or "manifest.yml" in the directory it is pointed; however, the manifest file can be called anything as long as it has a .yaml or .yml extension and passed on the command line using the ```-m``` flag.

For example, if you called your manifest "my_pkg_manifest.yml" you could simply provide the manifest file name as follows:
```
$ wskdeploy -p <my_directory> -m my_pkg_manifest.yaml
```
