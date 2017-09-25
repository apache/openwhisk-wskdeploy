## "Hello world" with basic input and output parameters

This use case extends the ÒHello worldÓ example with explicit input and output Parameter declarations.

This example:
- shows how to declare input and output parameters on the action Ôhello_worldÕ
using a simple, single-line grammar.</p>
- adds two input parameters, ÔnameÕ and ÔplaceÕ, both of type ÔstringÕ to the Ôhello_worldÕ action.
- adds one output parameter, ÔgreetingÕ of type string to the Ôhello_worldÕ action.</p>
Manifest File

### Example 2: ÒHello worldÓ with explicit input and output parameter declarations
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

- The default values for the ÔnameÕ and ÔplaceÕ inputs would be set to empty strings (i.e., ÒÓ), since they are of type ÔstringÕ, when passed to the Ôhello.jsÕ function.</p>
