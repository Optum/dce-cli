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
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"

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

func init() {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Global Flags
	// ---------------
	// --config flag, to specify path to dce.yml config
	// default to ~/.dce.yml
	RootCmd.PersistentFlags().StringVar(
		&cfgFile, "config",
		filepath.Join(homeDir, ".dce", constants.DefaultConfigFileName),
		"config file",
	)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
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

	// Check if the user has valid creds,
	// otherwise require authentication
	creds := Util.AWSSession.Config.Credentials
	_, _ = creds.Get()
	hasValidCreds := !creds.IsExpired()
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
	var serialized []byte
	serialized = []byte(fmt.Sprintf("%s\n", entry.Message))
	return serialized, nil
}

func onInit(cmd *cobra.Command, args []string) error {
	// Configure observation / logging
	initObservation()
	// Expose global `log` object for ease of use
	log = Observation.Logger
	Log = log

	fsUtil := &utl.FileSystemUtil{Config: Config, ConfigFile: cfgFile}

	// Initialize config
	// If config file does not exist,
	// run the `dce init` command
	if !fsUtil.IsExistingFile(cfgFile) {
		if cmd.Name() != initCmd.Name() {
			return errors.New("Config file not found. Please type 'dce init' to generate one.")
		}
	} else {
		// Load config from dce.yaml
		err := fsUtil.ReadInConfig()
		if err != nil {
			return fmt.Errorf("Failed to parse dce.yml: %s", err)
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
