package integration

import (
	"fmt"
	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/ptr"
	"testing"
)

func TestAuthCommand(t *testing.T) {

	t.Run("GIVEN auth command is run", func(t *testing.T) {

		t.Run("THEN API token should be saved to config file", func(t *testing.T) {
			cli := NewCLITest(t)

			// Setup a basic config file
			confFile := writeTempConfig(t, &configs.Root{
				API:         configs.API{
					Host:     ptr.String("dce.example.com"),
					BasePath: ptr.String("/api"),
				},
			})

			// Mock Weber.OpenURL()
			mockWeber := &mocks.Weber{}
			cli.Inject(func (input *injectorInput) {
				input.service.Util.Weber = mockWeber
			})
			mockWeber.On("OpenURL", "https://dce.example.com/api/auth")

			// Enter the API Token (IRL, would be provided by the web page)
			cli.AnswerBasic("Enter API Token: ", "my-api-token")

			// Run DCE auth
			err := cli.Execute([]string{"auth", "--config", confFile})
			require.Nil(t, err)

			cli.AssertAllPrompts()
			mockWeber.AssertExpectations(t)

			output := cli.Output()
			require.Contains(t, output, "Opening web browser. Please Login and copy/paste the provided token into this terminal.")
			require.Contains(t, output, fmt.Sprintf("Saving API Token to %s", confFile))
		})

		t.Run("AND no API configuration is set", func(t *testing.T) {

			t.Run("THEN auth command should fail", func(t *testing.T) {
				cli := NewCLITest(t)

				// Setup a config file, missing API info
				confFile := writeTempConfig(t, &configs.Root{
					API:         configs.API{},
				})

				// Run DCE auth
				err := cli.Execute([]string{"auth", "--config", confFile})
				require.NotNil(t, err)

				require.Equal(t,
					"Unable to authenticate against DCE API: missing API configuration",
					err.Error(),
				)
			})

		})
	})
}

