package service

import (
	"encoding/json"
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
	params := &operations.PostAccountsParams{
		Account: operations.PostAccountsBody{
			ID:           &accountID,
			AdminRoleArn: &adminRoleARN,
		},
	}
	params.SetTimeout(5 * time.Second)
	_, err := apiClient.PostAccounts(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	} else {
		log.Infoln("Account added to DCE accounts pool")
	}
}

func (s *AccountsService) RemoveAccount(accountID string) {
	params := &operations.DeleteAccountsIDParams{
		ID: accountID,
	}
	params.SetTimeout(5 * time.Second)
	_, err := apiClient.DeleteAccountsID(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	} else {
		log.Infoln("Account removed from DCE accounts pool")
	}
}

func (s *AccountsService) GetAccount(accountID string) {
	params := &operations.GetAccountsIDParams{
		ID: accountID,
	}
	params.SetTimeout(5 * time.Second)
	res, err := apiClient.GetAccountsID(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	} else {
		log.Infoln(res)
	}
}

func (s *AccountsService) ListAccounts() {
	params := &operations.GetAccountsParams{}
	params.SetTimeout(5 * time.Second)
	res, err := apiClient.GetAccounts(params, nil)
	if err != nil {
		log.Fatalln("err: ", err)
	} else {
		jsonPayload, err := json.Marshal(res.GetPayload())
		if err != nil {
			log.Fatalln("err: ", err)
		}
		log.Infoln(string(jsonPayload))
	}
}
