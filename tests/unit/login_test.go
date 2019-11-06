package unit

import (
	"testing"
	"time"

	"github.com/Optum/dce-cli/client/operations"
	"github.com/Optum/dce-cli/configs"
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
			mockWeber.AssertNotCalled(t, "OpenURL", exprectedConsoleURL)

			service.LoginToLease(leaseID, openBrowser)
			t.Run("THEN don't open web browser", func(t *testing.T) {
				mockWeber.AssertExpectations(t)
				mockAPIer.AssertExpectations(t)
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
			t.Run("THEN open browser to URL from api response", func(t *testing.T) {
				mockWeber.AssertExpectations(t)
				mockAPIer.AssertExpectations(t)
			})
		})
	})
}
