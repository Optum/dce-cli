package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const accountsPath = "/accounts"

var accountID string
var adminRoleARN string
var newAcct bool

func init() {
	accountsCmd.AddCommand(accountsListCmd)

	accountsAddCmd.Flags().StringVarP(&accountID, "account-id", "a", "", "The ID of the existing account to add to the DCE accounts pool (WARNING: Account will be nuked.)")
	accountsAddCmd.Flags().StringVarP(&adminRoleARN, "admin-role-arn", "r", "", "The admin role arn to be assumed by the DCE master account. Trust policy must be configured with DCE master account as trusted entity.")
	accountsAddCmd.Flags().BoolVarP(&newAcct, "new", "n", false, "Create a new account rather than specifying an exiting one.")
	accountsCmd.AddCommand(accountsAddCmd)

	accountsRemoveCmd.Flags().StringVarP(&accountID, "account-id", "a", "", "The ID of the account to remove from the accounts pool.")
	accountsCmd.AddCommand(accountsRemoveCmd)
	accountsCmd.AddCommand(accountsDescribeCmd)
	RootCmd.AddCommand(accountsCmd)

	// TODO: Configure util for this command to use local env credentials
}

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage dce accounts",
}

var accountsDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "describe an account",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("TODO")
	},
}

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list accounts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("TODO")
	},
}

var accountsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add one or more accounts to the accounts pool.",
	Run: func(cmd *cobra.Command, args []string) {
		service.AddAccount(accountID, adminRoleARN)
	},
}

var accountsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove one or more accounts from the accounts pool.",
	Run: func(cmd *cobra.Command, args []string) {
		service.RemoveAccount(accountID)
	},
}
