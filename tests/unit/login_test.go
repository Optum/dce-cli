package unit

import (
	"testing"
	"time"

	"github.com/Optum/dce-cli/client/operations"
	"github.com/Optum/dce-cli/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLeaseLogin(t *testing.T) {
	emptyConfig := configs.Root{}
	t.Run("GIVEN lease id is valid and openBrowser is false", func(t *testing.T) {
		leaseID := "leaseID"
		openBrowser := false

		expectedAccessKeyID := "expectedAccessKeyID"
		expectedSecretAccessKey := "SecretAccessKey"
		expectedSessionToken := "expectedAccessKeyID"
		exprectedConsoleURL := "ConsoleURL"

		reqParams := &operations.PostLeasesIDAuthParams{
			ID: leaseID,
		}
		reqParams.SetTimeout(5 * time.Second)
		t.Run("WHEN LoginToLease", func(t *testing.T) {
			initMocks(emptyConfig)
			mockAPIer.On("PostLeasesIDAuth", reqParams, nil).Return(&operations.PostLeasesIDAuthCreated{
				Payload: &operations.PostLeasesIDAuthCreatedBody{
					AccessKeyID:     expectedAccessKeyID,
					SecretAccessKey: expectedSecretAccessKey,
					SessionToken:    expectedSessionToken,
					ConsoleURL:      exprectedConsoleURL,
				},
			}, nil)

			service.LoginToLease(leaseID, openBrowser)
			t.Run("THEN print credentials and don't open web browser", func(t *testing.T) {
				mockWeber.AssertExpectations(t)
				mockAPIer.AssertExpectations(t)
				mockWeber.AssertNotCalled(t, "OpenURL", mock.Anything)

				expectedOutput := "aws configure set aws_access_key_id " + expectedAccessKeyID +
					";aws configure set aws_secret_access_key " + expectedSecretAccessKey +
					";aws configure set aws_session_token " + expectedSessionToken
				assert.Equal(t, expectedOutput, spyLogger.Msg)
			})
		})
	})
	t.Run("GIVEN lease id is valid and openBrowser is true", func(t *testing.T) {
		leaseID := "leaseID"
		openBrowser := true

		expectedAccessKeyID := "expectedAccessKeyID"
		expectedSecretAccessKey := "SecretAccessKey"
		expectedSessionToken := "expectedAccessKeyID"
		exprectedConsoleURL := "ConsoleURL"

		reqParams := &operations.PostLeasesIDAuthParams{
			ID: leaseID,
		}
		reqParams.SetTimeout(5 * time.Second)
		t.Run("WHEN LoginToLease", func(t *testing.T) {
			initMocks(emptyConfig)
			mockAPIer.On("PostLeasesIDAuth", reqParams, nil).Return(&operations.PostLeasesIDAuthCreated{
				Payload: &operations.PostLeasesIDAuthCreatedBody{
					AccessKeyID:     expectedAccessKeyID,
					SecretAccessKey: expectedSecretAccessKey,
					SessionToken:    expectedSessionToken,
					ConsoleURL:      exprectedConsoleURL,
				},
			}, nil)
			mockWeber.On("OpenURL", exprectedConsoleURL)

			service.LoginToLease(leaseID, openBrowser)
			t.Run("THEN open browser to URL from api response and don't print credentials", func(t *testing.T) {
				mockWeber.AssertExpectations(t)
				mockAPIer.AssertExpectations(t)
				assert.Equal(t, "Opening AWS Console in Web Browser", spyLogger.Msg)
			})
		})
	})
}
