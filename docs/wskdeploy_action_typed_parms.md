# Actions

## Typed input and output parameters

This example extends the 'Hello world' example with typed input and output Parameters.

It shows how to:
- Declare input and output parameters on the action '```hello_world```' using a simple, single-line format.
- Add two input parameters, '```name```' and '```place```', both of type '```string```' to the '```hello_world```' action.
- Add an output parameter, '```greeting```' of type string to the '```hello_world```' action.

### Manifest File

#### Example 3: 'Hello world' with typed input and output parameter declarations
```yaml
package:
  name: hello_world_package
  ... # Package keys omitted for brevity
  actions:
    hello_world_typed_parms:
      function: src/hello.js
      inputs:
        name: string
        place: string
      outputs:
        greeting: string
```

### Deploying
```sh
$ wskdeploy -m docs/examples/manifest_hello_world_typed_parms.yaml
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world_typed_parms --blocking
```

### Result
```sh
"result": {
    "greeting": "Hello,  from "
},
```

### Discussion

In this example:

- The default values for the '```name```' and '```place```' inputs would be set to empty strings (i.e., ''), since they are of type 'string', when passed to the 'hello.js' function; therefore 'greeting' will appear a follows:
  - ```"greeting": "Hello, from "```

### Source code
The manifest file for this example can be found here:
- [manifest_hello_world_typed_parms.yaml](examples/manifest_hello_world_typed_parms.yaml)

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
    <td><a href="wskdeploy_action_fixed_parms.md#actions">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_action_advanced_parms.md#actions">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
