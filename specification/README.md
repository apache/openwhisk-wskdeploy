# OpenWhisk Packaging Specification v0.8

[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

## Purpose

This specification defines two file artifacts, along with YAML schema, that are used to describe OpenWhisk Packages for deployment to a target OpenWhisk platform; these include:

*	**Package Manifest file**: Contains the Package definition along with any included Action, Trigger or Rule definitions that comprise the package.  This file includes the schema of input and output data to each entity for validation purposes.
*	**Deployment file**: Contains the values and bindings used configure a Package to a target OpenWhisk platform provider’s environment and supply input parameter values for Packages, Actions and Triggers.  This can include Namespace bindings, security and policy information.

## Note
This specification is under development and in draft status; therefore is subject to change during this time.  We are seeking input from the OpenWhisk community to provide us to review the contents, suggest edits, provide interesting use cases, etc.  In general, make it a top-quality means to describe a complete OpenWhisk package without having to understand and API.  

## Formats

Once the draft progresses further (i.e., known outstanding "high priority" design issues have been addressed), we will make the document available in markdown.  At this time the following formats are provided for review:

* PDF: [openwhisk_v0.8.pdf](openwhisk_v0.8.pdf)
