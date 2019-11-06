package unit

import (
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
	"github.com/Optum/dce-cli/mocks"
	svc "github.com/Optum/dce-cli/pkg/service"
	"github.com/sirupsen/logrus"
)

type TestLogObservation struct {
	observ.LevelLogger
	Ended bool
	Msg   string
}

func (l *TestLogObservation) Infoln(args ...interface{}) {
	l.Msg = args[0].(string)
	l.Ended = true
}

func (l *TestLogObservation) Endf(format string, args ...interface{}) {
	l.Ended = true
}

func (l *TestLogObservation) End(args ...interface{}) {
	l.Ended = true
}

func (l *TestLogObservation) Endln(args ...interface{}) {
	l.Ended = true
}

var mockPrompter mocks.Prompter
var mockFileSystemer mocks.FileSystemer
var mockWeber mocks.Weber
var mockGithuber mocks.Githuber
var mockAwser mocks.AWSer
var mockTerraformer mocks.Terraformer
var mockAPIer mocks.APIer
var spyLogger TestLogObservation
var service *svc.ServiceContainer

func initMocks(config configs.Root) {
	mockPrompter = mocks.Prompter{}
	mockFileSystemer = mocks.FileSystemer{}
	mockWeber = mocks.Weber{}
	mockGithuber = mocks.Githuber{}
	mockAwser = mocks.AWSer{}
	mockTerraformer = mocks.Terraformer{}
	mockAPIer = mocks.APIer{}
	spyLogger = TestLogObservation{
		logrus.New(),
		false,
		"",
	}
	spyObservation := observ.ObservationContainer{
		Logger: &spyLogger,
	}
	mockUtil := utl.UtilContainer{
		Config:       &config,
		Prompter:     &mockPrompter,
		FileSystemer: &mockFileSystemer,
		Weber:        &mockWeber,
		Observation:  &spyObservation,
		Githuber:     &mockGithuber,
		AWSer:        &mockAwser,
		Terraformer:  &mockTerraformer,
		APIer:        &mockAPIer,
	}
	service = svc.New(&config, &spyObservation, &mockUtil)
}
