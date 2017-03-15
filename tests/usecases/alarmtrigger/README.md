# Test Case for Whisk Deploy

This is a test case for `wskdeploy`. This package demonstrates how to create alarm trigger. You have to specify `/whisk.system/alarms/alarm` as a `source` for alarm trigger. It takes one mandatory parameter `cron` in deployment file.

It can be deployed and tested with:

```bash
$ wskdeploy -p tests/usecases/alarmtrigger
$ wsk activation poll
$ wsk trigger fire Every12Hours
$ wsk activation get <HelloWorldActivationID>
```  

