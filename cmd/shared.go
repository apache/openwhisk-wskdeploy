// shared.go
package cmd

// name of manifest and deployment files
const ManifestFileNameYaml = "manifest.yaml"
const ManifestFileNameYml = "manifest.yml"
const DeploymentFileNameYaml = "deployment.yaml"
const DeploymentFileNameYml = "deployment.yml"

var cfgFile string
var CliVersion string
var CliBuild string

// used to configure service deployer for various commands
var Verbose bool
var projectPath string
var deploymentPath string
var manifestPath string
var useDefaults bool
var useInteractive bool
