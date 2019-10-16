package configs

// Root contains config
type Root struct {
	System struct {
		Auth struct {
			LoginURL *string
		}
		MasterAccount struct {
			Credentials struct {
				AwsAccessKeyID     *string
				AwsSecretAccessKey *string
			}
		}
	}
	API struct {
		BaseURL     *string
		Credentials struct {
			AwsAccessKeyID     *string
			AwsSecretAccessKey *string
			AwsSessionToken    *string
		}
	}
	Region      *string
	GithubToken *string
}

var Regions = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"}
