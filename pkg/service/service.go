package service

import (
	"encoding/json"

	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
)

// ServiceContainer is a service that injects its config and util into other services
type ServiceContainer struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
	Deployer
	Accounter
	Leaser
	Initer
	Authenticater
	Usager
}

var log observ.Logger
var apiClient utl.APIer

// New returns a new ServiceContainer given config
func New(config *configs.Root, observation *observ.ObservationContainer, util *utl.UtilContainer) *ServiceContainer {

	log = observation.Logger
	apiClient = util.APIer

	serviceContainer := ServiceContainer{
		Config:        config,
		Observation:   observation,
		Util:          util,
		Deployer:      &DeployService{Config: config, Util: util},
		Accounter:     &AccountsService{Config: config, Util: util},
		Leaser:        &LeasesService{Config: config, Util: util},
		Initer:        &InitService{Config: config, Util: util},
		Authenticater: &AuthService{Config: config, Util: util},
		Usager:        &UsageService{Config: config, Util: util},
	}

	return &serviceContainer
}

type DeployOverrides struct {
	AWSRegion                         string
	GlobalTags                        []string
	Namespace                         string
	BudgetNotificationFromEmail       string
	BudgetNotificationBCCEmails       []string
	BudgetNotificationTemplateHTML    string
	BudgetNotificationTemplateText    string
	BudgetNotificationTemplateSubject string
}
type Deployer interface {
	Deploy(deployLocal string, overrides *DeployOverrides)
}

type Usager interface {
	GetUsage(startDate, endDate float64)
}

type Accounter interface {
	AddAccount(accountID, adminRoleARN string)
	RemoveAccount(accountID string)
	GetAccount(accountID string)
	ListAccounts()
}

type Leaser interface {
	CreateLease(principalID string, budgetAmount float64, budgetCurrency string, email []string, expiresOn string)
	EndLease(accountID, principalID string)
	LoginToLease(leaseID, profile string, loginOpenBrowser, loginPrintCreds bool)
	ListLeases(acctID, principalID, nextAcctID, nextPrincipalID, leaseStatus string, pagLimit int64)
	GetLease(leaseID string)
}

type Initer interface {
	InitializeDCE()
}

type Authenticater interface {
	Authenticate() error
}

type ResponseWithPayload interface {
	GetPayload() interface{}
}

func printResponsePayload(res ResponseWithPayload) {
	jsonPayload, err := json.MarshalIndent(res.GetPayload(), "", "\t")
	if err != nil {
		log.Fatalln("err: ", err)
	}
	log.Infoln(string(jsonPayload))
}
