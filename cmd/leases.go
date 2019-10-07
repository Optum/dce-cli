package cmd

import (
	"fmt"
	"io/ioutil"

	"encoding/json"

	api "github.com/Optum/dce-cli/internal/api"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var loginAcctID string
var loginLeaseID string
var loginPrintCredentials bool

func init() {
	leasesCmd.AddCommand(leasesDescribeCmd)
	leasesCmd.AddCommand(leasesListCmd)
	leasesCmd.AddCommand(leasesCreateCmd)
	leasesCmd.AddCommand(leasesDestroyCmd)

	leasesLoginCmd.Flags().StringVarP(&loginAcctID, "acctount-id", "a", "", "Account ID to login to")
	leasesLoginCmd.Flags().StringVarP(&loginLeaseID, "lease-id", "l", "", "Lease ID for the account to login to")
	leasesLoginCmd.Flags().BoolVarP(&loginPrintCredentials, "print-credentials", "c", true, "Prints temporary credentials instead of opening AWS console in a web browser")
	leasesCmd.AddCommand(leasesLoginCmd)

	RootCmd.AddCommand(leasesCmd)
}

var leasesCmd = &cobra.Command{
	Use:   "leases",
	Short: "Manage dce leases",
}

var leasesDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "describe a lease",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Describe command")
	},
}

var leasesListCmd = &cobra.Command{
	Use:   "list",
	Short: "list leases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("List command")
	},
}

var leasesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a lease.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Create command")
	},
}

var leasesDestroyCmd = &cobra.Command{
	Use:   "end",
	Short: "Cause a lease to immediately expire",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Destroy command")
	},
}

var leasesLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a leased DCE account",
	Run: func(cmd *cobra.Command, args []string) {
		if loginAcctID != "" && loginLeaseID != "" {
			fmt.Println("Please specify either --lease-id or --acctount-id, not both.")
			return
		}
		if loginAcctID == "" && loginLeaseID == "" {
			fmt.Println("Please specify either --lease-id or --acctount-id")
			return
		}
		fmt.Println("Logging into a leased DCE account")

		var leaseLoginURL string
		if loginAcctID != "" {
			leaseLoginURL = *config.API.BaseURL + "?accountID=" + loginAcctID
		}
		if loginLeaseID != "" {
			leaseLoginURL = *config.API.BaseURL + "?leaseID=" + loginAcctID
		}

		response := api.Request(&api.ApiRequestInput{
			Method: "GET",
			Url:    leaseLoginURL,
			Region: *config.API.Region,
		})

		leaseCreds := struct {
			AwsAccessKeyID     string
			AwsSecretAccessKey string
			AwsSessionToken    string
		}{}

		body, _ := ioutil.ReadAll(response.Body)

		// Some test data. Remove once integrated with api.
		body = []byte("{\"AwsAccessKeyID\": \"AKD\", \"AwsSecretAccessKey\": \"ASK\", \"AwsSessionToken\": \"AST\" }")
		json.Unmarshal(body, &leaseCreds)

		if loginPrintCredentials {
			fmt.Println("Requesting leased account credentials from: ", leaseLoginURL)
			fmt.Println(leaseCreds)
		} else {
			fmt.Println("Requesting leased account credentials from: ", leaseLoginURL)
			fmt.Println("Opening AWS Console in Web Browser")
			var consoleURL string

			// Build aws console url here

			browser.OpenURL(consoleURL)
		}
	},
}
