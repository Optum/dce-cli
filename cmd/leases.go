package cmd

import (
	"github.com/spf13/cobra"
)

//LeasesPath path to lease endpoint
const LeasesPath = "/leases"

var acctID string
var loginOpenBrowser bool

var principleID string
var budgetAmount float64
var budgetCurrency string
var email []string

var pagLimit int64
var nextAcctID string
var nextPrincipalID string
var leaseStatus string

func init() {
	leasesCmd.AddCommand(leasesDescribeCmd)

	var defaultPagLiimt int64 = 25
	leasesListCmd.Flags().StringVarP(&acctID, "account-id", "a", "", "An AWS Account ID")
	leasesListCmd.Flags().Int64VarP(&pagLimit, "limit", "l", defaultPagLiimt, "Max number of leases to return at once. Will include url to next page if there is one.")
	leasesListCmd.Flags().StringVarP(&nextAcctID, "next-account-id", "", "", "Account ID with which to begin the scan operation. This is used to traverse through paginated results.")
	leasesListCmd.Flags().StringVarP(&nextPrincipalID, "next-principal-id", "", "", "Principal ID with which to begin the scan operation. This is used to traverse through paginated results.")
	leasesListCmd.Flags().StringVarP(&principleID, "principle-id", "p", "", "Principle ID of a user")
	leasesListCmd.Flags().StringVarP(&leaseStatus, "status", "s", "", "Lease status")
	leasesCmd.AddCommand(leasesListCmd)

	leasesCreateCmd.Flags().StringVarP(&principleID, "principle-id", "p", "", "Principle ID for the user of the leased account")
	leasesCreateCmd.Flags().Float64VarP(&budgetAmount, "budget-amount", "b", 0, "The leased accounts budget amount")
	leasesCreateCmd.Flags().StringVarP(&budgetCurrency, "budget-currency", "c", "USD", "The leased accounts budget currency")
	leasesCreateCmd.Flags().StringArrayVarP(&email, "email", "e", nil, "The email address that budget notifications will be sent to")
	leasesCreateCmd.MarkFlagRequired("principle-id")
	leasesCreateCmd.MarkFlagRequired("budget-amount")
	leasesCreateCmd.MarkFlagRequired("budget-currency")
	leasesCreateCmd.MarkFlagRequired("email")
	leasesCmd.AddCommand(leasesCreateCmd)

	leasesEndCmd.Flags().StringVarP(&principleID, "principle-id", "p", "", "Principle ID for the user of the leased account")
	leasesEndCmd.Flags().StringVarP(&accountID, "account-id", "a", "", "Account ID associated with the lease you wish to end")
	leasesEndCmd.MarkFlagRequired("principle-id")
	leasesEndCmd.MarkFlagRequired("account-id")
	leasesCmd.AddCommand(leasesEndCmd)

	leasesLoginCmd.Flags().BoolVarP(&loginOpenBrowser, "open-browser", "b", false, "Opens web broswer to AWS console instead of printing credentials")
	leasesCmd.AddCommand(leasesLoginCmd)

	RootCmd.AddCommand(leasesCmd)
}

var leasesCmd = &cobra.Command{
	Use:   "leases",
	Short: "Manage dce leases",
}

var leasesDescribeCmd = &cobra.Command{
	Use:   "describe [Lease ID]",
	Short: "describe a lease",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.GetLease(args[0])
	},
}

var leasesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List leases using various query filters.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		service.ListLeases(acctID, principleID, nextAcctID, nextPrincipalID, leaseStatus, pagLimit)
	},
}

var leasesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a lease.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		service.CreateLease(principleID, budgetAmount, budgetCurrency, email)
	},
}

var leasesEndCmd = &cobra.Command{
	Use:   "end [Lease ID]",
	Short: "Cause a lease to immediately expire",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		service.EndLease(accountID, principleID)
	},
}

var leasesLoginCmd = &cobra.Command{
	Use:   "login [Lease ID]",
	Short: "Login to a leased DCE account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.LoginToLease(args[0], loginOpenBrowser)
	},
}
