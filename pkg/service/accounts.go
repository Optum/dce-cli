package service

import (
	"fmt"
	"io/ioutil"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/util"
	utl "github.com/Optum/dce-cli/internal/util"
)

const accountsPath = "/accounts"

type AccountsService struct {
	Config *configs.Root
	Util   *utl.UtilContainer
}

func (s *AccountsService) AddAccount(accountID, adminRoleARN string) {
	requestBody := &utl.CreateAccountRequest{
		ID:           accountID,
		AdminRoleArn: adminRoleARN,
	}

	accountsFullURL := *s.Config.API.BaseURL + accountsPath
	fmt.Println("Posting to: ", accountsFullURL)
	fmt.Println("Post body: ", requestBody)

	response := s.Util.Request(&util.ApiRequestInput{
		Method: "POST",
		Url:    accountsFullURL,
		Region: *s.Config.Region,
		Json:   requestBody,
	})

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("Response: ", response)
	fmt.Println("Response Body: ", body)
}

func (s *AccountsService) RemoveAccount(accountID string) {
	accountsFullURL := *s.Config.API.BaseURL + accountsPath + "/" + accountID
	fmt.Println("Posting to: ", accountsFullURL)

	response := s.Util.Request(&utl.ApiRequestInput{
		Method: "DELETE",
		Url:    accountsFullURL,
		Region: *s.Config.Region,
	})

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("Response: ", response)
	fmt.Println("Response Body: ", body)
}
