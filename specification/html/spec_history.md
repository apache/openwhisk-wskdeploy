<!--
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
-->

## Revision History

| Version | Date | Notes |
|:---|:---|:---|
| 0.8.1 | 2016-11-03 | </br><ul><li>Initial public point draft</ul> |
| 0.8.2 | 2016-12-12 | </br><ul><li>Begin adding use cases, examples to all sections.</ul> |
| 0.8.3 | 2017-02-02 | </br><ul><.i>Add more use cases, examples<li> Add description of "$" (dollar) notation</ul> |
| 0.8.4 | 2017-04-18 | </br><li>Support JSON parameter type;</br>Clarify use of Parameter single-line grammar and inferred types.<li>Add support for API Gateway mappings.<li>Add support for Web Actions</ul> |
| 0.8.5 | 2017-04-21 | </br><ul><li>Add support for “dependencies”, that is allow automatic deployment of other OpenWhisk packages (from GitHub) that the current package declares as a dependency.</ul> |
| 0.8.6 | 2017-07-25 | </br><ul><li>Clarified requirements for \$ dollar notation.<li>Updated conceptual Manifest/Deployment File processing images. |
| 0.8.7 | 2017-08-24 | </br><ul><li>Added explicit Application entity and grammar.<li>Added API listing to Package entity.<li>Cleaned up pseudo-grammar which contained various uses of credentials in places not intended.<li>Fixed Polygon Tracking example (indentation incorrect).</ul> |
| 0.8.8 | 2017-08-29 | </br><ul><li>Created a simplified API entity (i.e., “api”) grammar that allows multiple sets of named APIs for the same basepath.<li>Acknowledge PHP as supported runtime (kind).<li>Added “sequences” entity as a convenient way to declare action sequences in the manifest. Updated supported runtime values.</ul> |
| 0.8.9 | 2017-09-22 | Clarified “version” key requirements for Package (required) and Action (optional); removed from shared entity schema.</ul> |
| 0.8.9.1 | 2017-09-29 | </br><ul><li>Made “license” key optional for package.<li>keyword “package” (singular) and “packages” (plural) both allowed.<li>Adjusted use case examples to reflect these changes.<li>Rework of schema use cases into full, step-by-step examples.<li>Spellcheck, fixed bugs, update examples to match web-based version.</ul> |
| 0.8.9.1 | 2017-10-06 | </br><ul><li>Added grammar and example for concatenating string values on input parameters using environment variables.</ul> |
| 0.9.0 | 2017-11-23 | </br><ul><li>Identified new user scenarios including: clean, refresh, sync, pre/post processing</ul> |
| 0.9.1 | 2017-11-30 | </br><ul><li>Clarified “runtime” field on Action is equivalent to “kind” parameter used on the Apache OpenWhisk CLI for Actions.<li>Added “project” key as an synonym name for “application”.” key, moving application to become deprecated. Project name made required.<li>Support “public” (i.e., publish) key on Package.<li>Documented support for the “raw-http” annotation under Action.<li>Documented support for the “final” annotation under Action.<li>Documented support for the “main” field under Action.<li>Dollar Notation section becomes Interpolation / updates<li><ul><li>Supported beyond Parameter values<li>Package names can be interpolated<li>Annotations values can be interpolated<li>Multiple replacements supported in same value</ul><li>Usage scenarios 6-8 added, i.e., Clean, Project Sync, Tool chain support. |
| 0.9.2 | 2018-02-05 | </br><ul><li>Fixed and clarified the allowed values for “web-export” on Action.<li>Clarified use of “final” on Action.<li>Added support (planned) for “web-custom-options” and “require-whisk-auth. ” flags on Actions (annotations)<li>Deprecate ‘application’ and ‘package’ completely (no longer supported)</ul> |
| 0.9.2.2   2018-04-04 | </br><ul><li>Allow “web” key as an overload for “web-export” key for to indicate Web Actions.<li>Added Web Sequences, specify a sequence is a Web Action.<li>Added support for Conductor Actions, to align with OpenWhisk CLI support.<li>Added “docker” and “native” binary support under Action.<li>Added in-line “code” support under Action.<li>Support \$\$, double-dollar notation for string literals on parameter values.<li>Added support for “default” package (allowing all entities to be assigned directly under the user’s default namespace), that is not requiring a package name to be created.</ul> |


<!--
 Bottom Navigation
-->
---
<html>
<div align="center">
<a href="../README.md#index">Index</a>
</div>
</html>
