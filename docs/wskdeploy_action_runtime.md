# Actions

## Setting your Function's runtime

In the previous ["Hello world" example](), the ```wskdeploy``` utility used the file extension of the function "hello.js" to select the most current NodeJS runtime.

In most cases, allowing the utility to select the runtime works well using this implicit method. However, perhaps your code is dependent on a specific version of a language runtime and you want to explicitly set it?

This example shows how to:
- declare an explicit runtime for your Action's function.

### Manifest file
#### _Example: explicit selection of the NodeJS version 6 runtime_
```yaml
package:
  name: hello_world_package
  version: 1.0
  license: Apache-2.0
  actions:
    hello_world_runtime:
      function: src/hello.js
      runtime: nodejs@6
```

### Deploying

```sh
$ wskdeploy -m docs/examples/manifest_hello_world_runtime.yaml
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world_runtime --blocking
```

### Result
The invocation should return an 'ok' with a response that includes this result:

```json
"result": {
    "greeting": "Hello, undefined from undefined"
},
```

### Discussion

In the above example,
- The value for the '```runtime```' key was a valid name and version supported by OpenWhisk
  - Please see the current supported list here: **[Actions - Valid runtime names](../specification/html/spec_actions.md#valid-runtime-names)**

#### Runtime mismatch for function
If the language runtime you requested is not compatible with the function's language, then you will receive an error response when invoking the Action.  For example, the following manifest indicates a JavaScript (.js) function, but the runtime selected as "python":
```yaml
package:
  name: hello_world_package
  ...
  actions:
    hello_world_runtime:
      function: src/hello.js
      runtime: python
```

The result would a "failure" with a failed response:
```json
"response": {
    "result": {
        "error": "The action failed to generate or locate a binary. See logs for details."
    },
    "status": "action developer error",
    "success": false
```

### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world_runtime.yaml](examples/manifest_hello_world_runtime.yaml)
- [hello.js](examples/src/hello.js)

### Specification
For convenience, the Packages and Actions grammar can be found here:
- **[Actions](../specification/html/spec_actions.md#actions)**
- **[Actions - Valid runtime names](../specification/html/spec_actions.md#valid-runtime-names)**

### Notes

- If you use the following curl command, you can see the latest runtimes and version supported by the IBM Cloud Functions platform:
  - ```curl -k https://openwhisk.ng.bluemix.net```

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="wskdeploy_action_helloworld.md#actions">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_action_fixed_parms.md#actions">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
