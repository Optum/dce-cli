package constants

type ContextKey string

const (
	DeployConfig     ContextKey = ContextKey("deployConfig")
	DeployLogFileKey ContextKey = ContextKey("deployLogFile")
)
