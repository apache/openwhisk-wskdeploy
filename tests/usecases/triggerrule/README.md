# TriggerRule use case for wskdeploy

### Package description

This Package named `triggerrule` includes:
- An action named "greeting". It accepts two parameters "name" and "place" and will return a greeting message "Hello, name from place!"
- A trigger named as "locationUpdate"
- A rule named "myRule" to associate the trigger with the action.

### How to deploy and test

#### Step 1. Deploy the package.

```
$ wskdeploy -p tests/usecases/triggerrule
```

#### Step 2. Verify the action is installed.

```
$ wsk action list

{
  "name": "greeting",
  "publish": false,
  "annotations": [{
    "key": "exec",
    "value": "nodejs:6"
  }],
  "version": "0.0.1",
  "namespace": "<namespace>/triggerrule"
}
```

#### Step 3. Verify the action's last activation (ID) before invoking.

```
$ wsk activation list --limit 1 greeting

activations
3b80ec50d934464faedc04bc991f8673 greeting
```

#### Step 4. Invoke the trigger

```
$ wsk trigger fire locationUpdate --param name Bernie --param place "Washington, D.C."

ok: triggered /_/locationUpdate with id 5f8f26928c384afd85f47281bf85b302
```

#### Step 5. Verify the action is invoked

```
$ wsk activation list --limit 2 greeting

activations
a2b5cb1076b24c1885c018bb46f4e990 greeting
3b80ec50d934464faedc04bc991f8673 greeting
```

Indeed a new activation was listed.
