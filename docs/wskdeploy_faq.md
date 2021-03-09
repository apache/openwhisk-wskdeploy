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

# ```wskdeploy``` utility FAQ

### What if ```wskdeploy``` finds an error in my manifest?

- The ```wskdeploy``` utility will not attempt to deploy a package if an error in the manifest is detected, but will report as much information as it can to help you locate the error in the YAML file.

### What if ```wskdeploy``` encounters an error during deployment?

-  The ```wskdeploy``` utility will cease deploying as soon as it receives an error from the target platform and display what error information it receives to you.
- then it will attempt to undeploy any entities that it attempted to deploy.

### What is the order of precedence for OpenWhisk credentials?

- The ```wskdeploy``` utility finds the credentials (apihost, namespace, and auth) as well as the APIGW_ACCESS_TOKEN in the following precedence from highest to lowest:
  - ```wskdeploy``` command line (i.e. ```wskdeploy --apihost --namespace --auth```)
  - The deployment file
  - The manifest file
  - The .wskprops file
- So when filling out the credentials and the APIGW_ACCESS_TOKEN, it first looks to the command line. Anything not found on the command line is attempted to be filled by the deployment file. Next it searches for all the unfilled values in the manifest file, and finally any unfilled values are looked for in the .wskprops file. After failing to find a value in the places mentioned above, namespace will default to "_". In the case of a TravisCI environment, APIGW_ACCESS_TOKEN will be set to "DUMMY TOKEN" if it is not already set via one of the mechanisms described previously.
