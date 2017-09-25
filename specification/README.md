# OpenWhisk Packaging Specification

[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

Current version (link):  [openwhisk_v0.8.9.pdf](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/openwhisk_v0.8.9.pdf)

## Purpose

In general, the goal of this specification is to evolve a simple grammar to describe and deploy a complete OpenWhisk package without having to use any APIs.  Primarily, it is accomplished by coding a Package Manifest file and optionally a Deployment file using YAML.

# Programming Guide

If you want to learn how to write Packages and Applications by example using the specification and deploy them using the ```wskdeploy``` utility, please read the step-by-step guide:
- "[wskdeploy utility by example](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/docs/programming_guide.md#wskdeploy-utility-by-example)"

# OpenWhisk Package Specification (HTML)

Portions of the OpenWhisk Packaging Specification, for convenience, are made available here in HTML format. The canonical source for the specification is in PDF format and can be found within this directory.

## Index

- [Parameters](spec_parameters.md#parameters) - grammar, schema and examples for input and output parameters.
- [Parameter Types](spec_types.md#parameter-types) - supported YAML and OpenWhisk Types.
- [Shared Entity Schema](spec_shared_entity_schema.md#shared-entity-schema) - fields that are common among entities in the programming model.
- [Packages](spec_packages.md#packages) - grammar, schema and examples for Packages.
- [Actions](spec_actions.md#actions) - grammar, schema and examples for Actions.
- [Triggers and Rules](spec_trigger_rule.md#triggers-and-rules) - grammar, schema and examples for Triggers and Rules.
- [Sequences](spec_sequences.md#sequences) - shema to compose multiple Actions into a sequence.
- [Normative & Non-normative References](spec_normative_refs.md)

---

### Note
This specification is under development and in draft status; therefore, it is subject to change during this development period.  We are posting drafts seeking review, comments, suggestions and use cases from the OpenWhisk and greater Serverless community.

---

## Package Manifest and Deployment files

- **Package Manifest file**: Contains the Package definition along with any included Action, Trigger or Rule definitions that comprise the package.  This file includes the schema of input and output data to each entity for validation purposes.
- **Deployment file**: Contains the values and bindings used configure a Package to a target OpenWhisk platform provider’s environment and supply input parameter values for Packages, Actions and Triggers.  This can include Namespace bindings, security and policy information.

### Conceptual Manifest and Deployment file usage

The following images outline the basic process for creating and using both Manifest and Deployment files against a typical developer workstream:

#### Conceptual Manifest file creation
![Manifest file creation](images/OpenWhisk%20-%20Conceptual%20Manifest%20File%20Creation.png "image 1")

---

### Conceptual Manifest file deployment
![Manifest file deployment](images/OpenWhisk%20-%20Conceptual%20Manifest%20File%20Deployment.png "image 1")
