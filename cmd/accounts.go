package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/Optum/dce-cli/internal/api"
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

		requestBody := &api.CreateAccountRequest{
			ID:           accountID,
			AdminRoleArn: adminRoleARN,
		}

		accountsFullURL := *config.API.BaseURL + accountsPath
		fmt.Println("Posting to: ", accountsFullURL)
		fmt.Println("Post body: ", requestBody)

		response := api.Request(&api.ApiRequestInput{
			Method: "POST",
			Url:    accountsFullURL,
			Region: *config.API.Region,
			Json:   requestBody,
		})

		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Response: ", response)
		fmt.Println("Response Body: ", body)

	},
}

var accountsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove one or more accounts from the accounts pool.",
	Run: func(cmd *cobra.Command, args []string) {
		accountsFullURL := *config.API.BaseURL + accountsPath + "/" + accountID
		fmt.Println("Posting to: ", accountsFullURL)

		response := api.Request(&api.ApiRequestInput{
			Method: "DELETE",
			Url:    accountsFullURL,
			Region: *config.API.Region,
		})

		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Response: ", response)
		fmt.Println("Response Body: ", body)
	},
}
