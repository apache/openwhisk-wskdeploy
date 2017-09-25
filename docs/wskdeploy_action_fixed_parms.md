## Actions with fixed input parameters

This example builds upon the previous “hello world” example and shows how fixed values can be supplied to the input parameters of an Action.

It shows how to:
- declare input parameters on the action ‘```hello_world```’ using a single-line grammar.
- add the ‘```name```’ and ‘```place```’ input parameters, with fixed values “```Sam```” and “```the Shire```” respectively.

### Manifest File

#### Example 2: 'Hello world' with explicit input and output parameter declarations
```yaml
package:
  name: hello_world_package
  version: 1.0
  license: Apache-2.0
  actions:
    hello_world_2:
      function: src/hello.js
      inputs:
        name: Sam
        place: the Shire
```
### Deploying
```sh
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world_2 --blocking
```

### Discussion
This packaging specification grammar places an emphasis on simplicity for the casual developer who may wish to hand-code a Manifest File; however, it also provides a robust optional schema that can be advantaged when integrating with larger application projects using design and development tooling such as IDEs.

In this example:

- The default values for the '```name```' and '```place```' inputs would be set to empty strings (i.e., ''), since they are of type 'string', when passed to the 'hello.js' function; therefore 'greeting' will appear a follows:
  - ```"greeting": "Hello, from "```

### Source code
The manifest file for this example can be found here:
- [manifest_hello_world_2.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/examples/manifest_hello_world_2.yaml)

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
    <td><a href="">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
