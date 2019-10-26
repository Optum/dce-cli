package service

import (
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	"github.com/Optum/dce-cli/internal/util"
	utl "github.com/Optum/dce-cli/internal/util"
)

const accountsPath = "/accounts"

type AccountsService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
}

func (s *AccountsService) AddAccount(accountID, adminRoleARN string) {
	requestBody := &utl.CreateAccountRequest{
		ID:           accountID,
		AdminRoleArn: adminRoleARN,
	}

	accountsFullURL := *s.Config.API.BaseURL + accountsPath
	response := s.Util.Request(&util.ApiRequestInput{
		Method: "POST",
		Url:    accountsFullURL,
		Region: *s.Config.Region,
		Json:   requestBody,
	})

	if response.StatusCode == 201 {
		Log.Println("Account added to DCE accounts pool")
	} else {
		Log.Println("DCE Responded with an error: ", response)
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
		Log.Println("Account removed from DCE accounts pool")
	} else {
		Log.Println("DCE Responded with an error: ", response)
	}
}
