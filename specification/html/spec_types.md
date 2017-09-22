
<!--
********************************
  YAML Types
********************************
-->

<h2>YAML Types</H2>

<p>Many of the types we use in this profile are <i>built-in</i> types from the <a
href="http://www.yaml.org/spec/1.2/spec.html">YAML 1.2 specification</a>
(i.e., those identified by the 'tag:yaml.org,2002' version tag). </p>

<p>The following table declares the valid YAML type URIs and
aliases that SHALL be used when defining parameters or properties within an
OpenWhisk package manifest:<a> </a></p>

<table width="100%">
 <tr>
  <th>
  <p>Type Name</p>
  </th>
  <th>
  <p>Type URI</p>
  </th>
  <th>
  <p>Notes</p>
  </th>
 </tr>
 <tr>
  <td>
  <p><a>string</a></p>
  </td>
  <td>
  <p>tag:yaml.org,2002:str (default)</p>
  </td>
  <td>
  <p>Default type if no type provided</p>
  </td>
 </tr>
 <tr>
  <td>
  <p><a>integer</a></p>
  </td>
  <td>
  <p>tag:yaml.org,2002:int</p>
  </td>
  <td>
  <p>Signed. Includes large integers (i.e., long type)</p>
  </td>
 </tr>
 <tr>
  <td>
  <p><a>float</a></p>
  </td>
  <td>
  <p>tag:yaml.org,2002:float</p>
  </td>
  <td>
  <p>Signed. Includes large floating point values (i.e., double type)</p>
  </td>
 </tr>
 <tr>
  <td>
  <p><a>boolean</a></p>
  </td>
  <td>
  <p>tag:yaml.org,2002:bool</p>
  </td>
  <td>
  <p>This specification uses lowercase 'true' and lowercase 'false'</p>
  </td>
 </tr>
 <tr>
  <td>
  <p><a>timestamp</a></p>
  </td>
  <td>
  <p>tag:yaml.org,2002:timestamp (see <a href="spec_normative_refs.md#normative-references">YAML-TS-1.1</a>)</p>
  </td>
  <td>
  <p>ISO 8601 compatible.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p><a>null</a></p>
  </td>
  <td>
  <p>tag:yaml.org,2002:null</p>
  </td>
  <td>
  <p>Different meaning than an empty string, map, list, etc.</p>
  </td>
 </tr>
</table>

<h4>Requirements</h4>
<ul>
<li>The 'string' type SHALL be the default type when not specified on a parameter or property declaration.</li>
<li>All 'boolean' values SHALL be lowercased (i.e., 'true' or 'false').</li>
</ul>
<html>

<!--
********************************
  OpenWhisk Types
********************************
-->
<html>
<h2>OpenWhisk Types</h2>
<p>In addition to the YAML built-in types, OpenWhisk supports the types listed in the table below. A complete description of each of these types is provided below.</p>
<p>&nbsp;</p>
<table width="100%">
 <tr>
  <th width=20%>Type Name</th>
  <th width=30%>Description</th>
  <th width=50%>Notes</th>
 </tr>
 <tr>
  <td>
  <p>version</p>
  </td>
  <td>
  <p>tag:maven.apache.org:version (see <a href="spec_normative_refs.md#normative-references">Maven version</a>)</p>
  </td>
  <td>
  <p>Typically found in modern tooling (i.e., 'package@version' or 'package:version' format).</p>
  </td>
 </tr>
 <tr>
  <td>
  <p><a>string256</a></p>
  </td>
  <td>
  <p>long length strings (e.g., descriptions)</p>
  </td>
  <td>
  <p>A string type limited to 256 characters.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p><a>string64</a></p>
  </td>
  <td>
  <p>medium length strings (e.g., abstracts, hover text)</p>
  </td>
  <td>
  <p>A string type limited to 64 characters.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p><a>string16</a></p>
  </td>
  <td>
  <p>short length strings (e.g., small form-factor list displays)</p>
  </td>
  <td>
  <p>A string type limited to 16 characters.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>json</p>
  </td>
  <td>
  <p>The parameter value represents a JavaScript Object Notation (JSON) data object.</p>
  </td>
  <td>
  <p>The deploy tool will validate the corresponding parameter value against JSON schema.</p>
  <p>Note: The implied schema for JSON the JSON Schema (see <a href="http://json-schema.org/)">http://json-schema.org/</a>).</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>scalar-unit</p>
  </td>
  <td>
  <p>Convenience type for declaring common scalars that have an associated unit. For example, '10 msec.', '2 Gb', etc.)</p>
  </td>
  <td>
  <p>Currently, the following scalar-unit subtypes are supported:</p>
  <ul>
  <li>scalar-unit.size</li>
  <li>scalar-unit.time</li>
  </ul>
  <p>See description below for details.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>schema</p>
  </td>
  <td>
  <p>The parameter itself is an OpenAPI Specifcation v2.0 <b>Schema
  Object</b> (in YAML formatt) with self-defining schema.</p>
  </td>
  <td>
  <p>The schema declaration follows the <a href="#REF_SWAGGER_2_0">OpenAPI</a> v2.0 specification for Schema Objects (YAML format).</p>
  <p>Specifically, see https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#schemaObject</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>object</p>
  </td>
  <td>
  <p>The parameter itself is an object with the associated
  defined Parameters (schemas). </p>
  </td>
  <td>
  <p>Parameters of this type would include a declaration of its
  constituting Parameter schema.</p>
  </td>
 </tr>
</table>

<h2>scalar-unit types</h2>
<p>Scalar-unit types can be used to define scalar values along
with a unit from the list of recognized units (a subset of GNU units) provided
below.</p>
</html>

### Grammar
```
<scalar> <unit>
```

<html>
<h3>Requirements</h3>
<ul>
<li><b>Whitespace</b>: any number of spaces (including zero or none) SHALL
be allowed between the scalar value and the unit value.</p>
<p> It SHALL be considered an error if either the scalar or unit
portion is missing on a property or attribute declaration derived from any
scalar-unit type.</p>

<h3>Recognized units for sizes (i.e., scalar-unit.size)</h3>
<table width="100%">

  <tr>
   <th>
   <p>Unit</p>
   </th>
   <th>
   <p>Description</p>
   </th>
  </tr>

 <tr>
  <td>
  <p>B</p>
  </td>
  <td>
  <p>byte</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>kB</p>
  </td>
  <td>
  <p>kilobyte (1000 bytes)</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>MB</p>
  </td>
  <td>
  <p>megabyte (1000000 bytes)</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>GB</p>
  </td>
  <td>
  <p>gigabyte (1000000000 bytes)</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>TB</p>
  </td>
  <td>
  <p>terabyte (1000000000000 bytes)</p>
  </td>
 </tr>
</table>

<h3>Example</h3>

</html>

```yaml
inputs:
  memory_size:
    type: scalar-unit.size
    value: 256 MB
```

<html>
<h3>Recognized units for times (i.e., scalar-unit.time)</h3>
<table>
  <tr>
   <th width="20%">
   <p>Unit</p>
   </th>
   <th width="80%">
   <p>Description</p>
   </th>
  </tr>

 <tr>
  <td>
  <p>d</p>
  </td>
  <td>
  <p>days</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>h</p>
  </td>
  <td>
  <p>hours</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>m</p>
  </td>
  <td>
  <p>minutes</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>s</p>
  </td>
  <td>
  <p>seconds</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>ms</p>
  </td>
  <td>
  <p>milliseconds</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>us</p>
  </td>
  <td>
  <p>microseconds</p>
  </td>
 </tr>
</table>

<h3>Example</h3>
</html>
```yaml
inputs:
  max_execution_time:
    type: scalar-unit.time
    value: 600 s
```

<html>
<h3>Object type example</h3>

<p>The Object type allows for complex objects to be declared as parameters with an optional validatable schema.</p>
</html>

```yaml
inputs:
  person:
    type: object
    parameters: <schema>
```
