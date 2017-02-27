# Test case for WskTool.

This is a test case for WskTool. There is an action named as "hello" in this package. It accepts two parameters "name" and "place" and will return a greeting message "Hello, name from place!"

It can be tested as below:
$ wsktool deploy -p /tests/testcases/helloworld
$ wsk action list
$ wsk action invoke --blocking --result hello --param name Bernie --param place Vermont
