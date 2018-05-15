# Licensed to the Apache Software Foundation (ASF) under one or more contributor
# license agreements; and to You under the Apache License, Version 2.0.

FROM openwhisk/nodejs6action

RUN apt-get update \
  && apt-get install -y git \
  && npm install --save simple-git js-yaml

COPY wskdeploy /nodejsAction/
