package cmd

import (
	"context"

	cfg "github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	svc "github.com/Optum/dce-cli/pkg/service"
	"github.com/spf13/cobra"
)

var (
	dceRepoPath     string
	deployOverrides svc.DeployOverrides
	deployConfig    cfg.DeployConfig
)

const (
	DCETFInitOptionsEnvVar  string = "DCE_TF_INIT_OPTIONS"
	DCETFApplyOptionsEnvVar string = "DCE_TF_APPLY_OPTIONS"
	Empty                   string = ""
)

func init() {
	deployOverrides = svc.DeployOverrides{}
	deployConfig = cfg.DeployConfig{}
	systemDeployCmd.Flags().StringVarP(&deployConfig.DeployLocalPath, "local", "l", "", "Path to a local DCE repo to deploy.")
	systemDeployCmd.Flags().BoolVarP(&deployConfig.UseCached, "use-cached", "c", true, "Overwrite local backend state.")
	systemDeployCmd.Flags().BoolVarP(&deployConfig.BatchMode, "batch-mode", "b", false, "Skip prompting for resource creation.")
	systemDeployCmd.Flags().StringVar(&deployConfig.TFInitOptions, "tf-init-options", "", "Options to pass to the underlying `tf init` command.")
	systemDeployCmd.Flags().StringVar(&deployConfig.TFApplyOptions, "tf-apply-options", "", "Options to pass to the underlying `tf apply` command.")
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
	RunE: func(cmd *cobra.Command, args []string) error {
		// Before running the command, resolve the configuration by using
		// the default load order.
		deployConfig.TFInitOptions = *cfg.Coalesce(&deployConfig.TFInitOptions, Config.Terraform.TFInitOptions, stringp(DCETFInitOptionsEnvVar), stringp(Empty))
		deployConfig.TFApplyOptions = *cfg.Coalesce(&deployConfig.TFApplyOptions, Config.Terraform.TFApplyOptions, stringp(DCETFApplyOptionsEnvVar), stringp(Empty))

		ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)
		if err := Service.Deploy(ctx, &deployOverrides); err != nil {
			// log.Fatalln(err)
			return err
		}
		return nil
	},
	// If the command was successful, read the config file back in
	// and if it's different than the running config, then save it.
	// We don't want to do this unless the command was succesful,
	// though, because of cases like bad tf opts we don't want to
	// create an usuable state.
	PostRunE: func(cmd *cobra.Command, args []string) error {
		deployConfig.TFInitOptions = *cfg.Coalesce(&deployConfig.TFInitOptions, Config.Terraform.TFInitOptions, stringp(DCETFInitOptionsEnvVar), stringp(Empty))
		deployConfig.TFApplyOptions = *cfg.Coalesce(&deployConfig.TFApplyOptions, Config.Terraform.TFApplyOptions, stringp(DCETFApplyOptionsEnvVar), stringp(Empty))

		ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)
		if err := Service.PostDeploy(ctx); err != nil {
			return err
		}
		return nil
	},
}

func stringp(s string) *string {
	return &s
}
