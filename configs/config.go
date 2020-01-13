package configs

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
