package integration

import (
	"encoding/json"
	"github.com/Optum/dce-cli/client/operations"
	"github.com/Optum/dce-cli/mocks"
	"github.com/Optum/dce-cli/pkg/service"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type leaseCreateTestCase struct {
	// Args to pass to the dce-cli
	commandArgs []string

	// Request body we expect to bbe sent to the `POST /leases` endpoint
	expectedCreateLeaseRequest operations.PostLeasesBody

	// AccountID to return from the `POST /leases` endpoint
	mockAccountID string

	// JSON we expect the CLI operation to output to stdout
	expectedJSONOutput map[string]interface{}
}

func leasesCreateTest(t *testing.T, testCase *leaseCreateTestCase) {
	cli := NewCLITest(t)

	// Mock the DCE API
	api := &mocks.APIer{}

	// Should send a `POST /leases` request
	expectedLease := testCase.expectedCreateLeaseRequest
	api.On("PostLeases",
		mock.MatchedBy(func(params *operations.PostLeasesParams) bool {
			// Check that ExpiresOn is within ~5s of our expectation
			// to account for actual time passing
			assert.InDelta(t, expectedLease.ExpiresOn, params.Lease.ExpiresOn, 5)

			// Pull out ExpiresOn field,
			// so we can assert.Equals on the rest of the struct
			expectedLeaseCopy := expectedLease
			expectedLeaseCopy.ExpiresOn = 0
			actualLeaseCopy := params.Lease
			actualLeaseCopy.ExpiresOn = 0

			// Check the rest of the lease object
			assert.Equal(t, expectedLeaseCopy, actualLeaseCopy)
			return true
		}), nil).
		Return(func(params *operations.PostLeasesParams, wt runtime.ClientAuthInfoWriter) *operations.PostLeasesCreated {
			// Return a Lease, which looks like the requested lease
			// (but with an account ID)
			return &operations.PostLeasesCreated{
				Payload: &operations.PostLeasesCreatedBody{
					AccountID:                testCase.mockAccountID,
					BudgetAmount:             *params.Lease.BudgetAmount,
					BudgetCurrency:           *params.Lease.BudgetCurrency,
					BudgetNotificationEmails: params.Lease.BudgetNotificationEmails,
					ExpiresOn:                params.Lease.ExpiresOn,
					PrincipalID:              *params.Lease.PrincipalID,
				},
			}
		}, nil)

	// Mock the output writer
	out := &mocks.OutputWriter{}
	out.On("Write", mock.MatchedBy(func(out []byte) bool {
		var actualJSONOutput map[string]interface{}
		err := json.Unmarshal(out, &actualJSONOutput)
		assert.Nil(t, err, "output should be valid JSON")

		// Check that ExpiresOn is within ~5s of our expectation
		// to account for actual time passing
		_ = testCase.expectedJSONOutput
		expectedExpiresOn := testCase.expectedJSONOutput["expiresOn"].(int64)
		actualExpiresOn := actualJSONOutput["expiresOn"].(float64)
		assert.InDelta(t,
			expectedExpiresOn,
			actualExpiresOn,
			5,
		)

		// Pull out the expiresOn field,
		// so we can assert.Equals on the rest of JSON
		actualJSONOutput["expiresOn"] = int64(0)
		testCase.expectedJSONOutput["expiresOn"] = int64(0)

		assert.Equal(t, testCase.expectedJSONOutput, actualJSONOutput)

		return true
	})).Return(0, nil)

	// Mock the Authentication service (would pop open browser to auth user)
	authSvc := &mocks.Authenticater{}
	authSvc.On("Authenticate").Return(nil)

	cli.Inject(func(input *injectorInput) {
		service.ApiClient = api
		service.Out = out
		input.service.Authenticater = authSvc
	})

	// Run `dce leases create` command
	err := cli.Execute(testCase.commandArgs)
	require.Nil(t, err)

	api.AssertExpectations(t)
	out.AssertNumberOfCalls(t, "Write", 1)
}

func TestLeasesCreate(t *testing.T) {

	t.Run("with all required flags, only", func(t *testing.T) {
		leasesCreateTest(t, &leaseCreateTestCase{
			commandArgs: []string{
				"leases", "create",
				"-p", "test-user",
				"-b", "100", "-c", "USD",
				"-e", "a@example.com", "-e", "b@example.com",
			},
			expectedCreateLeaseRequest: operations.PostLeasesBody{
				BudgetAmount:             aws.Float64(100),
				BudgetCurrency:           aws.String("USD"),
				BudgetNotificationEmails: []string{"a@example.com", "b@example.com"},
				// Defaults to 7d
				ExpiresOn:   float64(time.Now().Add(time.Hour * 24 * 7).Unix()),
				PrincipalID: aws.String("test-user"),
			},
			mockAccountID: "123456789012",
			expectedJSONOutput: map[string]interface{}{
				"accountId":                "123456789012",
				"budgetAmount":             float64(100),
				"budgetCurrency":           "USD",
				"budgetNotificationEmails": []interface{}{"a@example.com", "b@example.com"},
				"principalId":              "test-user",
				// Defaults to 7d
				"expiresOn": time.Now().Add(time.Hour * 24 * 7).Unix(),
			},
		})
	})

	t.Run("with --expires-on flag", func(t *testing.T) {
		leasesCreateTest(t, &leaseCreateTestCase{
			commandArgs: []string{
				"leases", "create",
				"-p", "test-user",
				"-b", "100", "-c", "USD",
				"-e", "test@example.com",
				"--expires-on", "1h",
			},
			expectedCreateLeaseRequest: operations.PostLeasesBody{
				BudgetAmount:             aws.Float64(100),
				BudgetCurrency:           aws.String("USD"),
				BudgetNotificationEmails: []string{"test@example.com"},
				ExpiresOn:                float64(time.Now().Add(time.Hour).Unix()),
				PrincipalID:              aws.String("test-user"),
			},
			mockAccountID: "123456789012",
			expectedJSONOutput: map[string]interface{}{
				"accountId":                "123456789012",
				"budgetAmount":             float64(100),
				"budgetCurrency":           "USD",
				"budgetNotificationEmails": []interface{}{"test@example.com"},
				"principalId":              "test-user",
				"expiresOn":                time.Now().Add(time.Hour).Unix(),
			},
		})
	})

}
