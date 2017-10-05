FROM openwhisk/nodejs6action

RUN apt-get update \
  && apt-get install -y git \
  && npm install --save simple-git js-yaml

COPY wskdeploy /nodejsAction/
