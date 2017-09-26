# Actions

## The "Hello World" Action

As with most language introductions, here we show a simple "hello world" action as encoded in an OpenWhisk Package Manifest YAML file:

### Manifest file

#### _Example: “Hello world” using a NodeJS (JavaScript) action_
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
$ wskdeploy -m docs/examples/manifest_hello_world.yaml
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

The output parameter '```greeting```''s value included "_undefined_" values for the '```name```' and '```place```' input parameters as they were not provided in the manifest.

### Discussion

This "hello world" example represents the minimum valid Manifest file which includes only the required parts of the Package and Action descriptors.

In the above example,
- The Package and its Action were deployed to the user’s default namespace using the ‘package’ and 'action' names from the manifest.
  - ```/<default namespace>/hello_world_package/hello_world```
- The NodeJS default runtime (i.e., ```runtime: nodejs```) was selected automatically based upon the file extension '```.js```'' of the function's source file '```hello.js```'.


### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world.yaml](examples/manifest_hello_world.yaml)
- [hello.js](examples/src/hello.js)

### Specification
For convenience, the Packages and Actions grammar can be found here:
- **[Packages](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_packages.md#packages)**
- **[Actions](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_actions.md#actions)**

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
    <td><a href="wskdeploy_action_fixed_parms.md#actions">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
