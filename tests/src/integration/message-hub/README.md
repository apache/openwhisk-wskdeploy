# Test Case of message hub

This is a test case for message hub. Before running this use case, please make sure you have set the following
environment variables on your machine: MESSAGEHUB_ADMIN_HOST, KAFKA_BROKERS_SASL, SOURCE_TOPIC and DESTINATION_TOPIC.
 
The environment variables, SOURCE_TOPIC and DESTINATION_TOPIC, are two topic names in message hub service. Both of them
must be available in the message hub service you are about to use. The variable MESSAGEHUB_ADMIN_HOST specifies the url
link of the admin host for the message hub. The variable KAFKA_BROKERS_SASL specifies the array of the kafka brokers, e.g.
[kafka01-prod01.messagehub.services.net:9093 kafka02-prod01.messagehub.services.net:9093].

It can be deployed and tested with:

```bash
$ wskdeploy -p tests/src/integration/message-hub
```
