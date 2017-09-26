# Actions

## Advanced parameters

This example builds on the previous [“Hello world" with typed input and output parameters](wskdeploy_action_typed_parms.md#actions) example with more robust input and output parameter declarations by using a multi-line format for declaration.

This example:
- shows how to declare input and output parameters on the action ‘```hello_world```’ using a multi-line grammar.
- adds the ‘```name```’ and ‘```place```’ input parameters, both of type ‘```string```’, to the ‘```hello_world``’ action each also includes an associated ‘```description```’ value.
- adds the ‘```greeting```’ output parameter of explicit ‘```type```’ of ‘```string```’ to the ‘```hello_world```’ action with a ‘```description```’.

### Manifest File

If we want to do more than declare the type (i.e., ‘string’, ‘integer’, ‘float’, etc.) of the input parameter, we can use the multi-line grammar.

#### _Example: input and output parameters with explicit types and descriptions_
```yaml
package:
  name: hello_world_package
  ... # Package keys omitted for brevity
  actions:
    hello_world_advanced_parms:
      function: src/hello/hello_plus.js
      runtime: nodejs@6
      inputs:
        name:
          type: string
          description: name of person
          value: Sam
        place:
          type: string
          description: location of person
          default: unknown
        children:
          type: integer
          description: Number of children
          default: 0
        height:
          type: float
          description: height in meters
          default: 0.0
      outputs:
        greeting:
          type: string
          description: greeting string
```

### Deploying
```sh
$ wskdeploy -m docs/examples/manifest_hello_world_advanced_parms.yaml
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world_advanced_parms --blocking
```

### Result
The invocation should return an 'ok' with a response that includes this result:
```sh
"result": {
    "greeting": "Hello, Sam from unknown"
},
```

### Discussion
Describing the input and output parameter types, descriptions, defaults and other data:
- enables tooling to validate values users may input and prompt for missing values using the descriptions provided.
- allows verification that outputs of an Action are compatible with the expected inputs of another Action so that they can be composed in a sequence.

### Source code
The manifest file for this example can be found here:
- [manifest_hello_world_advanced_parms.yaml](examples/manifest_hello_world_advanced_parms.yaml)

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
    <td><a href="wskdeploy_action_typed_parms.md#actions">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_triggerrule_basic.md#triggers-and-rules">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
