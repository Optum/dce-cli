package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//LeasesPath path to lease endpoint
const LeasesPath = "/leases"

var loginAcctID string
var loginLeaseID string
var loginOpenBrowser bool

var principleID string
var budgetAmount float64
var budgetCurrency string
var email []string

func init() {
	leasesCmd.AddCommand(leasesDescribeCmd)
	leasesCmd.AddCommand(leasesListCmd)

	leasesCreateCmd.Flags().StringVarP(&principleID, "principle-id", "p", "", "Principle ID for the user of the leased account")
	leasesCreateCmd.Flags().Float64VarP(&budgetAmount, "budget-amount", "b", 0, "The leased accounts budget amount")
	leasesCreateCmd.Flags().StringVarP(&budgetCurrency, "budget-currency", "a", "USD", "The leased accounts budget currency")
	leasesCreateCmd.Flags().StringArrayVarP(&email, "email", "e", nil, "The email address that budget notifications will be sent to")
	leasesCmd.AddCommand(leasesCreateCmd)

	leasesEndCmd.Flags().StringVarP(&principleID, "principle-id", "p", "", "Principle ID for the user of the leased account")
	leasesEndCmd.Flags().StringVarP(&accountID, "account-id", "a", "", "Account ID associated with the lease you wish to end")
	leasesCmd.AddCommand(leasesEndCmd)

	leasesLoginCmd.Flags().StringVarP(&loginAcctID, "account-id", "a", "", "Account ID to login to")
	leasesLoginCmd.Flags().StringVarP(&loginLeaseID, "lease-id", "l", "", "Lease ID for the account to login to")
	leasesLoginCmd.Flags().BoolVarP(&loginOpenBrowser, "open-browser", "b", false, "Opens web broswer to AWS console instead of printing credentials")
	leasesCmd.AddCommand(leasesLoginCmd)

	RootCmd.AddCommand(leasesCmd)

	// TODO: Configure util for this command to use local env credentials
}

var leasesCmd = &cobra.Command{
	Use:   "leases",
	Short: "Manage dce leases",
}

var leasesDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "describe a lease",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("TODO")
	},
}

var leasesListCmd = &cobra.Command{
	Use:   "list",
	Short: "list leases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("TODO")
	},
}

var leasesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a lease.",
	Run: func(cmd *cobra.Command, args []string) {
		service.CreateLease(principleID, budgetAmount, budgetCurrency, email)
	},
}

var leasesEndCmd = &cobra.Command{
	Use:   "end",
	Short: "Cause a lease to immediately expire",
	Run: func(cmd *cobra.Command, args []string) {
		service.EndLease(accountID, principleID)
	},
}

var leasesLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a leased DCE account",
	Run: func(cmd *cobra.Command, args []string) {
		service.LoginToLease(loginAcctID, loginLeaseID, loginOpenBrowser)
	},
}
