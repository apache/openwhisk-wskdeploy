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

# Introduction

Apache OpenWhisk™ is an open source, distributed Serverless computing
project.

Specifically, it is able to execute application logic (*Actions*) in
response to events (*Triggers*) from external sources (*Feeds*) governed
by simple conditional logic (*Rules*) around the event data.

It provides a programming model for registering and managing *Actions*,
*Triggers* and *Rules* supported by a REST-based Command Line Interface
(CLI) along with tooling to support packaging and catalog services.

The project includes a catalog of built-in system and utility *Actions*
and *Feeds*, along with a robust set of samples that demonstrate how to
integrate OpenWhisk with various external service providers (e.g.,
GitHub, Slack, etc.) along with several platform and run-time Software
Development Kits (SDKs).

The code for the Actions, along with any support services implementing
*Feeds*, are packaged according to this specification to be compatible
with the OpenWhisk catalog and its tooling. It also serves as a means
for architects and developers to model OpenWhisk package Actions as part
of full, event-driven services and applications providing the necessary
information for artifact and data type validation along with package
management operations.

### Compatibility

This specification is intended to be compatible with the following
specifications:

-   *OpenWhisk API which is defined as an OpenAPI document: *

    -   <https://raw.githubusercontent.com/openwhisk/openwhisk/master/core/controller/src/main/resources/whiskswagger.json>

-   *OpenAPI Specification when defining REST APIs and parameters:*

    -   <https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md>
