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

# How to generate the file i18n_resources.go for internationalization

As a contributor to wskdeploy, the file *i18n_resources.go* needs to regenerated,
when you add or change any localized message.

### Install go-bindata
In order to generate i18n_resources.go, you need to install go-bindata first:

```
$ go get -u github.com/jteeuwen/go-bindata/...
```

### Generate i18n_resources.go
Then, go the HOME directory of wskdeploy and run the following command:

```
$ $GOPATH/bin/go-bindata -pkg wski18n -o wski18n/i18n_resources.go wski18n/resources;
```

Finally, add the default ASF license header to i18n_resources.go. Since each file of
source code starts with the ASF license header, you need to add it to i18n_resources.go
each time it is regenerated. You can find this license header in any other file of source
code, e.g. i18n.go.
