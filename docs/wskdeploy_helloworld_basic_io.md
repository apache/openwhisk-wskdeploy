## "Hello world" Action with input and output parameters

This use case extends the 'Hello world' example with explicit input and output Parameter declarations.

This example:
- shows how to declare input and output parameters on the action 'hello_world'
using a simple, single-line grammar.
- adds two input parameters, '```name```' and '```place```', both of type '```string```' to the '```hello_world```' action.
- adds one output parameter, '```greeting```' of type string to the '```hello_world```' action.

### Manifest File

#### Example 2: 'Hello world' with explicit input and output parameter declarations
```yaml
package:
  name: hello_world_package
  ... # Package keys omitted for brevity
  actions:
    hello_world_2:
      function: src/hello.js
      inputs:
        name: string
        place: string
      outputs:
        greeting: string
```

### Discussion
This packaging specification grammar places an emphasis on simplicity for the casual developer who may wish to hand-code a Manifest File; however, it also provides a robust optional schema that can be advantaged when integrating with larger application projects using design and development tooling such as IDEs.

In this example:

- The default values for the '```name```' and '```place```' inputs would be set to empty strings (i.e., ''), since they are of type 'string', when passed to the 'hello.js' function; therefore 'greeting' will appear a follows:
  - ```"greeting": "Hello, from "```

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
