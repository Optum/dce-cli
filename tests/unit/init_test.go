package unit

import (
	"testing"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitializeDCE(t *testing.T) {
	emptyConfig := configs.Root{}

	t.Run("GIVEN config path is specified", func(t *testing.T) {
		providedConfigPath := "providedConfigPath"

		t.Run("WHEN InitializeDCE AND asks for confirmation AND confirmation is approved", func(t *testing.T) {
			initMocks(emptyConfig)

			doesntMatter := "doesn't matter"
			mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&doesntMatter)
			mockPrompter.On("PromptSelect", mock.Anything, mock.Anything).Return(&doesntMatter)
			mockFileSystemer.On("WriteToYAMLFile", mock.Anything, mock.Anything)

			yes := "yes"
			mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&yes)

			// Act
			service.InitializeDCE(providedConfigPath)

			mockFileSystemer.AssertNotCalled(t, "GetDefaultConfigFile")
			if !mockFileSystemer.AssertExpectations(t) || !mockPrompter.AssertExpectations(t) {
				t.Skip()
			}
			t.Run("THEN write to specified config path", func(t *testing.T) {
				mockFileSystemer.AssertCalled(t, "WriteToYAMLFile", providedConfigPath, mock.Anything)
			})
		})

		t.Run("WHEN InitializeDCE AND asks for confirmation AND confirmation is not approved", func(t *testing.T) {
			initMocks(emptyConfig)

			doesntMatter := "doesn't matter"
			mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&doesntMatter)
			mockPrompter.On("PromptSelect", mock.Anything, mock.Anything).Return(&doesntMatter)
			mockFileSystemer.On("WriteToYAMLFile", mock.Anything, mock.Anything)

			notYes := "not yes"
			mockPrompter.On("PromptBasic", constants.PromptChangeConfigConfirmation, mock.Anything).Return(&notYes)

			// Act
			service.InitializeDCE(providedConfigPath)

			if !mockFileSystemer.AssertExpectations(t) || !mockPrompter.AssertExpectations(t) {
				t.Skip()
			}
			t.Run("THEN end process", func(t *testing.T) {
				assert.True(t, spyLogger.Ended, "Process not ended")
			})
		})
	})

	t.Run("GIVEN config path is not specified", func(t *testing.T) {
		t.Run("WHEN InitializeDCE AND ask for confirmation", func(t *testing.T) {
			initMocks(emptyConfig)

			doesntMatter := "doesn't matter"
			mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&doesntMatter)
			mockPrompter.On("PromptSelect", mock.Anything, mock.Anything).Return(&doesntMatter)
			mockFileSystemer.On("WriteToYAMLFile", mock.Anything, mock.Anything)
			defaultConfigPath := "defaultConfigPath"
			mockFileSystemer.On("GetDefaultConfigFile", mock.Anything, mock.Anything).Return(defaultConfigPath)

			t.Run("AND confirmation is approved", func(t *testing.T) {
				yes := "yes"
				mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&yes)
				service.InitializeDCE("")
				if !mockFileSystemer.AssertExpectations(t) || !mockPrompter.AssertExpectations(t) {
					t.Skip()
				}
				t.Run("THEN write to default config path", func(t *testing.T) {
					mockFileSystemer.AssertCalled(t, "WriteToYAMLFile", defaultConfigPath, mock.Anything)
				})
			})

			t.Run("AND confirmation is not approved", func(t *testing.T) {
				notYes := "not yes"
				mockPrompter.On("PromptBasic", constants.PromptChangeConfigConfirmation, mock.Anything).Return(&notYes)
				service.InitializeDCE("")
				if !mockFileSystemer.AssertExpectations(t) || !mockPrompter.AssertExpectations(t) {
					t.Skip()
				}
				t.Run("THEN end process", func(t *testing.T) {
					assert.True(t, spyLogger.Ended, "Process not ended")
				})
			})

		})
	})
}
