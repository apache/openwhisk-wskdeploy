// shared.go
package cmd

const ManifestFileName = "manifest\\.yaml|\\.yml"
const DeploymentFileName = "deployment\\.yaml|\\.yml"

var cfgFile string
var CliVersion string
var CliBuild string

// used to configure service deployer for various commands
var Verbose bool
var projectPath string
var deploymentPath string
var manifestPath string
var useDefaults string
var useInteractive bool
