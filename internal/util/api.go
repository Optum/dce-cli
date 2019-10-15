package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Optum/dce-cli/configs"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

type CreateAccountRequest struct {
	ID           string `json:"id"`
	AdminRoleArn string `json:"adminRoleArn"`
}

type LeaseRequest struct {
	PrincipalID              string   `json:"principalId"`
	AccountID                string   `json:"accountId"`
	BudgetAmount             float64  `json:"budgetAmount"`
	BudgetCurrency           string   `json:"budgetCurrency"`
	BudgetNotificationEmails []string `json:"budgetNotificationEmails"`
}

type ApiRequestInput struct {
	Method string
	Url    string
	Creds  *credentials.Credentials
	Region string
	Json   interface{}
}

type ApiResponse struct {
	http.Response
	json interface{}
}

type APIUtil struct {
	Config  *configs.Root
	Session *awsSession.Session
}

//Request sends sig4 signed requests to api
func (u *APIUtil) Request(input *ApiRequestInput) *ApiResponse {

	//TODO: Use a better pattern to set these
	input.Creds = credentials.NewStaticCredentials(
		*u.Config.API.Credentials.AwsAccessKeyID,
		*u.Config.API.Credentials.AwsSecretAccessKey,
		*u.Config.API.Credentials.AwsSessionToken,
	)

	if input.Region == "" {
		input.Region = "us-east-1"
	}

	// Create API request
	req, err := http.NewRequest(input.Method, input.Url, nil)

	// Sign our API request, using sigv4
	// See https://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html
	signer := sigv4.NewSigner(input.Creds)
	now := time.Now().Add(time.Duration(30) * time.Second)

	// If there's a json provided, add it when signing
	// Body does not matter if added before the signing, it will be overwritten
	if input.Json != nil {

		payload, err := json.Marshal(input.Json)
		fmt.Println("Marshalled Payload: ", string(payload))

		if err != nil {
			fmt.Println("Error marshaling json payload")
		}
		req.Header.Set("Content-Type", "application/json")
		_, err = signer.Sign(req, bytes.NewReader(payload),
			"execute-api", input.Region, now)
	} else {
		_, err = signer.Sign(req, nil, "execute-api",
			input.Region, now)
	}

	// Send the API requests
	// resp, err := http.DefaultClient.Do(req)
	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := httpClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	// Parse the JSON response
	apiResp := &ApiResponse{
		Response: *resp,
	}
	defer resp.Body.Close()
	var data interface{}

	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal([]byte(body), &data)
	if err == nil {
		apiResp.json = data
	}

	return apiResp
}
