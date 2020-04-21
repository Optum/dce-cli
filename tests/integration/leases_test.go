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

type cliInputs []string

type cliOutputs struct {
	httpRequest interface{}
	printedJson map[string]interface{}
	printedString string
}

func TestLeasesEnd(t *testing.T) {
	testLeaseID := "testLeaseId123"
	accountID := "accountID"
	principalID := "principalID"

	tests := []struct {
		name            string
		cliInputs       cliInputs
		leaseID         string
		principalID     string
		accountID       string
		expectedOutputs cliOutputs
		expErr          error

	}{
		{
			name: "only leaseID arg succeeds",
			leaseID: testLeaseID,
			cliInputs: []string{"leases", "end", testLeaseID},
			expectedOutputs: cliOutputs{
				printedString: "Lease ended",
			},
		},
		{
			name: "both accountID and principalID flags succeeds",
			accountID: accountID,
			principalID: principalID,
			cliInputs: []string{"leases", "end", "--account-id", accountID, "--principal-id", principalID},
			expectedOutputs: cliOutputs{
				printedString: "Lease ended",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cli := NewCLITest(t)

			api := &mocks.APIer{}

			if tt.leaseID != "" && (tt.accountID == "" && tt.principalID == "") {
				api.On("DeleteLeasesID",
					mock.MatchedBy(func(params *operations.DeleteLeasesIDParams) bool {
						return params.ID == tt.leaseID
					}), nil).Return( &operations.DeleteLeasesIDOK{
					Payload: &operations.DeleteLeasesIDOKBody{},
				}, nil)
			}

			if tt.leaseID == "" && (tt.accountID != "" && tt.principalID != "") {
				api.On("DeleteLeases",
					mock.MatchedBy(func(params *operations.DeleteLeasesParams) bool {
						return *params.Lease.AccountID == tt.accountID && *params.Lease.PrincipalID == tt.principalID
					}), nil).Return( &operations.DeleteLeasesOK{
					Payload: &operations.DeleteLeasesOKBody{},
				}, nil)
			}

			out := &mocks.OutputWriter{}
			out.On("Write", mock.MatchedBy(func(out []byte) bool {
				return string(out) == tt.expectedOutputs.printedString
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
			err := cli.Execute(tt.cliInputs)
			require.Nil(t, err)

			api.AssertExpectations(t)
			out.AssertNumberOfCalls(t, "Write", 1)
		})
	}
}


func TestLeasesCreate(t *testing.T) {

	tests := []struct {
		name string
		cliInputs
		cliOutputs
		// AccountID to return from the `POST /leases` endpoint
		mockAccountID string
	}{
		{
			"with --expires-on flag",
			[]string{
				"leases", "create",
				"-p", "test-user",
				"-b", "100", "-c", "USD",
				"-e", "test@example.com",
				"--expires-on", "1h",
			},
			cliOutputs{
				httpRequest: operations.PostLeasesBody{
					BudgetAmount:             aws.Float64(100),
					BudgetCurrency:           aws.String("USD"),
					BudgetNotificationEmails: []string{"test@example.com"},
					ExpiresOn:                float64(time.Now().Add(time.Hour).Unix()),
					PrincipalID:              aws.String("test-user"),
				},
				printedJson:   map[string]interface{}{
					"accountId":                "123456789012",
					"budgetAmount":             float64(100),
					"budgetCurrency":           "USD",
					"budgetNotificationEmails": []interface{}{"test@example.com"},
					"principalId":              "test-user",
					"expiresOn":                time.Now().Add(time.Hour).Unix(),
				},
			},
			"123456789012",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {

			cli := NewCLITest(t)

			// Mock the DCE API
			api := &mocks.APIer{}

			// Should send a `POST /leases` request
			expectedLease := testCase.cliOutputs.httpRequest.(operations.PostLeasesBody)
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
				_ = testCase.cliOutputs.printedJson
				expectedExpiresOn := testCase.cliOutputs.printedJson["expiresOn"].(int64)
				actualExpiresOn := actualJSONOutput["expiresOn"].(float64)
				assert.InDelta(t,
					expectedExpiresOn,
					actualExpiresOn,
					5,
				)

				// Pull out the expiresOn field,
				// so we can assert.Equals on the rest of JSON
				actualJSONOutput["expiresOn"] = int64(0)
				testCase.cliOutputs.printedJson["expiresOn"] = int64(0)

				assert.Equal(t, testCase.cliOutputs.printedJson, actualJSONOutput)

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
			err := cli.Execute(testCase.cliInputs)
			require.Nil(t, err)

			api.AssertExpectations(t)
			out.AssertNumberOfCalls(t, "Write", 1)
		})
	}
}
