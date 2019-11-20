package unit

import (
	"os"
	"strings"
	"testing"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	svc "github.com/Optum/dce-cli/pkg/service"
	"github.com/stretchr/testify/mock"
)

var doesntMatter = "doesntmatter"

type mockFileInfo struct {
	os.FileInfo
}

func (m *mockFileInfo) Name() string { return doesntMatter }
func (m *mockFileInfo) IsDir() bool  { return true }

func TestDeployTFOverrides(t *testing.T) {
	config := configs.Root{}

	t.Run("GIVEN region is overridden", func(t *testing.T) {
		region := "someRegion"
		deployOverrides := svc.DeployOverrides{
			AWSRegion: region,
		}
		deployLocalPath := ""
		t.Run("WHEN Deploy", func(t *testing.T) {
			initMocks(config)
			mockMethods()
			service.Deploy(deployLocalPath, &deployOverrides)
			mockTerraformer.AssertExpectations(t)
			t.Run("THEN terraform apply is executed with the provided region", func(t *testing.T) {
				calledCountdown := 2
				mockTerraformer.AssertCalled(t, "Apply", mock.MatchedBy(func(args []string) bool {
					isMatch := false
					if len(args) == 3 {
						isMatch = (args[0] == `aws_region=`+region) &&
							strings.Contains(args[2], `namespace=`)
					}
					if len(args) == 1 {
						isMatch = strings.Contains(args[0], `namespace=`)
					}
					calledCountdown--
					return isMatch && calledCountdown == 0
				}))
			})
		})
	})
	t.Run("GIVEN tags are added", func(t *testing.T) {
		deployOverrides := svc.DeployOverrides{
			GlobalTags: []string{"a:b", "c:d"},
		}
		deployLocalPath := ""
		t.Run("WHEN Deploy", func(t *testing.T) {
			initMocks(config)
			mockMethods()

			service.Deploy(deployLocalPath, &deployOverrides)
			mockTerraformer.AssertExpectations(t)

			t.Run("THEN terraform apply is executed with the provided tags and default tags", func(t *testing.T) {
				calledCountdown := 2
				mockTerraformer.AssertCalled(t, "Apply", mock.MatchedBy(func(args []string) bool {
					isMatch := false
					if len(args) == 2 {
						isMatch = (args[0] == `global_tags={`+constants.GlobalTFTagDefaults+`,"a":"b","c":"d"}`) &&
							strings.Contains(args[1], `namespace=`)
					}
					if len(args) == 1 {
						isMatch = strings.Contains(args[0], `namespace=`)
					}
					calledCountdown--
					return isMatch && calledCountdown == 0
				}))
				mockTerraformer.AssertNumberOfCalls(t, "Apply", 2)
			})
		})
	})
	t.Run("GIVEN namespace is overridden", func(t *testing.T) {
		expectedNamspace := "expectedNamspace"
		deployOverrides := svc.DeployOverrides{
			Namespace: expectedNamspace,
		}
		deployLocalPath := ""
		t.Run("WHEN Deploy", func(t *testing.T) {
			initMocks(config)
			mockMethods()

			service.Deploy(deployLocalPath, &deployOverrides)
			mockTerraformer.AssertExpectations(t)
			t.Run("THEN terraform apply is executed with the provided namespace", func(t *testing.T) {
				calledCountdown := 2
				mockTerraformer.AssertCalled(t, "Apply", mock.MatchedBy(func(args []string) bool {
					isMatch := false
					if len(args) == 2 {
						isMatch = (args[0] == `global_tags={`+constants.GlobalTFTagDefaults+`}`) &&
							strings.Contains(args[1], `namespace=`+expectedNamspace)
					}
					if len(args) == 1 {
						isMatch = strings.Contains(args[0], `namespace=`)
					}
					calledCountdown--
					return isMatch && calledCountdown == 0
				}))
			})
		})
	})
	t.Run("GIVEN BudgetNotificationFromEmail is overridden", func(t *testing.T) {
		emailTemplate := "emailTemplate"
		deployOverrides := svc.DeployOverrides{
			BudgetNotificationFromEmail: emailTemplate,
		}
		deployLocalPath := ""
		t.Run("WHEN Deploy", func(t *testing.T) {
			initMocks(config)
			mockMethods()

			service.Deploy(deployLocalPath, &deployOverrides)
			mockTerraformer.AssertExpectations(t)
			t.Run("THEN terraform apply is executed with the provided email", func(t *testing.T) {
				calledCountdown := 2
				mockTerraformer.AssertCalled(t, "Apply", mock.MatchedBy(func(args []string) bool {
					isMatch := false
					if len(args) == 3 {
						isMatch = (args[2] == `budget_notification_from_email=`+emailTemplate) &&
							strings.Contains(args[1], `namespace=`)
					}
					if len(args) == 1 {
						isMatch = strings.Contains(args[0], `namespace=`)
					}
					calledCountdown--
					return isMatch && calledCountdown == 0
				}))

				mockTerraformer.AssertNumberOfCalls(t, "Apply", 2)
			})
		})
	})
	t.Run("GIVEN BudgetNotificationBCCEmails is overridden", func(t *testing.T) {
		emailTemplate := "emailTemplate"
		deployOverrides := svc.DeployOverrides{
			BudgetNotificationBCCEmails: []string{emailTemplate},
		}
		deployLocalPath := ""
		t.Run("WHEN Deploy", func(t *testing.T) {
			initMocks(config)
			mockMethods()

			service.Deploy(deployLocalPath, &deployOverrides)
			mockTerraformer.AssertExpectations(t)
			t.Run("THEN terraform apply is executed with the provided email and default tags", func(t *testing.T) {
				calledCountdown := 2
				mockTerraformer.AssertCalled(t, "Apply", mock.MatchedBy(func(args []string) bool {
					isMatch := false
					if len(args) == 3 {
						isMatch = (args[2] == `budget_notification_bcc_emails=["`+emailTemplate+`"]`) &&
							strings.Contains(args[1], `namespace=`)
					}
					if len(args) == 1 {
						isMatch = strings.Contains(args[0], `namespace=`)
					}
					calledCountdown--
					return isMatch && calledCountdown == 0
				}))
			})
		})
	})
	t.Run("GIVEN BudgetNotificationTemplateHTML is overridden", func(t *testing.T) {
		emailTemplate := "emailTemplate"
		deployOverrides := svc.DeployOverrides{
			BudgetNotificationTemplateHTML: emailTemplate,
		}
		deployLocalPath := ""
		t.Run("WHEN Deploy", func(t *testing.T) {
			initMocks(config)
			mockMethods()

			service.Deploy(deployLocalPath, &deployOverrides)
			mockTerraformer.AssertExpectations(t)
			t.Run("THEN terraform apply is executed with the provided email and default tags", func(t *testing.T) {
				calledCountdown := 2
				mockTerraformer.AssertCalled(t, "Apply", mock.MatchedBy(func(args []string) bool {
					isMatch := false
					if len(args) == 3 {
						isMatch = (args[2] == `budget_notification_template_html=`+emailTemplate) &&
							strings.Contains(args[1], `namespace=`)
					}
					if len(args) == 1 {
						isMatch = strings.Contains(args[0], `namespace=`)
					}
					calledCountdown--
					return isMatch && calledCountdown == 0
				}))
			})
		})
	})
	t.Run("GIVEN BudgetNotificationTemplateText is overridden", func(t *testing.T) {
		emailTemplate := "emailTemplate"
		deployOverrides := svc.DeployOverrides{
			BudgetNotificationTemplateText: emailTemplate,
		}
		deployLocalPath := ""
		t.Run("WHEN Deploy", func(t *testing.T) {
			initMocks(config)
			mockMethods()

			service.Deploy(deployLocalPath, &deployOverrides)
			mockTerraformer.AssertExpectations(t)
			t.Run("THEN terraform apply is executed with the provided email and default tags", func(t *testing.T) {
				calledCountdown := 2
				mockTerraformer.AssertCalled(t, "Apply", mock.MatchedBy(func(args []string) bool {
					isMatch := false
					if len(args) == 3 {
						isMatch = (args[2] == `budget_notification_template_text=`+emailTemplate) &&
							strings.Contains(args[1], `namespace=`)
					}
					if len(args) == 1 {
						isMatch = strings.Contains(args[0], `namespace=`)
					}
					calledCountdown--
					return isMatch && calledCountdown == 0
				}))
			})
		})
	})
	t.Run("GIVEN BudgetNotificationTemplateSubject is overridden", func(t *testing.T) {
		emailTemplate := "emailTemplate"
		deployOverrides := svc.DeployOverrides{
			BudgetNotificationTemplateSubject: emailTemplate,
		}
		deployLocalPath := ""
		t.Run("WHEN Deploy", func(t *testing.T) {
			initMocks(config)
			mockMethods()

			service.Deploy(deployLocalPath, &deployOverrides)
			mockTerraformer.AssertExpectations(t)
			t.Run("THEN terraform apply is executed with the provided email and default tags", func(t *testing.T) {
				calledCountdown := 2
				mockTerraformer.AssertCalled(t, "Apply", mock.MatchedBy(func(args []string) bool {
					isMatch := false
					if len(args) == 3 {
						isMatch = (args[2] == `budget_notification_template_subject=`+emailTemplate) &&
							strings.Contains(args[1], `namespace=`)
					}
					if len(args) == 1 {
						isMatch = strings.Contains(args[0], `namespace=`)
					}
					calledCountdown--
					return isMatch && calledCountdown == 0
				}))
			})
		})
	})
}

func mockMethods() {
	mockGithuber.On("DownloadGithubReleaseAsset", mock.Anything)
	mockFileSystemer.On("Unarchive", mock.Anything, mock.Anything)
	mockFileSystemer.On("MvToTempDir", mock.Anything, mock.Anything).Return(doesntMatter, doesntMatter)
	mockFileSystemer.On("RemoveAll", mock.Anything, mock.Anything)
	mockFileSystemer.On("Chdir", mock.Anything, mock.Anything)

	mockFileSystemer.On("ReadDir", mock.Anything).Return([]os.FileInfo{&mockFileInfo{}})
	mockFileSystemer.On("WriteFile", mock.Anything, mock.Anything)
	mockTerraformer.On("Init", mock.Anything)
	mockTerraformer.On("Apply", mock.Anything)
	mockTerraformer.On("GetOutput", mock.Anything).Return(doesntMatter)

	mockAwser.On("UploadDirectoryToS3", mock.Anything, mock.Anything, mock.Anything).Return([]string{}, []string{})
	mockAwser.On("UpdateLambdasFromS3Assets", mock.Anything, mock.Anything, mock.Anything)
}
