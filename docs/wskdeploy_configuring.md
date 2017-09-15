# Configuring ```wskdeploy```

At minimum, the wskdeploy utility needs valid OpenWhisk APIHOST and AUTH variable to attempt deployment. In this case the default target namespace is assumed; otherwise, NAMESPACE also needs to be provided.

## Precedence order

Wskdeploy attempts to find these values in the following order:

1. **Command line**

Values supplied on the command line using the ```apihost```, ```auth``` and ```namespace``` flags will override values supplied elsewhere (below).

for example the following flags can be used:

```
$ wskdeploy --apihost <host> --auth <auth> --namespace <namespace>
```

Command line values are considered higher in precedence since they indicate an intentional user action to override values that may have been supplied in files.

2. **Deployment file**

Values supplied in a Deployment YAML file will override values found elsewhere (below).

3. **Manifest file**

Values supplied in a Manifest YAML file will override values found elsewhere (below).

4. **.wskprops**

Values set using the Whisk Command Line Interface (CLI) are stored in a ```.wskprops```, typically in your $HOME directory, will override values found elsewhere (below).

It assumes that you have setup and can run the wskdeploy as described in the project README. If so, then the utility will use the OpenWhisk APIHOST and AUTH variable values in your .wskprops file to attempt deployment.

5. **Interactice mode**

If interactive mode is enabled (i.e., using the ```-i``` or ```--allow-interactive``` flags) then wskdeploy will prompt for any missing (required) values.

for example:

```
$ wskdeploy -i -m manifest.yaml
```
