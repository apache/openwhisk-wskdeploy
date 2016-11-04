# OpenWhisk Packaging Specification v0.8

[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

This document defines two file artifacts that are used to deploy Packages to a target OpenWhisk platform; these include:
•	Package Manifest file: Contains the Package definition along with any included Action, Trigger or Rule definitions that comprise the package.  This file includes the schema of input and output data to each entity for validation purposes.
•	Deployment file: Contains the values and bindings used configure a Package to a target OpenWhisk platform provider’s environment and supply input parameter values for Packages, Actions and Triggers.  This can include Namespace bindings, security and policy information.

Conceptual Package creation and publishing
The following diagram illustates how a developer would create OpenWhisk code artifacts and associate a Package Manifest file that describes them for deployment and reuse.
