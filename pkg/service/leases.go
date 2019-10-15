package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Optum/dce-cli/configs"
	utl "github.com/Optum/dce-cli/internal/util"
	"github.com/pkg/browser"
)

const LeasesPath = "/leases"

type LeasesService struct {
	Config *configs.Root
	Util   *utl.UtilContainer
}

func (s *LeasesService) CreateLease(principleID string, budgetAmount float64, budgetCurrency string, email []string) {
	requestBody := &utl.LeaseRequest{
		PrincipalID:              principleID,
		BudgetAmount:             budgetAmount,
		BudgetCurrency:           budgetCurrency,
		BudgetNotificationEmails: email,
	}

	leasesFullURL := *s.Config.API.BaseURL + LeasesPath
	fmt.Println("Posting to: ", leasesFullURL)
	fmt.Println("Post body: ", requestBody)

	response := s.Util.Request(&utl.ApiRequestInput{
		Method: "POST",
		Url:    leasesFullURL,
		Region: *s.Config.Region,
		Json:   requestBody,
	})

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("Response: ", response)
	fmt.Println("Response Body: ", body)
}

func (s *LeasesService) EndLease(accountID, principleID string) {
	requestBody := &utl.LeaseRequest{
		AccountID:   accountID,
		PrincipalID: principleID,
	}

	leasesFullURL := *s.Config.API.BaseURL + LeasesPath
	fmt.Println("Posting to: ", leasesFullURL)
	fmt.Println("Post body: ", requestBody)

	response := s.Util.Request(&utl.ApiRequestInput{
		Method: "DELETE",
		Url:    leasesFullURL,
		Region: *s.Config.Region,
		Json:   requestBody,
	})

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("Response: ", response)
	fmt.Println("Response Body: ", body)
}

func (s *LeasesService) LoginToLease(loginAcctID, loginLeaseID string, loginOpenBrowser bool) {
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
		leaseLoginURL = *s.Config.API.BaseURL + "?accountID=" + loginAcctID
	}
	if loginLeaseID != "" {
		leaseLoginURL = *s.Config.API.BaseURL + "?leaseID=" + loginLeaseID
	}

	fmt.Println("Requesting leased account credentials from: ", leaseLoginURL)
	response := s.Util.Request(&utl.ApiRequestInput{
		Method: "GET",
		Url:    leaseLoginURL,
		Region: *s.Config.Region,
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
}
