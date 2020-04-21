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

	if err == nil && expiry > 0 {
		expiryf := float64(expiry)
		postBody.ExpiresOn = expiryf
	}

	params := &operations.PostLeasesParams{
		Lease: postBody,
	}
	params.SetTimeout(5 * time.Second)
	res, err := ApiClient.PostLeases(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	}
	jsonPayload, err := json.MarshalIndent(res.GetPayload(), "", "\t")
	if err != nil {
		log.Fatalln("err: ", err)
	}

	if _, err := Out.Write(jsonPayload); err != nil {
		log.Fatalln("err: ", err)

	}
}

func (s *LeasesService) EndLease(leaseID, accountID, principalID string) {
	var err error = nil
	if leaseID != "" {
		params := &operations.DeleteLeasesIDParams{
			ID: leaseID,
		}
		params.SetTimeout(5 * time.Second)
		_, err = ApiClient.DeleteLeasesID(params, nil)
	} else if accountID != "" && principalID != "" {
		params := &operations.DeleteLeasesParams{
			Lease: operations.DeleteLeasesBody{
				AccountID:   &accountID,
				PrincipalID: &principalID,
			},
		}
		params.SetTimeout(5 * time.Second)
		_, err = ApiClient.DeleteLeases(params, nil)
	}

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
	res, err := ApiClient.GetLeasesID(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	}
	jsonPayload, err := json.MarshalIndent(res.GetPayload(), "", "\t")
	if err != nil {
		log.Fatalln("err: ", err)
	}
	if _, err := Out.Write(jsonPayload); err != nil {
		log.Fatalln("err: ", err)
	}
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
	res, err := ApiClient.GetLeases(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	}
	jsonPayload, err := json.MarshalIndent(res.GetPayload(), "", "\t")
	if err != nil {
		log.Fatalln("err: ", err)
	}
	if _, err := Out.Write(jsonPayload); err != nil {
		log.Fatalln("err: ", err)
	}
}

type leaseCreds struct {
	AccessKeyID     string  `json:"accessKeyId,omitempty"`
	ConsoleURL      string  `json:"consoleUrl,omitempty"`
	ExpiresOn       float64 `json:"expiresOn,omitempty"`
	SecretAccessKey string  `json:"secretAccessKey,omitempty"`
	SessionToken    string  `json:"sessionToken,omitempty"`
}

func (s *LeasesService) Login(opts *LeaseLoginOptions) {
	log.Debugln("Requesting leased account credentials")

	params := &operations.PostLeasesAuthParams{}
	params.SetTimeout(20 * time.Second)
	res, err := ApiClient.PostLeasesAuth(params, nil)

	if err != nil {
		log.Fatal(err)
	}

	responsePayload := res.GetPayload()

	creds := leaseCreds(*responsePayload)
	s.loginWithCreds(&creds, opts)
}

func (s *LeasesService) LoginByID(leaseID string, opts *LeaseLoginOptions) {
	log.Debugln("Requesting leased account credentials")
	params := &operations.PostLeasesIDAuthParams{
		ID: leaseID,
	}
	params.SetTimeout(20 * time.Second)
	res, err := ApiClient.PostLeasesIDAuth(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	}

	responsePayload := res.GetPayload()

	creds := leaseCreds(*responsePayload)
	s.loginWithCreds(&creds, opts)
}

func (s *LeasesService) loginWithCreds(leaseCreds *leaseCreds, opts *LeaseLoginOptions) {
	if !(opts.OpenBrowser || opts.PrintCreds) {
		credsPath := filepath.Join(".aws", "credentials")
		log.Infoln("Adding credentials to " + credsPath + " using AWS CLI")
		s.Util.ConfigureAWSCLICredentials(leaseCreds.AccessKeyID,
			leaseCreds.SecretAccessKey,
			leaseCreds.SessionToken,
			opts.CliProfile)

	} else if opts.CliProfile != "default" {
		log.Infoln("Setting --profile has no effect when used with other flags.\n")
	}

	if opts.OpenBrowser {
		log.Infoln("Opening AWS Console in Web Browser")
		s.Util.OpenURL(leaseCreds.ConsoleURL)
	}

	if opts.PrintCreds {
		creds := fmt.Sprintf(constants.CredentialsExport,
			leaseCreds.AccessKeyID,
			leaseCreds.SecretAccessKey,
			leaseCreds.SessionToken)
		if _, err := Out.Write([]byte(creds)); err != nil {
			log.Fatalln("err: ", err)
		}
	}
}
