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
| 0.8.1 | 2016-11-03 | Initial public point draft, Working Draft 01 |
| 0.8.2 | 2016-12-12 | Working Draft 02, Add. Use cases, examples |
| 0.8.3 | 2017-02-02 | Working Draft 03, Add use cases, examples, "$" (dollar) notation |
| 0.8.4 | 2017-04-18 | Working Draft 04</br>
Support JSON parameter type;</br>
Clarify use of Parameter single-line grammar and inferred types.</br>Add support for API Gateway mappings.</br>Add support for Web Actions |
| 0.8.5 | 2017-04-21 | Add support for “dependencies”, that is allow automatic deployment of other OpenWhisk packages (from GitHub) that the current package declares as a dependency. |
| 0.8.6 | 2017-07-25 | Clarified requirements for \$ dollar notation. |
                      -   Updated conceptual Manifest/Deployment File processing images.
| 0.8.7 | 2017-08-24 | Added explicit Application entity and grammar.
                      -   Added API listing to Package entity.
                      -   Cleaned up pseudo-grammar which contained various uses of credentials in places not intended.
                      -   Fixed Polygon Tracking example (indentation incorrect). |
| 0.8.8 | 2017-08-29 | Created a simplified API entity (i.e., “api”) grammar that allows multiple sets of named APIs for the same basepath.-   Acknowledge PHP as supported runtime (kind).
                      -   Added “sequences” entity as a convenient way to declare action sequences in the manifest. Updated supported runtime values. |
| 0.8.9 | 2017-09-22 | Clarified “version” key requirements for Package (required) and Action (optional); removed from shared entity schema. |
| 0.8.9.1 | 2017-09-29 | Made “license” key optional for package.
                      -   keyword “package” (singular) and “packages” (plural) both allowed.
                      -   Adjusted use case examples to reflect these changes.
                      -   Rework of schema use cases into full, step-by-step examples.
                      -   Spellcheck, fixed bugs, update examples to match web-based version. |
| 0.8.9.1 | 2017-10-06 | Added grammar and example for concatenating string values on input parameters using environment variables. |
| 0.9.0 | 2017-11-23 | Identified new user scenarios including: clean, refresh, sync, pre/post processing |
| 0.9.1 | 2017-11-30 | Clarified “runtime” field on Action is equivalent to “kind” parameter used on the Apache OpenWhisk CLI for Actions.
                      -   Added “project” key as an synonym name for “application”.” key, moving application to become deprecated. Project name made required.
                      -   Support “public” (i.e., publish) key on Package.
                      -   Documented support for the “raw-http” annotation under Action.
                      -   Documented support for the “final” annotation under Action.
                      -   Documented support for the “main” field under Action.
                      -   Dollar Notation section becomes Interpolation / updates
                          -   Supported beyond Parameter values
                          -   Package names can be interpolated
                          -   Annotations values can be interpolated
                          -   Multiple replacements supported in same value
                    -   Usage scenarios 6-8 added, i.e., Clean, Project Sync, Tool chain support. |
| 0.9.2 | 2018-02-05 | Fixed and clarified the allowed values for “web-export” on Action.
                      -   Clarified use of “final” on Action.
                      -   Added support (planned) for “web-custom-options” and “require-whisk-auth. ” flags on Actions (annotations)
                      -   Deprecate ‘application’ and ‘package’ completely (no longer supported) |
| 0.9.2.2   2018-04-04 | Allow “web” key as an overload for “web-export” key for to indicate Web Actions.
                      -   Added Web Sequences, specify a sequence is a Web Action.
                      -   Added support for Conductor Actions, to align with OpenWhisk CLI support.
                      -   Added “docker” and “native” binary support under Action.
                      -   Added in-line “code” support under Action.
                      -   Support \$\$, double-dollar notation for string literals on parameter values.
                      -   Added support for “default” package (allowing all entities to be assigned directly under the user’s default namespace), that is not requiring a package name to be created. |


<!--
 Bottom Navigation
-->
---
<html>
<div align="center">
<a href="../README.md#index">Index</a>
</div>
</html>
