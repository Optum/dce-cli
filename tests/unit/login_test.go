package unit

import (
	"fmt"
	service2 "github.com/Optum/dce-cli/pkg/service"
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
	opts          *service2.LeaseLoginOptions
	expectedOut   string
	isWeberCalled bool
}{
	{
		name:    "GIVEN no flags THEN do not print anything",
		leaseID: "doesntMatter",
		opts: &service2.LeaseLoginOptions{
			CliProfile: "default",
		},
	},
	{
		name:    "GIVEN openBrowser THEN Weber should open browser",
		leaseID: "doesntMatter",
		opts: &service2.LeaseLoginOptions{
			OpenBrowser: true,
			CliProfile:  "default",
		},
		expectedOut:   "Opening AWS Console in Web Browser",
		isWeberCalled: true,
	},
	{
		name:    "GIVEN printCreds THEN Weber should open browser",
		leaseID: "doesntMatter",
		opts: &service2.LeaseLoginOptions{
			PrintCreds: true,
			CliProfile: "default",
		},
		expectedOut: credsOutput,
	},
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
			reqParams.SetTimeout(20 * time.Second)
			mockAPIer.On("PostLeasesIDAuth", reqParams, nil).Return(&operations.PostLeasesIDAuthCreated{
				Payload: &operations.PostLeasesIDAuthCreatedBody{
					AccessKeyID:     expectedAccessKeyID,
					SecretAccessKey: expectedSecretAccessKey,
					SessionToken:    expectedSessionToken,
					ConsoleURL:      expectedConsoleURL,
				},
			}, nil)
			if !(tc.opts.OpenBrowser || tc.opts.PrintCreds) {
				mockAwser.On("ConfigureAWSCLICredentials", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
			if tc.isWeberCalled {
				mockWeber.On("OpenURL", expectedConsoleURL)
			}
			// Act
			service.LoginByID(tc.leaseID, tc.opts)

			// Assert
			mockWeber.AssertExpectations(t)
			mockAPIer.AssertExpectations(t)

			assert.Contains(t, spyLogger.Msg, tc.expectedOut)
		})
	}
}

func TestLeaseLoginNoID(t *testing.T) {
	initMocks(configs.Root{})

	// Mock the `POST /leases/auth` endpoint
	reqParams := &operations.PostLeasesAuthParams{}
	reqParams.SetTimeout(20 * time.Second)
	mockAPIer.On("PostLeasesAuth", reqParams, nil).
		Return(&operations.PostLeasesAuthCreated{
			Payload: &operations.PostLeasesAuthCreatedBody{
				AccessKeyID:     "access-key-id",
				SecretAccessKey: "secret-access-key",
				SessionToken:    "session-token",
				ConsoleURL:      "console-url",
			},
		}, nil)

	// Weber.OpenURL() should be called with the
	// ConsoleURL returned by the API
	mockWeber.On("OpenURL", "console-url")

	// Run the login command
	service.Login(&service2.LeaseLoginOptions{
		OpenBrowser: true,
	})

	// Check that we called Weber.OpenURL()
	mockWeber.AssertExpectations(t)
}
