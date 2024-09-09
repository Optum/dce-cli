//go:build !windows
// +build !windows

package constants

//nolint:gosec
const CredentialsExport string = `export AWS_ACCESS_KEY_ID=%s
export AWS_SECRET_ACCESS_KEY=%s
export AWS_SESSION_TOKEN=%s`
