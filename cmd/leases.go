package cmd

import (
	"fmt"
	"io/ioutil"

	"encoding/json"

	"github.com/Optum/dce-cli/internal/util/api"
	"github.com/pkg/browser"
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

		requestBody := &api.LeaseRequest{
			PrincipalID:              principleID,
			BudgetAmount:             budgetAmount,
			BudgetCurrency:           budgetCurrency,
			BudgetNotificationEmails: email,
		}

		leasesFullURL := *config.API.BaseURL + LeasesPath
		fmt.Println("Posting to: ", leasesFullURL)
		fmt.Println("Post body: ", requestBody)

		response := api.Request(&api.ApiRequestInput{
			Method: "POST",
			Url:    leasesFullURL,
			Region: *config.API.Region,
			Json:   requestBody,
		})

		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Response: ", response)
		fmt.Println("Response Body: ", body)
	},
}

var leasesEndCmd = &cobra.Command{
	Use:   "end",
	Short: "Cause a lease to immediately expire",
	Run: func(cmd *cobra.Command, args []string) {
		requestBody := &api.LeaseRequest{
			AccountID:   accountID,
			PrincipalID: principleID,
		}

		leasesFullURL := *config.API.BaseURL + LeasesPath
		fmt.Println("Posting to: ", leasesFullURL)
		fmt.Println("Post body: ", requestBody)

		response := api.Request(&api.ApiRequestInput{
			Method: "DELETE",
			Url:    leasesFullURL,
			Region: *config.API.Region,
			Json:   requestBody,
		})

		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Response: ", response)
		fmt.Println("Response Body: ", body)
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
			leaseLoginURL = *config.API.BaseURL + "?leaseID=" + loginLeaseID
		}

		fmt.Println("Requesting leased account credentials from: ", leaseLoginURL)
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

		if loginOpenBrowser {
			fmt.Println("Opening AWS Console in Web Browser")
			var consoleURL string

			// Build aws console url here
			consoleURL = "https://amazon.com"

			browser.OpenURL(consoleURL)
		} else {
			fmt.Println(leaseCreds)
		}
	},
}
