## The "Hello World" Action

As with most language introductions, here we show a minimal "hello world" action as encoded in an OpenWhisk Package Manifest YAML file:

### Manifest file
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

### Deploying

You can actually deploy the "hello world" manifest from the incubator-openwhisk-wskdeploy project directory if you have downloaded it from GitHub:

```sh
$ ./wskdeploy -m docs/examples/manifest_helloworld.yaml
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world --blocking
```

### Result
The invocation should return an 'ok' with a response that includes this result:

```sh
  "result": {
      "greeting": "Hello, undefined from undefined"
  },
```

### Discussion

This "hello world" example represents the minimum valid Manifest file which includes only the required parts of the Package and Action descriptors.

In the above example,
- The Package and its Action were deployed to the user’s default namespace using the ‘package’ name.
  - ```/<default namespace>/hello_world_package/hello_world```
- The default runtime (i.e., ```runtime: nodejs```) was automatically selected based upon the ‘```.js```’ extension on the Action function’s source file '```hello.js```'.


- The package and action were deployed to the user’s default namespace:
  - ```/<default namespace>/hello_world_package/hello_world```
- The NodeJS runtime was selected automatically based upon the file extension '```.js```'' of the function's source file '```hello.js```'.
- The parameters 'name' and 'place' in this example are not set to any default values; therefore 'greeting' will appear a follows:
  - ```"greeting": "Hello, undefined from undefined"```
  we will explore setting these values in later examples.

### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_helloworld.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/manifest_hello_world_1.yaml)
- [hello.js](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/src/hello.js)

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="wskdeploy_package_minimal.md#packages">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
