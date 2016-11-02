# Test case for WskTool.

This is a test case for WskTool. The package includes:
- A feed named as "webhook_feed". It will add a webhook to a specified github repository. 
- A trigger named as "webhook_trigger". When there is some specified actions happened, the trigger will be fired.

It can be tested as below:
Step 1. Deploy the package.
$ wsktool deploy -p /tests/testcases/webhook

Step 2. Verify the action is installed.
$ wsk action list

Step 3. Push a patch to the repository

Step 4. Verify the trigger is invoked

