package observation

type ObservationContainer struct {
	Logger
}

func New(levelLogger LevelLogger) *ObservationContainer {

	logger := &LogObservation{
		LevelLogger: levelLogger,
	}

	return &ObservationContainer{
		Logger: logger,
	}
}

type Logger interface {
	LevelLogger
	// Endf logs and exits successfully (e.g. exit 0)
	Endf(format string, args ...interface{})
	// End logs and exits successfully (e.g. exit 0)
	End(args ...interface{})
	// Endln logs and exits successfully (e.g. exit 0)
	Endln(args ...interface{})
}

//LevelLogger contains common functions for printing logs at different levels
type LevelLogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
}
