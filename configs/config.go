package configs

import "os"

// Root contains config
type Root struct {
	API       API
	Region    *string
	Terraform Terraform
}

type API struct {
	Host     *string
	BasePath *string
	// Token for authenticating against the API
	// token is base64 encoded JSON, containing an STS token.
	Token *string `yaml:"token,omitempty"`
}

// Terraform contains configuration for the underlying terraform
// command used to provision the DCE infrastructure.
type Terraform struct {
	Bin            *string
	Source         *string
	TFInitOptions  *string `yaml:"initOptions,omitempty"`
	TFApplyOptions *string `yaml:"applyOptions,omitempty"`
}

// DeployConfig holds configuration values for the `system deploy`
// command
type DeployConfig struct {
	// UseCached, if true, tells DCE to use files already in the
	// `~/.dce/.cache` folder
	UseCached bool
	// DeployLocalPath, if set, specifies a path from which to pull
	// local resources
	DeployLocalPath string
	// BatchMode, if enabled, forces DCE to run non-interactively
	// and supplies Terraform with -auto-approve and input=false
	BatchMode bool
	// TFInitOptions are options passed along to `terraform init`
	TFInitOptions string
	// TFApplyOptions are options passed along to `terraform apply`
	TFApplyOptions string
	// SaveTFOptions, if yes, will save the provided terraform options to the config file.
	SaveTFOptions bool
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
