package configs

// Root contains config
type Root struct {
	System struct {
		Auth struct {
			LoginURL *string `yaml:"loginURL"`
		} `yaml:"auth"`
		MasterAccount struct {
			Credentials struct {
				AwsAccessKeyID     *string `yaml:"AWS_ACCESS_KEY_ID"` // TODO: Figure out why these tags aren't working
				AwsSecretAccessKey *string `yaml:"AWS_SECRET_ACCESS_KEY"`
				AwsSessionToken    *string `yaml:"AWS_SESSION_TOKEN"`
			} `yaml:"credentials"`
		} `yaml:"masterAccount"`
	} `yaml:"system"`
	API struct {
		BaseURL *string `yaml:"baseURL"`
		Region  *string `yaml:"region"`
	}
}

var Regions = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"}
