<!--
#
# Licensed to the Apache Software Foundation (ASF) under one or more contributor
# license agreements.  See the NOTICE file distributed with this work for additional
# information regarding copyright ownership.  The ASF licenses this file to you
# under the Apache License, Version 2.0 (the # "License"); you may not use this
# file except in compliance with the License.  You may obtain a copy of the License
# at:
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed
# under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
# CONDITIONS OF ANY KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations under the License.
#
-->

# ```wskdeploy``` utility FAQ

### What if ```wskdeploy``` finds an error in my manifest?

- The ```wskdeploy``` utility will not attempt to deploy a package if an error in the manifest is detected, but will report as much information as it can to help you locate the error in the YAML file.

### What if ```wskdeploy``` encounters an error during deployment?

-  The ```wskdeploy``` utility will cease deploying as soon as it receives an error from the target platform and display what error information it receives to you.
- then it will attempt to undeploy any entities that it attempted to deploy.
  - If "interactive mode" was used to deploy, then you will be prompted to confirm you wish to undeploy.
