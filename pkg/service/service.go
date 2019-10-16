package service

import (
	"github.com/Optum/dce-cli/configs"
	utl "github.com/Optum/dce-cli/internal/util"
)

// ServiceContainer is a service that injects its config and util into other services
type ServiceContainer struct {
	Config *configs.Root
	Util   *utl.UtilContainer
	Deployer
	Accounter
	Leaser
	Initer
	Authenticater
}

// New returns a new ServiceContainer given config
func New(config *configs.Root, util *utl.UtilContainer) *ServiceContainer {
	return &ServiceContainer{
		Config:        config,
		Util:          util,
		Deployer:      &DeployService{Config: config, Util: util},
		Accounter:     &AccountsService{Config: config, Util: util},
		Leaser:        &LeasesService{Config: config, Util: util},
		Initer:        &InitService{},
		Authenticater: &AuthService{Config: config},
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
