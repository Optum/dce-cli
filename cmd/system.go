package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dceVersion string

func init() {

	systemDeployCmd.Flags().StringVarP(&dceVersion, "dce-version", "s", "", "DCE version to deploy (Defaults to latest)")
	systemCmd.AddCommand(systemDeployCmd)

	systemLogsCmd.AddCommand(systemLogsAccountsCmd)
	systemLogsCmd.AddCommand(systemLogsLeasesCmd)
	systemLogsCmd.AddCommand(systemLogsUsageCmd)
	systemLogsCmd.AddCommand(systemLogsResetCmd)
	systemCmd.AddCommand(systemLogsCmd)

	systemUsersCmd.AddCommand(systemUsersAddCmd)
	systemUsersCmd.AddCommand(systemUsersRemoveCmd)
	systemCmd.AddCommand(systemUsersCmd)

	systemCmd.AddCommand(systemInitCmd)

	RootCmd.AddCommand(systemCmd)
}

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Deploy and configure the DCE system",
}

var systemInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a DCE configuration file at the specified location (Defaults to ~/.dce.yaml)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Init command")
	},
}

/*
Deploy Namespace
*/

var systemDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the DCE system",
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
		fmt.Print("Accounts command")
	},
}

var systemLogsLeasesCmd = &cobra.Command{
	Use:   "leases",
	Short: "View lease logs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Leases command")
	},
}

var systemLogsUsageCmd = &cobra.Command{
	Use:   "usage",
	Short: "View usage logs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Usage command")
	},
}

var systemLogsResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "View reset logs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Reset command")
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
		fmt.Print("Add command")
	},
}

var systemUsersRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove users",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Remove command")
	},
}
