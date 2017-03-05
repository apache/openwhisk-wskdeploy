# Test case for Whisk Deploy

This is a test case for wskdeploy. There is an action named as "hello" in this package. It accepts two parameters "name" and "place" and will return a greeting message "Hello, name from place!"

It can be tested as below:
$ wskdeploy -p tests/testcases/helloworld
$ wsk action list
$ wsk action invoke --blocking --result helloworld/hello --param name Bernie --param place Vermont
