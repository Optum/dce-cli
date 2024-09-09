// Created by interfacer; DO NOT EDIT

package util

import (
	"github.com/Optum/dce-cli/client/operations"
	"github.com/go-openapi/runtime"
)

// APIer is an interface generated for "github.com/Optum/dce-cli/client/operations.Client".
type APIer interface {
	DeleteAccountsID(*operations.DeleteAccountsIDParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.DeleteAccountsIDNoContent, error)
	DeleteLeases(*operations.DeleteLeasesParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.DeleteLeasesOK, error)
	DeleteLeasesID(*operations.DeleteLeasesIDParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.DeleteLeasesIDOK, error)
	GetAccounts(*operations.GetAccountsParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.GetAccountsOK, error)
	GetAccountsID(*operations.GetAccountsIDParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.GetAccountsIDOK, error)
	GetAuth(*operations.GetAuthParams, ...operations.ClientOption) (*operations.GetAuthOK, error)
	GetAuthFile(*operations.GetAuthFileParams, ...operations.ClientOption) (*operations.GetAuthFileOK, error)
	GetLeases(*operations.GetLeasesParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.GetLeasesOK, error)
	GetLeasesID(*operations.GetLeasesIDParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.GetLeasesIDOK, error)
	GetUsage(*operations.GetUsageParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.GetUsageOK, error)
	PostAccounts(*operations.PostAccountsParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.PostAccountsCreated, error)
	PostLeases(*operations.PostLeasesParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.PostLeasesCreated, error)
	PostLeasesIDAuth(*operations.PostLeasesIDAuthParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.PostLeasesIDAuthCreated, error)
	PutAccountsID(*operations.PutAccountsIDParams, runtime.ClientAuthInfoWriter, ...operations.ClientOption) (*operations.PutAccountsIDOK, error)
	SetTransport(runtime.ClientTransport)
}
