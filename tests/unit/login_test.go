package unit

import (
	"fmt"
	"testing"
	"time"

	"github.com/Optum/dce-cli/client/operations"
	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var expectedAccessKeyID = "expectedAccessKeyID"
var expectedSecretAccessKey = "SecretAccessKey"
var expectedSessionToken = "expectedAccessKeyID"
var expectedConsoleURL = "ConsoleURL"
var credsOutput = fmt.Sprintf(constants.CredentialsExport,
	expectedAccessKeyID,
	expectedSecretAccessKey,
	expectedSessionToken)

var testCases = []struct {
	name          string
	leaseID       string
	openBrowser   bool
	printCreds    bool
	profile       string
	expectedOut   string
	isWeberCalled bool
}{
	{"GIVEN no flags THEN do not print anything", "doesntMatter", false, false, "default", "", false},
	{"GIVEN openBrowser THEN Weber should open browser", "doesntMatter", true, false, "default", "Opening AWS Console in Web Browser", true},
	{"GIVEN printCreds THEN Weber should open browser", "doesntMatter", false, true, "default", credsOutput, false},
}

func TestLeaseLoginGivenFlags(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			emptyConfig := configs.Root{}
			initMocks(emptyConfig)
			reqParams := &operations.PostLeasesIDAuthParams{
				ID: tc.leaseID,
			}
			reqParams.SetTimeout(5 * time.Second)
			mockAPIer.On("PostLeasesIDAuth", reqParams, nil).Return(&operations.PostLeasesIDAuthCreated{
				Payload: &operations.PostLeasesIDAuthCreatedBody{
					AccessKeyID:     expectedAccessKeyID,
					SecretAccessKey: expectedSecretAccessKey,
					SessionToken:    expectedSessionToken,
					ConsoleURL:      expectedConsoleURL,
				},
			}, nil)
			if !(tc.openBrowser || tc.printCreds) {
				mockAwser.On("ConfigureAWSCLICredentials", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
			if tc.isWeberCalled {
				mockWeber.On("OpenURL", expectedConsoleURL)
			}
			// Act
			service.LoginToLease(tc.leaseID, tc.profile, tc.openBrowser, tc.printCreds)

			// Assert
			mockWeber.AssertExpectations(t)
			mockAPIer.AssertExpectations(t)

			assert.Contains(t, spyLogger.Msg, tc.expectedOut)
		})
	}
}
