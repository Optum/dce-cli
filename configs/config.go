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

type Terraform struct {
	Bin            *string
	Source         *string
	TFInitOptions  *string `yaml:"tfInitOptions,omitempty"`
	TFApplyOptions *string `yaml:"tfApplyOptions,omitempty"`
}

type DeployConfig struct {
	UseCached       bool
	DeployLocalPath string
	BatchMode       bool
	TFInitOptions   string
	TFApplyOptions  string
}

var Regions = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"}

// Coalesce returns the first non-empty vluae, but takes into account a loading order,
// which is CLI args > config file > environment variables > some default
func Coalesce(arg *string, config *string, envvar *string, def *string) *string {

	if arg != nil && len(*arg) > 0 {
		return arg
	}

	if config != nil && len(*config) > 0 {
		return config
	}

	if envvar != nil {
		envval, ok := os.LookupEnv(*envvar)

		if ok && len(envval) > 0 {
			return &envval
		}
	}

	return def
}
