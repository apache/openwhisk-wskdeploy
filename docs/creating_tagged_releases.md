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

# Creating Tagged Releases of ```wskdeploy```

The most convenient way to create a tagged release for wskdeploy is to build the binaries by adding tag to upstream master.


1. Add a tag to a commit id: ```git tag -a <tag/version> <commit hash>```

for example, using the leading characters on commit hash (minimum of 7 characters):

```
$ git tag -a 1.0.9 c08b0f3
```

2. Push the tag upstream: ```git push -f upstream <tag/version>```

for example:
```
$ git push -f upstream 1.0.9
```

will start a Travis build of 1.0.9 automatically from seeing the tag creation event.

If the travis build passed, binaries will be pushed into releases page.

If we modify the tag by pointing to a different commit, use ```git push -f upstream 1.0.9<tag>``` to overwrite the old tag. New binaries from travis build will overwrite the old binaries as well.

You can download the binaries, and delete them from the releases page in GitHub if we do not want them to be public.

# Publishing Tagged Release to Homebrew

[Homebrew](https://brew.sh) is used to install `wskdeploy` locally. Once we release a new version of `wskdeploy` we should update its version in homebrew.

Get the new release SHA256 checksum by downloading the Source Code (tar.gz) from the [releases page](https://github.com/apache/openwhisk-wskdeploy/releases) and running `shasum -a 256 X.Y.Z.tar.gz` on the tarball.

Update brew formula with the automation command `brew bump-formula-pr`:
```bash
$ brew bump-formula-pr \
  --url='https://github.com/apache/openwhisk-wskdeploy/archive/X.Y.Z.tar.gz' \
  --sha256='PASTE THE SHA256 CHECKSUM HERE' \
  --version='X.Y.Z' \
  wskdeploy
```
