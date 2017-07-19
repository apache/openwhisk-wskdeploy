# WebHook test case for wskdeploy

### Package description

The package manifest file includes:
- An action named as `webhook` which will register a webhook in a specific github repository. The url of `webhook_trigger` is set as the webhook's configured URL. When one of the specified events is happened, a HTTP POST payload will be sent to the webhook's configured URL, i.e. the `webhook_trigger` URL.
- A feed named as `webhook_feed` which is a FEED ACTION by setting the field `action` as an action name `webhook`.
- A trigger named as `webhook_trigger`, whose URL is set as the webhook's configured URL.

The package deployment file includes:
- The namesapce
- The credential
- The parameters of feed `webhook_feed`
 - username: the username of github
 - repository: the specified repository in github
 - accessToken: the access token of github
 - events: the event of github,seperated by comma e.g. `fork`, `push`. The full list can be found [here](https://developer.github.com/webhooks/#events).
 - endpoint: the name or IP address of the OpenWhisk server where the package is deployed.

### How to deploy and test

Step 1. Update the parameters in deployment.yml to the correct values

Step 2. Deploy the package.
```
$ wskdeploy -p /tests/usecases/webhook
```

Step 3. Verify the action is installed.
```
$ wsk action list
```

Step 4. Trigger the event configured in deployment.yml

Step 5. Verify the trigger is invoked
