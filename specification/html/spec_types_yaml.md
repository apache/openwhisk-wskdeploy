<h2>YAML Types</H2>
<p>Many of the types we use in this profile are <i>built-in</i> types from the <a
href="http://www.yaml.org/spec/1.2/spec.html">YAML 1.2 specification</a>
(i.e., those identified by the Òtag:yaml.org,2002Ó version tag). </p>
<p>&nbsp;</p>
<p>The following table declares the valid YAML type URIs and
aliases that SHALL be used when defining parameters or properties within an
OpenWhisk package manifest:<a> </a></p>
<p>&nbsp;</p>
<table width="98%">
 <tr>
  <td width=65>
  <p>Type Name</p>
  </td>

  <td width=154>
  <p>Type URI</p>
  </td>

  <td width=241>
  <p>Notes</p>
  </td>

 </tr>
 <tr>
  <td width=65>
  <p><a>string</a></p>
  </td>

  <td width=154>
  <p>tag:yaml.org,2002:str (default)</p>
  </td>

  <td width=241>
  <p>Default type if no type provided</p>
  </td>

 </tr>
 <tr>
  <td width=65>
  <p><a>integer</a></p>
  </td>

  <td width=154>
  <p>tag:yaml.org,2002:int</p>
  </td>

  <td width=241>
  <p>Signed. Includes large integers (i.e., long type)</p>
  </td>

 </tr>
 <tr>
  <td width=65>
  <p><a>float</a></p>
  </td>

  <td width=154>
  <p>tag:yaml.org,2002:float</p>
  </td>

  <td width=241>
  <p>Signed. Includes large floating point values (i.e., double type)</p>
  </td>

 </tr>
 <tr>
  <td width=65>
  <p><a>boolean</a></p>
  </td>

  <td width=154>
  <p>tag:yaml.org,2002:bool</p>
  </td>

  <td width=241>
  <p>This specification uses lowercase 'true' and lowercase 'false'</p>
  </td>

 </tr>
 <tr>
  <td width=65>
  <p><a>timestamp</a></p>
  </td>

  <td width=154>
  <p>tag:yaml.org,2002:timestamp (see <a href="#REF_YAML_TIMESTAMP_1_1">YAML-TS-1.1</a>)</p>
  </td>

  <td width=241>
  <p>ISO 8601 compatible.</p>
  </td>

 </tr>
 <tr>
  <td width=65>
  <p><a>null</a></p>
  </td>

  <td width=154>
  <p>tag:yaml.org,2002:null</p>
  </td>

  <td width=241>
  <p>Different meaning than an empty string, map, list, etc.</p>
  </td>

 </tr>
 <tr>
  <td width=65>
  <p><a></a><a href="#REF_MAVEN_VERSION">version</a></p>
  </td>

  <td width=154>
  <p>tag:maven.apache.org:version (see <a href="#REF_MAVEN_VERSION">Maven version</a>)</p>
  </td>

  <td width=241>
  <p>Typically found in modern tooling (i.e., 'package@version' or 'package:version' format).</p>
  </td>

 </tr>
 <tr>
  <td width=65>
  <p><a
 >string256</a></p>
  </td>

  <td width=154>
  <p>long
  &nbsp;length strings (e.g.,
  descriptions)</p>
  </td>

  <td width=241>
  <p>A
  string type limited to 256 characters.</p>
  </td>

 </tr>
 <tr>
  <td width=65>
  <p><a
 >string64</a></p>
  </td>

  <td width=154>
  <p>medium
  length strings (e.g., abstracts, hover text)</p>
  </td>

  <td width=241>
  <p>A
  string type limited to 64 characters.</p>
  </td>

 </tr>
 <tr>
  <td width=65>
  <p><a
 >string16</a></p>
  </td>

  <td width=154>
  <p>short
  length strings (e.g., small form-factor list displays)</p>
  </td>

  <td width=241>
  <p>A
  string type limited to 16 characters.</p>
  </td>

 </tr>
</table>
<h4>Requirements</h4>
<ul>
<li>The 'string' type SHALL be the default type when not specified on a parameter or property declaration.</li>
<li>All 'boolean' values SHALL be lowercased (i.e., 'true' or 'false').</li>
</ul>
