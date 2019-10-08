package cmd

import (
	"fmt"
	"io/ioutil"

	"encoding/json"

	"github.com/Optum/dce-cli/internal/api"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

const LeasesPath = "/leases"

var loginAcctID string
var loginLeaseID string
var loginOpenBrowser bool

func init() {
	leasesCmd.AddCommand(leasesDescribeCmd)
	leasesCmd.AddCommand(leasesListCmd)
	leasesCmd.AddCommand(leasesCreateCmd)
	leasesCmd.AddCommand(leasesDestroyCmd)

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

		type requestBody struct {
			PrincipalID              string   `json:"principalId"`
			AccountID                string   `json:"accountId"`
			BudgetAmount             float64  `json:"budgetAmount"`
			BudgetCurrency           string   `json:"budgetCurrency"`
			BudgetNotificationEmails []string `json:"budgetNotificationEmails"`
		}

		postBody := &requestBody{
			PrincipalID:              "abc",
			BudgetAmount:             350,
			BudgetCurrency:           "USD",
			BudgetNotificationEmails: []string{"test@test.com"},
		}

		leasesPath := *config.API.BaseURL + LeasesPath
		fmt.Println("Posting to: ", leasesPath)
		fmt.Println("Post body: ", postBody)

		response := api.Request(&api.ApiRequestInput{
			Method: "POST",
			Url:    *config.API.BaseURL + LeasesPath,
			Region: *config.API.Region,
			Json:   postBody,
		})

		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Response: ", response)
		fmt.Println("Response Body: ", body)
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
