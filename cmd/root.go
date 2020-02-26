/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"os"
	"path/filepath"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
	svc "github.com/Optum/dce-cli/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string
var Config = &configs.Root{}
var Service *svc.ServiceContainer
var Util *utl.UtilContainer
var Observation *observ.ObservationContainer

// Expose logger as global for ease of use
var log observ.Logger
var Log observ.Logger

// Allow injecting a PostInit phase.
// Useful for mocking global services in tests
var PostInit func(cmd *cobra.Command, args []string) error

func init() {
	// Global Flags
	// ---------------
	// --config flag, to specify path to dce.yml config
	// default to ~/.dce/config.yaml
	RootCmd.PersistentFlags().StringVar(
		&cfgFile, "config",
		"",
		"config file (default is \"$HOME/.dce/config.yaml\")",
	)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command {
	Use:   "dce",
	Short: "Disposable Cloud Environment (DCE)",
	Long: `Disposable Cloud Environment (DCE) 

  The DCE cli allows:

  - Admins to provision DCE to a master account and administer said account
  - Users to lease accounts and execute commands against them`,
	PersistentPreRunE: preRun,
}

func preRun(cmd *cobra.Command, args []string) error {
	err := onInit(cmd, args)
	if err != nil {
		return err
	}

	if PostInit != nil {
		err = PostInit(cmd, args)
		if err != nil {
			return err
		}
	}

	// Check if the requested command is for a version check
	// If it is, return here, as no creds are needed
	if cmd.Name() == versionCmd.Name() {
		return nil
	}

	// Check if the user has valid creds,
	// otherwise require authentication
	hasValidCreds := areCredsValid(Util.AWSSession.Config.Credentials)
	isAuthCommand := cmd.Name() == authCmd.Name()
	isInitCommand := cmd.Name() == initCmd.Name()
	if !hasValidCreds && !isAuthCommand && !isInitCommand {
		log.Print("No valid DCE credentials found")
		err := Service.Authenticate()
		if err != nil {
			return err
		}
	}

	return nil
}

func areCredsValid(creds *credentials.Credentials) bool {
	// Check if the user has valid creds,
	// otherwise require authentication
	_, err := creds.Get()
	if err != nil {
		return false
	}
	areExpired := creds.IsExpired()
	if areExpired {
		return false
	}

	// There's a bug in the AWS SDK, that show expired creds
	// as non-expired: https://github.com/aws/aws-sdk-go/issues/3163
	// To verify valid creds, make a `sts.GetCallerIdentity` call
	sess, err := session.NewSession(&aws.Config{Credentials: creds})
	if err != nil {
		return false
	}
	stsSvc := sts.New(sess)
	_, err = stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	return err == nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	// Print an extra newline when we're done,
	// so users terminal prompt shows up on a new line
	fmt.Println("")
}

type FmtOutputFormatter struct {
}

func (f *FmtOutputFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	serialized := []byte(fmt.Sprintf("%s\n", entry.Message))
	return serialized, nil
}

func onInit(cmd *cobra.Command, args []string) error {
	// Configure observation / logging
	initObservation()
	// Expose global `log` object for ease of use
	log = Observation.Logger
	Log = log

	if len(cfgFile) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		cfgFile = filepath.Join(homeDir, ".dce", constants.DefaultConfigFileName)
	}

	fsUtil := &utl.FileSystemUtil{Config: Config, ConfigFile: cfgFile}

	// Initialize config
	// If config file does not exist,
	// run the `dce init` command
	if !fsUtil.IsExistingFile(cfgFile) {
		if cmd.Name() == versionCmd.Name() {
			return nil
		} else if cmd.Name() != initCmd.Name() {
			return errors.New("Config file not found. Please type 'dce init' to generate one.")
		}
	} else {
		// Load config from the configuration file
		err := fsUtil.ReadInConfig()
		if err != nil {
			return fmt.Errorf("Failed to parse configuration file: %s", err)
		}
	}

	// initialize utilities and interfaces to external things
	Util = utl.New(Config, cfgFile, Observation)

	// initialize business logic services
	Service = svc.New(Config, Observation, Util)

	return nil
}

// initialize anything related to logging, metrics, or tracing
func initObservation() {
	logrusInstance := logrus.New()
	logrusInstance.SetOutput(os.Stderr)

	//TODO: Make configurable
	var logLevel logrus.Level
	switch os.Getenv("DCE_LOG_LEVEL") {
	case "TRACE":
		logLevel = logrus.TraceLevel
	case "DEBUG":
		logLevel = logrus.DebugLevel
	case "INFO":
		logLevel = logrus.InfoLevel
	case "WARN":
		logLevel = logrus.WarnLevel
	case "ERROR":
		logLevel = logrus.ErrorLevel
	case "FATAL":
		logLevel = logrus.FatalLevel
	case "PANIC":
		logLevel = logrus.PanicLevel
	default:
		logLevel = logrus.InfoLevel
	}

	logrusInstance.SetLevel(logLevel)
	if logLevel == logrus.InfoLevel {
		logrusInstance.SetFormatter(&FmtOutputFormatter{})
	} else {
		logrusInstance.SetFormatter(&logrus.TextFormatter{})
	}

	Observation = observ.New(logrusInstance)
}
