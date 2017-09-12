# Whisk Deploy Tool by Example
_A step-by-step guide for deploying Apache OpenWhisk applications using a Package Manifest file_

This guide, by example, will walk you through how to describe OpenWhisk applications using the OpenWhisk Packaging specification and deply them through the wskdeploy utility to any OpenWhisk Serverless provider.  

## Describe your OpenWhisk Actions in a Package Manifest

The wskdeploy utility maninly uses a single file, that uses a YAML syntax, called a "Package Manifest", to describe all the OpenWhisk components that make up your OpenWhisk Package including Actions, Triggers, Rules, etc. 

The most 
```

```

_Please note that the _wskdeploy utility will automatically look for any file named "manifest.yaml" or "manifest.yml" in the directory it is pointed to; however, the manifest file can be called anything as long as it has a .yaml or .yml extension and passed on the command line using the "-m" argument._

