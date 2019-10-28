package service

import (
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
	"github.com/Optum/dce-cli/mocks"
	"github.com/sirupsen/logrus"
)

type TestLogObservation struct {
	observ.LevelLogger
	Ended bool
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
var spyLogger TestLogObservation
var service *ServiceContainer

func initMocks(config configs.Root) {
	mockPrompter = mocks.Prompter{}
	mockFileSystemer = mocks.FileSystemer{}
	mockWeber = mocks.Weber{}
	spyLogger = TestLogObservation{
		logrus.New(),
		false,
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
	}
	service = New(&config, &spyObservation, &mockUtil)
}
