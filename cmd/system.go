package cmd

import (
	"github.com/spf13/cobra"
)

var deployLocalPath string
var deployNamespace string
var dceRepoPath string

func init() {
	systemDeployCmd.Flags().StringVarP(&deployLocalPath, "local", "l", "", "Path to a local DCE repo to deploy.")
	systemDeployCmd.Flags().StringVarP(&deployNamespace, "namespace", "n", "", "Set a custom terraform namespace (Optional)")
	systemCmd.AddCommand(systemDeployCmd)

	RootCmd.AddCommand(systemCmd)

}

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Deploy and configure the DCE system",
}

/*
Deploy Namespace
*/

var systemDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the DCE system",
	Run: func(cmd *cobra.Command, args []string) {
		service.Deploy(deployNamespace, deployLocalPath)
	},
}
