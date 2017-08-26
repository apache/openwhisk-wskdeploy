# Integration Test - Dependencies

`wskdeploy` supports dependencies where it allows you to declare other OpenWhisk
packages that your application or project (manifest) is dependent on. With declaring
dependent packages, `wskdeploy` supports automatic deployment of those dependent
packages.

For example, our `root-app` application is dependent on `child-app` package located
at https://github.com/pritidesai/child-app.
We can declare this dependency in `manifest.yaml` with:

```yaml
package:
  name: root-app
  namespace: guest
  dependencies:
    child-app:
      location: github.com/pritidesai/child-app
  triggers:
    trigger1:
  rules:
    rule1:
      trigger: trigger1
      action: child-app/hello
```

**Note:**

1. Package name of the dependent package `child-app` should match GitHub repo
name `github.com/pritidesai/child-app`.
2. `wskdeploy` creates a directory named `Packages` and clones GitHub repo of
dependent package under `Packages`. Now the repo is cloned and renamed to
`<package>-<branch>`. Here in our example, `wskdeploy` clones repo under
`Packages` and renames it to `child-app-master`.
3. Dependent packages must have `manifest.yml`.
