# ```wskdeploy``` utility FAQ

### What if ```wskdeploy``` finds an error in my manifest?

- The ```wskdeploy``` utility will not attempt to deploy a package if an error in the manifest is detected, but will report as much information as it can to help you locate the error in the YAML file.

### What if ```wskdeploy``` encounters an error during deployment?

-  The ```wskdeploy``` utility will cease deploying as soon as it receives an error from the target platform and display what error information it receives to you.
- then it will attempt to undeploy any entities that it attempted to deploy.
  - If "interactive mode" was used to deploy, then you will be prompted to confirm you wish to undeploy.
