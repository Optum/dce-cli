package service

import (
	"github.com/Optum/dce-cli/configs"
	utl "github.com/Optum/dce-cli/internal/util"
)

// ServiceContainer is a service that injects its config and util into other services
type ServiceContainer struct {
	Config   *configs.Root
	Util     *utl.UtilContainer
	Deployer Deployer
}

// New returns a new ServiceContainer given config
func New(config *configs.Root, util *utl.UtilContainer) *ServiceContainer {
	return &ServiceContainer{
		Config:   config,
		Util:     util,
		Deployer: &DeployService{Config: config, Util: util},
	}
}

// Deployer deploys the DCE application
type Deployer interface {
	Deploy(namespace string)
}
