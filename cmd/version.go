package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of openwhisk-wskdeploy",
	Long:  `Print the version number of openwhisk-wskdeploy`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("openwhisk-wskdeploy version is %s--%s\n", CliBuild, CliVersion)
	},
}
