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

# Integration Test - Dependencies

`wskdeploy` supports dependencies where it allows you to declare other OpenWhisk
packages that your application or project (manifest) is dependent on. With declaring
dependent packages, `wskdeploy` supports automatic deployment of those dependent
packages.

For example, our `root-app` application is dependent on `child-app` package located
at https://github.com/pritidesai/child-app.
We can declare this dependency in `manifest.yaml` with:

```yaml
packages:
  root-app:
    namespace: guest
    dependencies:
      child-app:
        location: github.com/pritidesai/child-app
    triggers:
      trigger1:
    rules:
      rule1:
        trigger: trigger1
        action: child-app/hello
```

**Note:**

1. Package name of the dependent package `child-app` should match GitHub repo
name `github.com/pritidesai/child-app`.
2. `wskdeploy` creates a directory named `Packages` and clones GitHub repo of
dependent package under `Packages`. Now the repo is cloned and renamed to
`<package>-<branch>`. Here in our example, `wskdeploy` clones repo under
`Packages` and renames it to `child-app-master`.
3. Dependent packages must have `manifest.yml`.
