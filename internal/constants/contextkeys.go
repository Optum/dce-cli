package constants

type ContextKey string

const (
	DeployConfig  ContextKey = ContextKey("deployConfig")
	DeployLogFile ContextKey = ContextKey("deployLogFile")
)
