# Test case for WskTool.

This is a test case for WskTool. The package includes:
- An action named as "hello". It accepts two parameters "name" and "place" and will return a greeting message "Hello, name from place!"
- A trigger named as "locationUpdate"
- A rule to associate the trigger with the action.

It can be tested as below:
Step 1. Deploy the package.
$ wsktool deploy -p /tests/testcases/helloworld

Step 2. Verify the action is installed.
$ wsk action list

Step 3. Verify the action is invoked.
$ wsk activation list --limit 1 hello
$ wsk activation result <your action ID>

Step 4. Invoke the trigger
$ wsk trigger fire locationUpdate --param name Bernie --param place "Washington, D.C."

Step 5. Verify the action is invoked
$ wsk activation list --limit 1 hello
$ wsk activation result <your action ID>

