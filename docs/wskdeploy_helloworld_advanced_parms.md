## Action with advanced parameters

If we want to do more than declare the type (i.e., ‘string’, ‘integer’, ‘float’, etc.) of the input parameter, we can use the multi-line grammar.

### Manifest File

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


### Source code
The manifest file for this example can be found here:
- [manifest_hello_world_3.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/manifest_hello_world_2.yaml)

### Specification
For convenience, the schema and grammar for declaring an **Parameters** can be found here:
[https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_parameters.md#parameters](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_parameters.md#parameters)

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
