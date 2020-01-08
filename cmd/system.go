package cmd

import (
	"context"

	cfg "github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	svc "github.com/Optum/dce-cli/pkg/service"
	"github.com/spf13/cobra"
)

type ExitFunc func(code int)

var (
	dceRepoPath     string
	deployOverrides svc.DeployOverrides
	deployConfig    cfg.DeployConfig
)

func init() {
	deployOverrides = svc.DeployOverrides{}
	deployConfig = cfg.DeployConfig{}
	systemDeployCmd.Flags().StringVarP(&deployConfig.DeployLocalPath, "local", "l", "", "Path to a local DCE repo to deploy.")
	systemDeployCmd.Flags().BoolVarP(&deployConfig.Overwrite, "overwrite", "o", false, "Overwrite local backend state.")
	systemDeployCmd.Flags().BoolVar(&deployConfig.NoPrompt, "noprompt", false, "Skip prompting for resource creation.")
	systemDeployCmd.Flags().StringVarP(&deployOverrides.Namespace, "namespace", "n", "", "Set a custom terraform namespace (Optional)")
	systemDeployCmd.Flags().StringVarP(&deployOverrides.AWSRegion, "region", "r", "", "The aws region that DCE will be deployed to (Default: us-east-1)")
	systemDeployCmd.Flags().StringArrayVarP(&deployOverrides.GlobalTags, "tag", "t", []string{}, "Tags to be placed on all DCE resources. E.g. `dce system deploy --tag key1:value1 --tag key2:value2`")
	systemDeployCmd.Flags().StringVar(&deployOverrides.BudgetNotificationFromEmail, "budget-notification-from-email", "", "Email address from which budget notifications will be sent")
	systemDeployCmd.Flags().StringArrayVar(&deployOverrides.BudgetNotificationBCCEmails, "budget-notification-bcc-emails", []string{}, "Email address from which budget notifications will be sent")
	systemDeployCmd.Flags().StringVar(&deployOverrides.BudgetNotificationTemplateHTML, "budget-notification-template-html", "", "HTML template for budget notification emails")
	systemDeployCmd.Flags().StringVar(&deployOverrides.BudgetNotificationTemplateText, "budget-notification-template-text", "", "Text template for budget notification emails")
	systemDeployCmd.Flags().StringVar(&deployOverrides.BudgetNotificationTemplateSubject, "budget-notification-template-subject", "", "Subjet for budget notification emails")
	systemCmd.AddCommand(systemDeployCmd)

	RootCmd.AddCommand(systemCmd)

}

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Deploy and configure the DCE system",
}

var systemDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy DCE to a new master account",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)
		if err := Service.Deploy(ctx, &deployOverrides); err != nil {
			log.Fatalln(err)
		}
	},
}
