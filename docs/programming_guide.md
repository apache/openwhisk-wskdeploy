# Whisk Deploy Tool by Example
_A step-by-step guide for deploying Apache OpenWhisk applications using a Package Manifest file_

This guide, by example, will walk you through how to describe OpenWhisk applications using the OpenWhisk Packaging specification and deply them through the wskdeploy utility to any OpenWhisk Serverless provider.  

It assumes that you have setup and can run the ```wskdeploy``` as described in the [project README](https://github.com/apache/incubator-openwhisk-wskdeploy).  If so, then the utility will use the OpenWhisk APIHOST and AUTH variable values in your .wskprops file to attempt deployment.

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

Simply execute the wskdeploy binary against the directory you saved your "manifest.yaml" file in.

```
$ wskdeploy <my_directory>
```

_Please note that the wskdeploy utility will automatically look for any file named "manifest.yaml" or "manifest.yml" in the directory it is pointed to; however, the manifest file can be called anything as long as it has a .yaml or .yml extension and passed on the command line using the "-m" argument._

