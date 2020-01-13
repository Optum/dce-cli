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
	Bin     *string
	Source  *string
	Backend map[string]interface{} `yaml:"backend,omitempty"`
}

type DeployConfig struct {
	UseCached       bool
	DeployLocalPath string
	BatchMode       bool
}

var Regions = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"}
