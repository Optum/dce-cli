package configs

// Root contains config
type Root struct {
	System struct {
		Auth struct {
			LoginURL *string
		}
	}
	API struct {
		Host     *string
		BasePath *string
	}
	Region      *string
	GithubToken *string
}

var Regions = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"}
