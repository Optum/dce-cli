package service

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Optum/dce-cli/client/operations"
	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
)

const LeasesPath = "/leases"

type LeasesService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
}

func (s *LeasesService) CreateLease(principalID string, budgetAmount float64, budgetCurrency string, email []string, expiresOn string) {
	postBody := operations.PostLeasesBody{
		PrincipalID:              &principalID,
		BudgetAmount:             &budgetAmount,
		BudgetCurrency:           &budgetCurrency,
		BudgetNotificationEmails: email,
	}

	expiry, err := s.Util.ExpandEpochTime(expiresOn)

	if err != nil && expiry > 0 {
		expiryf := float64(expiry)
		postBody.ExpiresOn = expiryf
	}

	params := &operations.PostLeasesParams{
		Lease: postBody,
	}
	params.SetTimeout(5 * time.Second)
	res, err := apiClient.PostLeases(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	}
	jsonPayload, err := json.MarshalIndent(res.GetPayload(), "", "\t")
	if err != nil {
		log.Fatalln("err: ", err)
	}
	log.Infoln("Lease created:", string(jsonPayload))
}

func (s *LeasesService) EndLease(accountID, principalID string) {
	params := &operations.DeleteLeasesParams{
		Lease: operations.DeleteLeasesBody{
			AccountID:   &accountID,
			PrincipalID: &principalID,
		},
	}
	params.SetTimeout(5 * time.Second)
	_, err := apiClient.DeleteLeases(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	}
	log.Infoln("Lease ended")
}

func (s *LeasesService) GetLease(leaseID string) {
	params := &operations.GetLeasesIDParams{
		ID: leaseID,
	}
	params.SetTimeout(5 * time.Second)
	res, err := apiClient.GetLeasesID(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	}
	jsonPayload, err := json.MarshalIndent(res.GetPayload(), "", "\t")
	if err != nil {
		log.Fatalln("err: ", err)
	}
	log.Infoln(string(jsonPayload))

}

func (s *LeasesService) ListLeases(acctID, principalID, nextAcctID, nextPrincipalID, leaseStatus string, pagLimit int64) {
	params := &operations.GetLeasesParams{
		AccountID:       &acctID,
		Limit:           &pagLimit,
		NextAccountID:   &nextAcctID,
		NextPrincipalID: &nextPrincipalID,
		PrincipalID:     &principalID,
		Status:          &leaseStatus,
	}
	params.SetTimeout(5 * time.Second)
	res, err := apiClient.GetLeases(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	}
	jsonPayload, err := json.MarshalIndent(res.GetPayload(), "", "\t")
	if err != nil {
		log.Fatalln("err: ", err)
	}
	log.Infoln(string(jsonPayload))
}

func (s *LeasesService) LoginToLease(leaseID, loginProfile string, loginOpenBrowser, loginPrintCreds bool) {
	log.Debugln("Requesting leased account credentials")
	params := &operations.PostLeasesIDAuthParams{
		ID: leaseID,
	}
	params.SetTimeout(5 * time.Second)
	res, err := apiClient.PostLeasesIDAuth(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	} else {
		jsonPayload, err := json.MarshalIndent(res.GetPayload(), "", "\t")
		if err != nil {
			log.Fatalln("err: ", err)
		}
		log.Debug(string(jsonPayload))
	}

	responsePayload := res.GetPayload()

	if !(loginOpenBrowser || loginPrintCreds) {
		credsPath := filepath.Join(".aws", "credentials")
		log.Infoln("Adding credentials to " + credsPath + " using AWS CLI")
		s.Util.ConfigureAWSCLICredentials(responsePayload.AccessKeyID,
			responsePayload.SecretAccessKey,
			responsePayload.SessionToken,
			loginProfile)

	} else if loginProfile != "default" {
		log.Infoln("Setting --profile has no effect when used with other flags.\n")
	}

	if loginOpenBrowser {
		log.Infoln("Opening AWS Console in Web Browser")
		s.Util.OpenURL(responsePayload.ConsoleURL)
	}

	if loginPrintCreds {
		creds := fmt.Sprintf(constants.CredentialsExport,
			responsePayload.AccessKeyID,
			responsePayload.SecretAccessKey,
			responsePayload.SessionToken)
		log.Infoln(creds)
	}
}
