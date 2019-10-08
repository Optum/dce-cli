package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

type CreateAccountRequest struct {
	ID           string `json:"id"`
	AdminRoleArn string `json:"adminRoleArn"`
}

type CreateLeaseRequest struct {
	PrincipalID              string   `json:"principalId"`
	AccountID                string   `json:"accountId"`
	BudgetAmount             float64  `json:"budgetAmount"`
	BudgetCurrency           string   `json:"budgetCurrency"`
	BudgetNotificationEmails []string `json:"budgetNotificationEmails"`
}

type GetLeaseRequest struct {
	PrincipalID string `json:"principalId"`
	AccountID   string `json:"accountId"`
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

//Request sends sig4 signed requests to api
func Request(input *ApiRequestInput) *ApiResponse {
	// Set defaults
	if input.Creds == nil {
		input.Creds = credentials.NewChainCredentials([]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{Filename: "", Profile: ""},
		})
	}
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
