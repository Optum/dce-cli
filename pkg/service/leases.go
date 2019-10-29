package service

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"path/filepath"

	"github.com/Optum/dce-cli/client/operations"
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
)

const LeasesPath = "/leases"

type LeasesService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
}

func (s *LeasesService) CreateLease(principleID string, budgetAmount float64, budgetCurrency string, email []string) {
	params := &operations.PostLeasesParams{
		Lease: operations.PostLeasesBody{
			PrincipalID:              &principleID,
			BudgetAmount:             &budgetAmount,
			BudgetCurrency:           &budgetCurrency,
			BudgetNotificationEmails: email,
		},
	}
	params.SetTimeout(5 * time.Second)
	res, err := apiClient.PostLeases(params, nil)
	if err != nil {
		log.Errorln("err: ", err)
	} else {
		jsonPayload, err := json.Marshal(res)
		if err != nil {
			log.Fatalln("err: ", err)
		}
		log.Infoln(string(jsonPayload))
	}
}

func (s *LeasesService) EndLease(accountID, principleID string) {
	params := &operations.DeleteLeasesParams{
		Lease: operations.DeleteLeasesBody{
			AccountID:   &accountID,
			PrincipalID: &principleID,
		},
	}
	params.SetTimeout(5 * time.Second)
	res, err := apiClient.DeleteLeases(params, nil)
	if err != nil {
		log.Errorln("err: ", err)
	} else {
		jsonPayload, err := json.Marshal(res)
		if err != nil {
			log.Fatalln("err: ", err)
		}
		log.Infoln(string(jsonPayload))
	}
}

func (s *LeasesService) GetLease(leaseID string) {
	params := &operations.GetLeasesIDParams{
		ID: leaseID,
	}
	params.SetTimeout(5 * time.Second)
	res, err := apiClient.GetLeasesID(params, nil)
	if err != nil {
		log.Errorln("err: ", err)
	} else {
		log.Infoln(res)
	}
}

func (s *LeasesService) ListLeases(acctID, principleID, nextAcctID, nextPrincipalID, leaseStatus string, pagLimit int64) {
	params := &operations.GetLeasesParams{
		AccountID:       &acctID,
		Limit:           &pagLimit,
		NextAccountID:   &nextAcctID,
		NextPrincipalID: &nextPrincipalID,
		PrincipalID:     &principleID,
		Status:          &leaseStatus,
	}
	params.SetTimeout(5 * time.Second)
	res, err := apiClient.GetLeases(params, nil)
	if err != nil {
		log.Errorln("err: ", err)
	} else {
		jsonPayload, err := json.Marshal(res.GetPayload())
		if err != nil {
			log.Fatalln("err: ", err)
		}
		log.Infoln(string(jsonPayload))
	}
}

func (s *LeasesService) LoginToLease(args []string, loginOpenBrowser bool) {

	loginLeaseID := args[0]
	leaseLoginURL := filepath.Join(*s.Config.API.Host, LeasesPath, loginLeaseID, "auth")

	log.Debugln("Requesting leased account credentials from: ", leaseLoginURL)
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

		s.Util.OpenURL(consoleURL)
	} else {
		log.Println(leaseCreds)
	}
}
