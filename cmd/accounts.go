package cmd

import (
	"github.com/spf13/cobra"
)

const accountsPath = "/accounts"

var accountID string
var adminRoleARN string

func init() {
	accountsCmd.AddCommand(accountsListCmd)

	accountsAddCmd.Flags().StringVarP(&accountID, "account-id", "a", "", "The ID of the existing account to add to the DCE accounts pool (WARNING: Account will be nuked.)")
	accountsAddCmd.Flags().StringVarP(&adminRoleARN, "admin-role-arn", "r", "", "The admin role arn to be assumed by the DCE master account. Trust policy must be configured with DCE master account as trusted entity.")
	accountsCmd.AddCommand(accountsAddCmd)

	accountsCmd.AddCommand(accountsRemoveCmd)
	accountsCmd.AddCommand(accountsDescribeCmd)
	RootCmd.AddCommand(accountsCmd)
}

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage dce accounts",
}

var accountsDescribeCmd = &cobra.Command{
	Use:   "describe [Accound ID]",
	Short: "describe an account",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.GetAccount(args[0])
	},
}

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list accounts",
	Run: func(cmd *cobra.Command, args []string) {
		service.ListAccounts()
	},
}

var accountsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an account to the accounts pool",
	Run: func(cmd *cobra.Command, args []string) {
		service.AddAccount(accountID, adminRoleARN)
	},
}

var accountsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an account from the accounts pool.",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.RemoveAccount(args[0])
	},
}
