package service

import (
	"time"

	"github.com/Optum/dce-cli/client/operations"
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
)

const accountsPath = "/accounts"

type AccountsService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
}

func (s *AccountsService) AddAccount(accountID, adminRoleARN string) {
	client := s.Util.SwaggerAPIClient
	params := &operations.PostAccountsParams{
		Account: operations.PostAccountsBody{
			ID:           &accountID,
			AdminRoleArn: &adminRoleARN,
		},
	}
	params.SetTimeout(5 * time.Second)
	_, err := client.PostAccounts(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	}

}

func (s *AccountsService) RemoveAccount(accountID string) {
	accountsFullURL := *s.Config.API.BaseURL + accountsPath + "/" + accountID
	response := s.Util.Request(&utl.ApiRequestInput{
		Method: "DELETE",
		Url:    accountsFullURL,
		Region: *s.Config.Region,
	})

	if response.StatusCode == 204 {
		log.Println("Account removed from DCE accounts pool")
	} else {
		log.Println("DCE Responded with an error: ", response)
	}
}
