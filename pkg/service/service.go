package service

import (
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
}

var log observ.Logger

// New returns a new ServiceContainer given config
func New(config *configs.Root, observation *observ.ObservationContainer, util *utl.UtilContainer) *ServiceContainer {

	log = observation.Logger

	return &ServiceContainer{
		Config:        config,
		Observation:   observation,
		Util:          util,
		Deployer:      &DeployService{Config: config, Util: util},
		Accounter:     &AccountsService{Config: config, Util: util},
		Leaser:        &LeasesService{Config: config, Util: util},
		Initer:        &InitService{Config: config, Util: util},
		Authenticater: &AuthService{Config: config, Util: util},
	}
}

// Deployer deploys the DCE application
type Deployer interface {
	Deploy(namespace string)
}

type Accounter interface {
	AddAccount(accountID, adminRoleARN string)
	RemoveAccount(accountID string)
}

type Leaser interface {
	CreateLease(principleID string, budgetAmount float64, budgetCurrency string, email []string)
	EndLease(accountID, principleID string)
	LoginToLease(loginAcctID, loginLeaseID string, loginOpenBrowser bool)
}

type Initer interface {
	InitializeDCE(cfgFile string)
}

type Authenticater interface {
	Authenticate(authUrl string)
}
