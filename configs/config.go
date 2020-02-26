package configs

import "os"

// Root contains config
type Root struct {
	API       API
	Region    *string
	Deploy    Deploy `yaml:"deploy,omitempty"`
	Terraform Terraform
}

type API struct {
	Host     *string
	BasePath *string
	// Token for authenticating against the API
	// token is base64 encoded JSON, containing an STS token.
	Token *string `yaml:"token,omitempty"`
}

type Deploy struct {
	// Path to the DCE repo to deploy
	// May be a local file path (eg. /path/to/dce)
	// or a github repo (eg. github.com/Optum/dce)
	Location *string `yaml:"location,omitempty"`
	// Version of DCE to deploy, eg 0.12.3
	Version *string `yaml:"version,omitempty"`
	// Deployment logs will be written to this location
	LogFile *string `yaml:"logFile,omitempty"`
	// AWS Region in which to deploy DCE
	AWSRegion *string `yaml:"region,omitempty"`
	// Namespace used as naming suffix for AWS resources
	Namespace                   *string `yaml:"namespace,omitempty"`
	BudgetNotificationFromEmail *string `yaml:"budgetNotificationFromEmail,omitempty"`
}

// Terraform contains configuration for the underlying terraform
// command used to provision the DCE infrastructure.
type Terraform struct {
	Bin            *string
	Source         *string // URL from which the Terraform release was downloaded
	TFInitOptions  *string `yaml:"initOptions,omitempty"`
	TFApplyOptions *string `yaml:"applyOptions,omitempty"`
}

var Regions = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"}

// Coalesce returns the first non-empty vluae, but takes into account a loading order,
// which is CLI args > environment variables > config file > some default
func Coalesce(arg *string, config *string, envvar *string, def *string) *string {

	if arg != nil && len(*arg) > 0 {
		return arg
	}

	if envvar != nil {
		envval, ok := os.LookupEnv(*envvar)

		if ok && len(envval) > 0 {
			return &envval
		}
	}

	if config != nil && len(*config) > 0 {
		return config
	}

	return def
}
