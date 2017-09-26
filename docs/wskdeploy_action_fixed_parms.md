# Actions

## Adding fixed input parameters

This example builds upon the previous “hello world” example and shows how fixed values can be supplied to the input parameters of an Action.

It shows how to:
- Declare input parameters on the action ‘```hello_world```’ using a single-line grammar.
- Add ‘```name```’ and ‘```place```’ as input parameters with the fixed values “```Sam```” and “```the Shire```” respectively.

### Manifest File

#### _Example: “Hello world” with fixed input values for ‘name’ and ‘place’_
```yaml
package:
  name: hello_world_package
  version: 1.0
  license: Apache-2.0
  actions:
    hello_world_fixed_parms:
      function: src/hello.js
      inputs:
        name: Sam
        place: the Shire
```

### Deploying
```sh
$ wskdeploy -m docs/examples/manifest_hello_world_fixed_parms.yaml
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world_fixed_parms --blocking
```

### Result
The invocation should return an 'ok' with a response that includes this result:
```sh
"result": {
    "greeting": "Hello, Sam from the Shire"
},
```

### Discussion

In this example:
- The value for the ‘```name```’ input parameter would be set to “```Sam```” and
- The value for the ‘```place```’ input parameter would be set to “```the Shire```”.
- The wskdeploy utility would infer both of their Types were '```string```'.

### Source code
The manifest file for this example can be found here:
- [manifest_hello_world_fixed_parms.yaml](examples/manifest_hello_world_fixed_parms.yaml)

### Specification
For convenience, the Actions and Parameters grammar can be found here:
- **[Actions](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_actions.md#actions)**
- **[Parameters](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_parameters.md#parameters)**

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_action_typed_parms.md#actions">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
