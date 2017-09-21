# ```wskdeploy``` utility by example
_A step-by-step guide for deploying Apache OpenWhisk applications using Package Manifest files._

This guide will walk you through how to describe OpenWhisk applications using the [OpenWhisk Packaging Specification](https://github.com/apache/incubator-openwhisk-wskdeploy/tree/master/specification#openwhisk-packaging-specification) and deploy them through the Whisk Deploy utility.  Please use the specification as the ultimate reference for all Manifest file grammar and syntax.

### Setting up your Host and Credentials
In order to deploy your OpenWhisk package, at minimum, the ```wskdeploy``` utility needs valid OpenWhisk APIHOST and AUTH variable to attempt deployment. Please read the [Configuring wskdeploy](wskdeploy_configuring.md#configuring-wskdeploy)

# Creating a "hello world" application

As with most language introductions, here we show a minimal "hello world" application as encoded in an OpenWhisk Package Manifest YAML file:

```yaml
package:
  name: hello_world_package
  version: 1.0
  license: Apache-2.0
  actions:
    hello_world:
      function: src/hello.js
```

where "hello.js" contains the following JavaScript code:
```javascript
function main(params) {
    msg = "Hello, " + params.name + " from " + params.place;
    return { payload:  msg };
}
```

#### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/manifest_hello_world.yaml)
- [hello.js](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/src/hello.js)


## Deploying "hello world"

You can actually deploy the "hello world" manifest from the incubator-openwhisk-wskdeploy project directory if you have downloaded it from GitHub:

```sh
$ ./wskdeploy -m docs/examples/manifest_hello_world.yaml
```

if you want to simply verify your manifest is able to read and parse your manifest file, you can add the ```-i``` or ```--allow-interactive``` flag:

```sh
$ ./wskdeploy -i -m docs/examples/manifest_hello_world.yaml
```

and the utility will stop, show you all the OpenWhisk package components it will deply and ask you if you want to deploy them or not.

```
         ____      ___                   _    _ _     _     _
        /\   \    / _ \ _ __   ___ _ __ | |  | | |__ (_)___| | __
   /\  /__\   \  | | | | '_ \ / _ \ '_ \| |  | | '_ \| / __| |/ /
  /  \____ \  /  | |_| | |_) |  __/ | | | |/\| | | | | \__ \   <
  \   \  /  \/    \___/| .__/ \___|_| |_|__/\__|_| |_|_|___/_|\_\
   \___\/              |_|

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

## Examining the "hello world" Manifest

The "hello world" example represents the minimum valid Manifest file which includes only the required parts of the Package and Action desciptors.

If you want to explore all possible Package and Action declarations (i.e., their schema) you can read:
- [Exploring Packages](wskdeploy_packages.md) - step-by-step guide on the **Package** schema.
- [Exploring Actions](wskdeploy_actions.md) - step-by-step guide on the **Action** schema.

# Debugging your Package Manifests

In addition to the normal output the wskdeploy utility provides, you may enable additional information that may further assist you in debugging. Please read the [Debugging Whisk Deploy](wskdeploy_debugging.md#debugging-wskdeploy) document.
