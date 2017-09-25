# Debugging ```wskdeploy```

The Whisk Deploy utility provides several ways to help you in debugging your OpenWhisk application or package during parsing, deployment or undeployment.

## Enabling Verbose mode

The first thing you should do is turn on _"verbose mode"_ using the flag ```-v``` or ```--verbose```.  This will assure that all Informational messages within the code will be displayed.

```
$ wskdeploy -v -m manifest.yaml
```

## Enable console logging in your Action

You may call ```console.log(<text>)``` within your Action (function) code to aid in debugging.  For example, in NodeJS (JavaScript) you could output your entire JSON payload before returning it:
```
function main(params) {
    msg = "Hello, " + params.name + " from " + params.place;
    console.log(msg)
    return { payload:  msg };
}
```

## Enable additional trace in Go Client

Wskdeploy uses the OpenWhisk GoLang Client to format and invoke OpenWhisk's APIs which has additional debug tracing available.

To enable this trace, set the following environment variable in Bash:

```
# set to any value > 0
WSK_CLI_DEBUG=1
```

## Pay attention to Named error messages

Wskpdeloy uses named errors that describe the type of any error found along with additional values that correspond with an error.

For example, if you have an error in your Manifest's YAML, you may see an error such as:
```
[50]: Invalid input of Yaml file =====> incubator-openwhisk-wskdeploy/parsers/manifest_parser.go
[98]: Failed to parse the yaml file manifest_bad_yaml.yaml
 =====> yaml: line 13: could not find expected ':'
```

The named error **NewInputYamlFormatError** provides direct indication of both where in the utilities GoLang code the error was reported, but also details provided from the YAML parser regarding where the Manifest file may contain a grammar error.


All current named errors supported by the utility can be found in the latest ```wskdeployerror.go``` source file:
[wskdeployerror.go](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/utils/wskdeployerror.go)
