## Creating a "hello world" package

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

### Deploying "hello world"

You can actually deploy the "hello world" manifest from the incubator-openwhisk-wskdeploy project directory if you have downloaded it from GitHub:

```sh
$ ./wskdeploy -m docs/examples/manifest_hello_world.yaml
```

### Running
```sh
$ wsk action invoke hello_world_package/hello_world --blocking
```
should return an 'ok' with a response that looks like:
```sh
    "response": {
        "result": {
            "greeting": "Hello, undefined from undefined"
        },
        "status": "success",
        "success": true
    },
```

### Discussion

This "hello world" example represents the minimum valid Manifest file which includes only the required parts of the Package and Action descriptors.

- The package and action were deployed to the user’s default namespace:
  - ```/<default namespace>/hello_world_package/hello_world```
- The NodeJS runtime was selected automatically based upon the file extension '```.js```'' of the function's source file '```hello.js```'.
- The parameters 'name' and 'place' in this example are not set to any default values; therefore 'greeting' will appear a follows:
  - ```"greeting": "Hello, undefined from undefined"```
  we will explore setting these values in later examples.

#### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/manifest_hello_world.yaml)
- [hello.js](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/src/hello.js)

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="wskdeploy_hello_world.md#creating-a-hello-world-package">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Index</a></td>
    <td><a href="">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
