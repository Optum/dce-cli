package configs

// Config contains config
type Config struct {
	System struct {
		Auth struct {
			LoginURL *string
		}
		MasterAccount struct {
			Credentials struct {
				AwsAccessKeyID     *string `yaml:"AWS_ACCESS_KEY_ID"` // TODO: Figure out why these tags aren't working
				AwsSecretAccessKey *string `yaml:"AWS_SECRET_ACCESS_KEY"`
				AwsSessionToken    *string `yaml:"AWS_SESSION_TOKEN"`
			}
		}
	}
	API struct {
		BaseURL *string
		Region  *string
	}
}
