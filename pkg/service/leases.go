package service

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
	"github.com/pkg/browser"
)

const LeasesPath = "/leases"

type LeasesService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
}

func (s *LeasesService) CreateLease(principleID string, budgetAmount float64, budgetCurrency string, email []string) {
	requestBody := &utl.LeaseRequest{
		PrincipalID:              principleID,
		BudgetAmount:             budgetAmount,
		BudgetCurrency:           budgetCurrency,
		BudgetNotificationEmails: email,
	}

	leasesFullURL := *s.Config.API.BaseURL + LeasesPath
	// log.Println("Posting to: ", leasesFullURL)
	// log.Println("Post body: ", requestBody)

	response := s.Util.Request(&utl.ApiRequestInput{
		Method: "POST",
		Url:    leasesFullURL,
		Region: *s.Config.Region,
		Json:   requestBody,
	})

	// body, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode == 201 {
		log.Println("Lease created for jdoe99")
	} else {
		log.Println("DCE Responded with an error: ", response)
	}
}

func (s *LeasesService) EndLease(accountID, principleID string) {
	requestBody := &utl.LeaseRequest{
		AccountID:   accountID,
		PrincipalID: principleID,
	}

	leasesFullURL := *s.Config.API.BaseURL + LeasesPath

	response := s.Util.Request(&utl.ApiRequestInput{
		Method: "DELETE",
		Url:    leasesFullURL,
		Region: *s.Config.Region,
		Json:   requestBody,
	})

	if response.StatusCode == 200 {
		log.Println("Lease ended")
	} else {
		log.Println("DCE Responded with an error: ", response)
	}
}

func (s *LeasesService) LoginToLease(loginAcctID, loginLeaseID string, loginOpenBrowser bool) {
	if loginAcctID != "" && loginLeaseID != "" {
		log.Println("Please specify either --lease-id or --acctount-id, not both.")
		return
	}
	if loginAcctID == "" && loginLeaseID == "" {
		log.Println("Please specify either --lease-id or --acctount-id")
		return
	}
	log.Println("Logging into a leased DCE account")

	var leaseLoginURL string
	if loginAcctID != "" {
		leaseLoginURL = *s.Config.API.BaseURL + "?accountID=" + loginAcctID
	}
	if loginLeaseID != "" {
		leaseLoginURL = *s.Config.API.BaseURL + "?leaseID=" + loginLeaseID
	}

	log.Println("Requesting leased account credentials from: ", leaseLoginURL)
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
		log.Println("Opening AWS Console in Web Browser")
		var consoleURL string

		// Build aws console url here
		consoleURL = "https://amazon.com"

		browser.OpenURL(consoleURL)
	} else {
		log.Println(leaseCreds)
	}
}
