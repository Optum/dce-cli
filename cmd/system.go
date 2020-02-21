package cmd

import (
	"fmt"
	"github.com/Optum/dce-cli/pkg/service"
	"github.com/aws/aws-sdk-go/aws"
	"math/rand"
	"strings"
	"time"

	cfg "github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	"github.com/spf13/cobra"
)

var (
	DeployConfig *service.DeployConfig
)

const (
	DCETFInitOptionsEnvVar               string = "DCE_TF_INIT_OPTIONS"
	DCETFApplyOptionsEnvVar              string = "DCE_TF_APPLY_OPTIONS"
	DCEDeployLogEnvVar                   string = "DCE_DEPLOY_LOG_FILE"
	DCELocationEnvVar                    string = "DCE_LOCATION"
	DCEVersionEnvVar                     string = "DCE_VERSION"
	DCENamespaceEnvVar                   string = "DCE_NAMESPACE"
	DCEBudgetNotificationFromEmailEnvVar string = "DCE_BUDGET_NOTIFICATION_FROM_EMAIL"
	Empty                                string = ""
)

func init() {
	DeployConfig = &service.DeployConfig{}
	systemDeployCmd.Flags().StringVar(&DeployConfig.Location, "local", "", "[DEPRECATED use --location instead] Path to a local DCE repo to deploy.")
	systemDeployCmd.Flags().StringVarP(&DeployConfig.Location, "location", "l", "", "Path to the DCE repo. May be a local path (eg. /path/to/dce) or a github location (eg. github.com/Optum/dce)")
	systemDeployCmd.Flags().StringVarP(&DeployConfig.Version, "dce-version", "d", "", fmt.Sprintf("Version of DCE to deploy. Defaults to v%s", constants.DefaultDCEVersion))
	systemDeployCmd.Flags().BoolVarP(aws.Bool(true), "use-cached", "c", true, "[DEPRECATED] Uses locally-cached files, if available.")
	systemDeployCmd.Flags().BoolVarP(&DeployConfig.BatchMode, "batch-mode", "b", false, "Skip prompting for resource creation.")
	systemDeployCmd.Flags().StringVar(&DeployConfig.TFInitOptions, "tf-init-options", "", "Options to pass to the underlying \"tf init\" command.")
	systemDeployCmd.Flags().StringVar(&DeployConfig.TFApplyOptions, "tf-apply-options", "", "Options to pass to the underlying \"tf apply\" command.")
	systemDeployCmd.Flags().BoolVar(&DeployConfig.SaveTFOptions, "save-options", false, "If specified, saves the values provided by \"--tf-init-options\" and \"--tf-apply-options\" in the config file.")
	systemDeployCmd.Flags().StringVarP(&DeployConfig.Namespace, "namespace", "n", "", "Set a custom terraform namespace (Optional)")
	systemDeployCmd.Flags().StringVarP(&DeployConfig.AWSRegion, "region", "r", "", "The aws region that DCE will be deployed to (Default: us-east-1)")
	systemDeployCmd.Flags().StringArrayVarP(&DeployConfig.GlobalTags, "tag", "t", []string{}, "Tags to be placed on all DCE resources. E.g. \"dce system deploy --tag key1:value1 --tag key2:value2\"")
	systemDeployCmd.Flags().StringVar(&DeployConfig.BudgetNotificationFromEmail, "budget-notification-from-email", "", "Email address from which budget notifications will be sent")
	systemDeployCmd.Flags().StringArrayVar(&DeployConfig.BudgetNotificationBCCEmails, "budget-notification-bcc-emails", []string{}, "Email address from which budget notifications will be sent")
	systemDeployCmd.Flags().StringVar(&DeployConfig.BudgetNotificationTemplateHTML, "budget-notification-template-html", "", "HTML template for budget notification emails")
	systemDeployCmd.Flags().StringVar(&DeployConfig.BudgetNotificationTemplateText, "budget-notification-template-text", "", "Text template for budget notification emails")
	systemDeployCmd.Flags().StringVar(&DeployConfig.BudgetNotificationTemplateSubject, "budget-notification-template-subject", "", "Subjet for budget notification emails")
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
		// Coalesce the deployment configuration values,
		// loading in from CLI, YAML, env vars, and default values
		DeployConfig = coalesceDeployConfig(DeployConfig)

		if err := Service.Deploy(DeployConfig); err != nil {
			cmd.SilenceUsage = true
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
		return Service.PostDeploy(DeployConfig)
	},
}

func coalesceDeployConfig(deployConfig *service.DeployConfig) *service.DeployConfig {
	// Before running the command, resolve the configuration by using
	// the default load order.
	//
	// Config loader order:
	//		1. CLI args
	//		2. YAML config
	// 		3. Env vars
	// 		4. Default value
	deployConfig.Location = *cfg.Coalesce(
		&deployConfig.Location,
		Config.Deploy.Location,
		stringp(DCELocationEnvVar),
		stringp(constants.DefaultDCELocation),
	)
	deployConfig.Version = *cfg.Coalesce(
		&deployConfig.Version,
		Config.Deploy.Version,
		stringp(DCEVersionEnvVar),
		stringp(constants.DefaultDCEVersion),
	)
	// Normalize version ("v1.2.3" --> "1.2.3")
	// to make our CLI more forgiving
	if strings.HasPrefix(deployConfig.Version, "v") {
		deployConfig.Version = deployConfig.Version[1:]
	}

	deployConfig.TFInitOptions = *cfg.Coalesce(
		&deployConfig.TFInitOptions,
		Config.Terraform.TFInitOptions,
		stringp(DCETFInitOptionsEnvVar),
		stringp(Empty),
	)
	deployConfig.TFApplyOptions = *cfg.Coalesce(
		&deployConfig.TFApplyOptions,
		Config.Terraform.TFApplyOptions,
		stringp(DCETFApplyOptionsEnvVar),
		stringp(Empty),
	)
	deployConfig.DeployLogFile = *cfg.Coalesce(
		&deployConfig.DeployLogFile,
		Config.Deploy.LogFile,
		stringp(DCEDeployLogEnvVar),
		stringp(Util.GetLogFile()),
	)
	deployConfig.AWSRegion = *cfg.Coalesce(
		&deployConfig.AWSRegion,
		Config.Deploy.AWSRegion,
		stringp("AWS_REGION"),
		stringp("us-east-1"),
	)
	deployConfig.Namespace = *cfg.Coalesce(
		&deployConfig.Namespace,
		Config.Deploy.Namespace,
		stringp(DCENamespaceEnvVar),
		stringp("dce-"+getRandString(8)),
	)
	deployConfig.BudgetNotificationFromEmail = *cfg.Coalesce(
		&deployConfig.BudgetNotificationFromEmail,
		Config.Deploy.BudgetNotificationFromEmail,
		stringp(DCEBudgetNotificationFromEmailEnvVar),
		stringp("no-reply@example.com"),
	)

	return deployConfig
}

func stringp(s string) *string {
	return &s
}

func getRandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
