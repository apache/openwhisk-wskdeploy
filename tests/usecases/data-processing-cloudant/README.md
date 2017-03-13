# Test case for data-processing-cloudant.

This is use case for the data-processing-cloudant 101 use case.

It can be tested as below:
Step 1. Update the parameters in deployment.yml to the correct values

Step 2. Deploy the package.
$ go run main.go  -d deployment.yaml  -m manifest.yaml

Step 3. Verify the package is binded to /whisk.system/cloudant, and actions is installed.
$ wsk action list

Step 4. Trigger the event configured in deployment.yml

Step 5. Verify the trigger is invoked