package cmd

import (
	"github.com/spf13/cobra"
)

var deployNamespace string
var dceRepoPath string

func init() {
	systemDeployCmd.Flags().StringVarP(&deployNamespace, "namespace", "n", "", "Set a custom terraform namespace (Optional)")
	systemDeployCmd.Flags().StringVarP(&dceRepoPath, "path", "p", "", "Path to local DCE repo")
	systemCmd.AddCommand(systemDeployCmd)

	systemLogsCmd.AddCommand(systemLogsAccountsCmd)
	systemLogsCmd.AddCommand(systemLogsLeasesCmd)
	systemLogsCmd.AddCommand(systemLogsUsageCmd)
	systemLogsCmd.AddCommand(systemLogsResetCmd)
	systemCmd.AddCommand(systemLogsCmd)

	systemUsersCmd.AddCommand(systemUsersAddCmd)
	systemUsersCmd.AddCommand(systemUsersRemoveCmd)
	systemCmd.AddCommand(systemUsersCmd)

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
		service.Deploy(deployNamespace)
	},
}

/*
Logs Namespace
*/

var systemLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View logs",
}

var systemLogsAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "View account logs",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("TODO")
	},
}

var systemLogsLeasesCmd = &cobra.Command{
	Use:   "leases",
	Short: "View lease logs",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("TODO")
	},
}

var systemLogsUsageCmd = &cobra.Command{
	Use:   "usage",
	Short: "View usage logs",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("TODO")
	},
}

var systemLogsResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "View reset logs",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("TODO")
	},
}

/*
Users Namespace
*/
var systemUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
}

var systemUsersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add users",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("TODO")
	},
}

var systemUsersRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove users",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("TODO")
	},
}
