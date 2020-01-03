package configs

// Root contains config
type Root struct {
	API         API
	Region      *string
	GithubToken *string `yaml:"githubToken,omitempty"`
}

type API struct {
	Host     *string
	BasePath *string
	// Token for authenticating against the API
	// token is base64 encoded JSON, containing an STS token.
	Token *string `yaml:"token,omitempty"`
}

type DeployConfig struct {
	Overwrite       bool
	DeployLocalPath string
	NoPrompt        bool
}

var Regions = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"}
