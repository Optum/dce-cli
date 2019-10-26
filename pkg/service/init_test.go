package service

import (
	"testing"

	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
	"github.com/Optum/dce-cli/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestLogObservation struct {
	observ.LevelLogger
	ended bool
}

func (l *TestLogObservation) Endf(format string, args ...interface{}) {
	l.ended = true
}

func (l *TestLogObservation) End(args ...interface{}) {
	l.ended = true
}

func (l *TestLogObservation) Endln(args ...interface{}) {
	l.ended = true
}

func TestInitializeDCE(t *testing.T) {

	t.Run("GIVEN a config file does not exists", func(t *testing.T) {
		mockPrompter := &mocks.Prompter{}
		mockFileSystemer := &mocks.FileSystemer{}
		spyLogger := &TestLogObservation{
			logrus.New(),
			false,
		}
		testObservation := &observ.ObservationContainer{
			Logger: spyLogger,
		}
		mockUtils := &utl.UtilContainer{
			Prompter:     mockPrompter,
			FileSystemer: mockFileSystemer,
			Observation:  testObservation,
		}

		service := InitService{
			Config:      nil,
			Util:        mockUtils,
			Observation: testObservation,
		}

		empty := ""
		defaultConfig := "defaultConfig"
		mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&empty)
		mockPrompter.On("PromptSelect", mock.Anything, mock.Anything).Return(&empty)
		mockFileSystemer.On("GetDefaultConfigFile", mock.Anything, mock.Anything).Return(defaultConfig)
		mockFileSystemer.On("WriteToYAMLFile", mock.Anything, mock.Anything)
		mockFileSystemer.On("IsExistingFile", mock.Anything).Return(false)

		t.Run("WHEN config path is not specified", func(t *testing.T) {
			service.InitializeDCE("")

			t.Run("THEN write config to default path", func(t *testing.T) {
				mockFileSystemer.AssertCalled(t, "WriteToYAMLFile", defaultConfig, mock.Anything)

				mockPrompter.AssertExpectations(t)
				mockFileSystemer.AssertExpectations(t)
			})
		})

		t.Run("WHEN config path is specified", func(t *testing.T) {
			specifiedPath := "somePath"
			service.InitializeDCE(specifiedPath)

			t.Run("THEN generate new config file at specified path", func(t *testing.T) {
				mockFileSystemer.AssertCalled(t, "WriteToYAMLFile", specifiedPath, mock.Anything)

				mockPrompter.AssertExpectations(t)
				mockFileSystemer.AssertExpectations(t)
			})
		})
	})

	t.Run("GIVEN a config file exists AND config path is not specified", func(t *testing.T) {
		// Arrange
		mockPrompter := &mocks.Prompter{}
		mockFileSystemer := &mocks.FileSystemer{}
		spyLogger := &TestLogObservation{
			logrus.New(),
			false,
		}
		testObservation := &observ.ObservationContainer{
			Logger: spyLogger,
		}
		mockUtils := &utl.UtilContainer{
			Prompter:     mockPrompter,
			FileSystemer: mockFileSystemer,
			Observation:  testObservation,
		}

		service := InitService{
			Config:      nil,
			Util:        mockUtils,
			Observation: testObservation,
		}

		doesntMatter := "doesn't matter"
		mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&doesntMatter)
		mockPrompter.On("PromptSelect", mock.Anything, mock.Anything).Return(&doesntMatter)
		mockFileSystemer.On("WriteToYAMLFile", mock.Anything, mock.Anything)
		defaultConfigPath := "defaultConfigPath"
		mockFileSystemer.On("GetDefaultConfigFile", mock.Anything, mock.Anything).Return(defaultConfigPath)
		mockFileSystemer.On("IsExistingFile", defaultConfigPath).Return(true)

		t.Run("WHEN command executes AND ask for confirmation", func(t *testing.T) {

			t.Run("AND confirmation is approved", func(t *testing.T) {
				yes := "yes"
				mockPrompter.On("PromptBasic", constants.PromptOverwiteConfig, mock.Anything).Return(&yes)
				service.InitializeDCE("")
				if !mockFileSystemer.AssertExpectations(t) || !mockPrompter.AssertExpectations(t) {
					t.Skip()
				}
				t.Run("THEN overwite existing config file", func(t *testing.T) {
					mockFileSystemer.AssertCalled(t, "WriteToYAMLFile", defaultConfigPath, mock.Anything)
				})
			})

			t.Run("AND confirmation is not approved", func(t *testing.T) {
				notYes := "not yes"
				mockPrompter.On("PromptBasic", constants.PromptOverwiteConfig, mock.Anything).Return(&notYes)
				service.InitializeDCE("")
				if !mockFileSystemer.AssertExpectations(t) || !mockPrompter.AssertExpectations(t) {
					t.Skip()
				}
				t.Run("THEN end process", func(t *testing.T) {
					require.True(t, spyLogger.ended, "Process not ended")
				})
			})

		})
	})
}
