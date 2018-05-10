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

# Configuring ```wskdeploy```

At minimum, the wskdeploy utility needs valid OpenWhisk APIHOST and AUTH variable to attempt deployment. In this case the default target namespace is assumed; otherwise, NAMESPACE also needs to be provided.

## Precedence order

Wskdeploy attempts to find these values in the following order:

1. **Command line**

Values supplied on the command line using the ```apihost```, ```auth``` and ```namespace``` flags will override values supplied elsewhere (below).

for example the following flags can be used:

```
$ wskdeploy --apihost <host> --auth <auth> --namespace <namespace>
```

Command line values are considered higher in precedence since they indicate an intentional user action to override values that may have been supplied in files.

2. **Deployment file**

Values supplied in a Deployment YAML file will override values found elsewhere (below).

3. **Manifest file**

Values supplied in a Manifest YAML file will override values found elsewhere (below).

4. **.wskprops**

Values set using the Whisk Command Line Interface (CLI) are stored in a ```.wskprops```, typically in your $HOME directory, will override values found elsewhere (below).

It assumes that you have setup and can run the wskdeploy as described in the project README. If so, then the utility will use the OpenWhisk APIHOST and AUTH variable values in your .wskprops file to attempt deployment.

