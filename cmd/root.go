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
	"fmt"
	"os"
	"path/filepath"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
	svc "github.com/Optum/dce-cli/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var config *configs.Root = &configs.Root{}
var service *svc.ServiceContainer
var util *utl.UtilContainer
var observation *observ.ObservationContainer

func init() {
	cobra.OnInitialize(initObservation, initConfig, initUtil, initService)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dce.yaml)")
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dce",
	Short: "Disposable Cloud Environment (DCE)",
	Long: `Disposable Cloud Environment (DCE) 

  The DCE cli allows:

  - Admins to provision DCE to a master account and administer said account
  - Users to lease accounts and execute commands against them`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initialize anything related to logging, metrics, or tracing
func initObservation() {
	logrusInstance := logrus.New()
	observation = observ.New(logrusInstance)
}

// Utils we need before they are normally instantiated
var log observ.Logger
var fsUtil utl.FileSystemer

// initialize config from file or tell user to run 'dce init' if none exists
func initConfig() {
	tempUtil := utl.New(config, observation)
	fsUtil = tempUtil.FileSystemer
	log = observation.Logger

	var configFileUsed string
	if cfgFile != "" {
		configFileUsed = cfgFile
		viper.SetConfigFile(configFileUsed)
	} else {
		home := fsUtil.GetHomeDir()
		configFileUsed = filepath.Join(home, constants.DefaultConfigFileName)
		viper.SetConfigFile(configFileUsed)
	}

	if !fsUtil.IsExistingFile(configFileUsed) {
		if len(os.Args) < 2 || os.Args[1] != initCmd.Name() {
			log.Endln("Config file not found. Please type 'dce init' to generate one.")
		}
	} else {
		if err := viper.ReadInConfig(); err == nil {
			viper.BindEnv("api.credentials.awsaccesskeyid", "AWS_ACCESS_KEY_ID")
			viper.BindEnv("api.credentials.awssecretaccesskey", "AWS_SECRET_ACCESS_KEY")
			viper.BindEnv("api.credentials.awssessiontoken", "AWS_SESSION_TOKEN")
			viper.BindEnv("githubtoken", "GITHUB_TOKEN")
			viper.Unmarshal(config)
		}
	}
}

// initialize utilities and interfaces to external things
func initUtil() {
	util = utl.New(config, observation)
}

// initialize business logic
func initService() {
	service = svc.New(config, observation, util)
}
