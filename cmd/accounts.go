package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var existingAcct string

func init() {
	accountsCmd.AddCommand(accountsListCmd)
	accountsAddCmd.Flags().StringVarP(&existingAcct, "existing", "E", "", "Bring an existing account into the accounts pool. (WARNING: Account will be nuked.)")
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
	Use:   "describe",
	Short: "describe an account",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Describe command")
	},
}

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list accounts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("List command")
	},
}

var accountsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add one or more accounts to the accounts pool.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Create command")
	},
}

var accountsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove one or more accounts from the accounts pool.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Destroy command")
	},
}
