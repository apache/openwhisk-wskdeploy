## Actions with advanced parameters

This example builds on the previous “Hello world with basic input and output parameters” example with more robust input and output parameter declarations by using a multi-line format for declaration.

This example:
- shows how to declare input and output parameters on the action ‘```hello_world```’ using a multi-line grammar.
- adds the ‘```name```’ and ‘```place```’ input parameters, both of type ‘```string```’, to the ‘```hello_world``’ action each also includes an associated ‘```description```’ value.
- adds the ‘```greeting```’ output parameter of explicit ‘```type```’ of ‘```string```’ to the ‘```hello_world```’ action with a ‘```description```’.

### Manifest File

If we want to do more than declare the type (i.e., ‘string’, ‘integer’, ‘float’, etc.) of the input parameter, we can use the multi-line grammar.

#### Example 3: input and output parameters with explicit types and descriptions
```yaml
package:
  name: hello_world_package
  ... # Package keys omitted for brevity
  actions:
    hello_world_3:
      function: src/hello/hello.js
      runtime: nodejs@6
      inputs:
        name:
          type: string
          description: name of person
        place:
          type: string
          description: location of person
      outputs:
        greeting:
          type: string
          description: greeting string
```

### Discussion
Describing the input and output parameter types, descriptions and other meta-data helps
- tooling validate values users may input and prompt for missing values using the descriptions provided.
- allows verification that outputs of an Action are compatible with the expected inputs of another Action so that they can be composed in a sequence.

### Source code
The manifest file for this example can be found here:
- [manifest_hello_world_3.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/manifest_hello_world_2.yaml)

### Specification
For convenience, links to the schema and grammar for Actions and Parameters:
- **Actions**: [https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_actions.md#actions](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_actions.md#actions)
- **Parameters** [https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_parameters.md#parameters](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_parameters.md#parameters)

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="wskdeploy_helloworld_basic_parms.md#actions-with-basic-parameters">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Index</a></td>
    <td><a href="">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
