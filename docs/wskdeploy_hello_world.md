## Creating a "hello world" application

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
    return { greeting:  msg };
}
```

#### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/manifest_hello_world.yaml)
- [hello.js](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/src/hello.js)


### Deploying "hello world"

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

### Examining the "hello world" Manifest

The "hello world" example represents the minimum valid Manifest file which includes only the required parts of the Package and Action desciptors.
