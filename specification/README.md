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

# OpenWhisk Packaging Specification

[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

## Purpose

In general, the goal of this specification is to evolve a simple grammar to describe and deploy a complete OpenWhisk package without having to use any APIs.  Primarily, it is accomplished by coding a Package Manifest file and optionally a Deployment file using YAML.

# Programming Guide

If you want to learn how to write Packages and Applications by example using the specification and deploy them using the ```wskdeploy``` utility, please read the step-by-step guide:
- "[wskdeploy utility by example](../docs/programming_guide.md#wskdeploy-utility-by-example)"

# Package Specification

Portions of the OpenWhisk Packaging Specification, for convenience, are made available here in Markdown/HTML format. The canonical source for the specification is in PDF format and can be found within the [archive](archive) directory.

- Current version (link): [openwhisk_v0.9.2.pdf](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/archive/openwhisk_v0.9.2.pdf)

#### Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC
2119](http://www.ietf.org/rfc/rfc2119.txt).

The OpenWhisk packaging specification is licensed under [The Apache License,
Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).

### Index

#### Overview
- [Introduction](html/spec_intro.md#introduction) - an overview of the goals for the packaging specification.
    - [Compatibility](html/spec_intro.md#compatibility) - describes intent to be compatible with other standards.
    - [Revision History](html/spec_history.md#revision-history) - lists changes to specification by version/revision.
- [Programming Model](html/spec_programming_model.md#programming-model) - an overview of the OpenWhisk programming model.
- [Package Processing](html/spec_package_processing.md#package-processing) - an overview of the OpenWhisk programming model.

#### Schema
- [Parameters](html/spec_parameters.md#parameters) - grammar, schema and examples for input and output parameters.
- [Parameter Types](html/spec_types.md#parameter-types) - supported YAML and OpenWhisk Types.
- [Shared Entity Schema](html/spec_shared_entity_schema.md#shared-entity-schema) - fields that are common among entities in the programming model.
- [Packages](html/spec_packages.md#packages) - grammar, schema and examples for Packages.
- [Actions](html/spec_actions.md#actions) - grammar, schema and examples for Actions.
- [Triggers and Rules](html/spec_trigger_rule.md#triggers-and-rules) - grammar, schema and examples for Triggers and Rules.
- [Sequences](html/spec_sequences.md#sequences) - shema to compose multiple Actions into a sequence.
- [Normative & Non-normative References](html/spec_normative_refs.md)

---

### Note
This specification is under development and in draft status; therefore, it is subject to change during this development period.  We are posting drafts seeking review, comments, suggestions and use cases from the OpenWhisk and greater Serverless community.
