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

# Apache OpenWhisk WskDeploy

# 1.1.0
  * Upgrade to Go 1.15 (#1105)
  * Added interpolation to docker option in manifest.yaml (#1100)
  * Migrate to using Go Module to manage dependencies (#1101)
  * doc: spec_apis.md  basePath, relativePath (#1096)
  * Remove misleading web-export key (#1092)
  * Fix link to spec. docs for triggers and rules (#1093)
  * Support alt namespace uuid as tenant id on swagger api creates/deletes (#1091)
  * Fix max int. comparison for 386 archs (#1090)
  * replace the last few uses of nodejs@6 with nodejs@10 (#1089)
  * Enable Web Secure token for sequences on API create (#1087)
  * Enable setting user-supplied auth tokens on API create via require-whisk-auth annotation (#1083)
  * Add the ability to specify api via swagger in manifest (#1078)
  * Updating runtimes cache to match server side (#1086)
  * Deprecate nodejs:6 and update runtime metadata (#1085)
  * Dummy token issue (#1081)
  * Support alt. namespace resource uuid as tenant id on undeploy (#1080)
  * Support alt. namespace resource uuid as tenant id to API gatway service (#1076)
  * gradle version 5 (#1075)

## 1.0.0
  * Auto supply a dummy API GW token (#1068)
  * Support Description field in corresponding entities (#1057)
  * Honor WSK_CONFIG_FILE if variable is set (#1054)
  * Update version of whisk modules (#1067)
  * Update openwhisk-client-go version (#1059); Fixes #1050.
  * Fixes export project with feed bug (#1052)
  * Added date and commit SHA to the version command (#1047)
  * Account for variability in Account Limits from various provider impls. (#1048)

## 0.10.0-incubating
  * Fix language:default runtime setting (#1039)
  * Link api schema to specification overview page (#1030)
  * Add API entity schema (#1029)
  * Only print info messages when the verbose flag is enabled (#1027)
  * Enable programatic support for additional request headers (#1023)
  * Add swift 4.2 (#1022)
  * Add PHP 7.3 runtime (#1021)
  * Add support for .Net dotnet:2.2 action kind (#1019)
  * Add fallback method to find wskprops when go-client fails (#1015)
  * Add ruby to specifications file (#1012)
  * Add nodejs:10 kind for wskdeploy (#1011)
  * Add support to parse the type slice (#1010)
  * Add pkg and action version number with interpolation (#1009)
  * Enable export verbose output (#996)
  * Add go-runtime (#1006)
  * Upgrade the Go version to 1.9 (#997)
  * Introducing include and exclude in zip action (#991)
  * Add ruby runtime (#983)
  * Skipping response data in case of http request was successful (#981)
  * Bug fix. Export shouldn't fail when ApiGW missing (#979)
  * Added HTTP response documentation (#976)
  * Fixed apigateway docs and example manifests (#974)

## 0.9.8-incubating
  Initial Apache Release
