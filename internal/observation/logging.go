package observation

import "os"

type LogObservation struct {
	LevelLogger
}

func (l *LogObservation) Endf(format string, args ...interface{}) {
	l.Infof(format, args)
	os.Exit(0)
}

func (l *LogObservation) End(args ...interface{}) {
	l.Info(args)
	os.Exit(0)
}

func (l *LogObservation) Endln(args ...interface{}) {
	l.Infoln(args)
	os.Exit(0)
}
