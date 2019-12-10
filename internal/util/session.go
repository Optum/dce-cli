package util

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

func NewAWSSession(token *string) (*session.Session, error) {
	// Setup the AWS credentials provider chain.
	// First, we'll check for credentials in the
	// dce.yaml's `api.token` config.
	// then we'll use AWS's standard chain (env vars, ~/aws/credentials file)
	creds := credentials.NewChainCredentials([]credentials.Provider{
		NewAPITokenProvider(token),
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
	})
	return session.NewSession(&aws.Config{
		Credentials: creds,
	})
}

// APITokenProvider is a custom AWS Credentials provider
// which uses a base64 encoded token containing a STS credentials as JSON
// See https://docs.aws.amazon.com/sdk-for-go/api/aws/credentials/#hdr-Custom_Provider
//
// Using a custom Provider does a few things for us:
// - Allows chaining credentials, so we can fall back to env vars, creds fil
// - The `Retrieve` method is cached by the client, so we don't need to re-parse our API token at every call
// - Provides a mechanism for handling expired creds
type APITokenProvider struct {
	token      *string
	expiration int64
}

const APITokenProviderName = "APITokenProvider"

func NewAPITokenProvider(token *string) credentials.Provider {
	return &APITokenProvider{
		token: token,
	}
}

func (t *APITokenProvider) Retrieve() (credentials.Value, error) {
	if t.token == nil {
		return credentials.Value{ProviderName: APITokenProviderName},
			errors.New("no API token is configured")
	}

	stsTokenJSON, err := base64.StdEncoding.DecodeString(*t.token)
	if err != nil {
		return credentials.Value{ProviderName: APITokenProviderName},
			errors.New("failed to decode token")
	}

	// Unmarshal the STS Token JSON
	var tokenValue APITokenValue
	err = json.Unmarshal(stsTokenJSON, &tokenValue)
	if err != nil {
		return credentials.Value{ProviderName: APITokenProviderName},
			errors.New("decoded token contains invalid JSON")
	}

	// Remember the tokens `expired` time,
	// so we can implement `IsExpired`
	t.expiration = tokenValue.Expiration

	if t.IsExpired() {
		return credentials.Value{ProviderName: APITokenProviderName},
			errors.New("token is expired")
	}

	return credentials.Value{
		AccessKeyID:     tokenValue.AccessKeyID,
		SecretAccessKey: tokenValue.SecretAccessKey,
		SessionToken:    tokenValue.SessionToken,
		ProviderName:    APITokenProviderName,
	}, nil
}

func (t *APITokenProvider) IsExpired() bool {
	return time.Now().Unix() > t.expiration
}

type APITokenValue struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	Expiration      int64 `json:"expireTime"`
}
