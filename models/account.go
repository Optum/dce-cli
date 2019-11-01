// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Account Account Details
// swagger:model account
type Account struct {

	// Status of the Account.
	// "Ready": The account is clean and ready for lease
	// "NotReady": The account is in "dirty" state, and needs to be reset before it may be leased.
	// "Leased": The account is leased to a principal
	//
	// Enum: [Ready NotReady Leased]
	AccountStatus string `json:"accountStatus,omitempty"`

	// ARN for an IAM role within this AWS account. The Redbox master account will assume this IAM role to execute operations within this AWS account. This IAM role is configured by the client, and must be configured with [a Trust Relationship with the Redbox master account.](/https://docs.aws.amazon.com/IAM/latest/UserGuide/tutorial_cross-account-with-roles.html)
	AdminRoleArn string `json:"adminRoleArn,omitempty"`

	// Epoch timestamp, when account record was created
	CreatedOn int64 `json:"createdOn,omitempty"`

	// AWS Account ID
	ID string `json:"id,omitempty"`

	// Epoch timestamp, when account record was last modified
	LastModifiedOn int64 `json:"lastModifiedOn,omitempty"`

	// Any organization specific data pertaining to the account that needs to be persisted
	Metadata interface{} `json:"metadata,omitempty"`

	// The S3 object ETag used to apply the Principal IAM Policy within this AWS account.  This policy is created by the Redbox master account, and is assumed by people with access to principalRoleArn.
	PrincipalPolicyHash string `json:"principalPolicyHash,omitempty"`

	// ARN for an IAM role within this AWS account. This role is created by the Redbox master account, and may be assumed by principals to login to their AWS Redbox account.
	PrincipalRoleArn string `json:"principalRoleArn,omitempty"`
}

// Validate validates this account
func (m *Account) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccountStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var accountTypeAccountStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["Ready","NotReady","Leased"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		accountTypeAccountStatusPropEnum = append(accountTypeAccountStatusPropEnum, v)
	}
}

const (

	// AccountAccountStatusReady captures enum value "Ready"
	AccountAccountStatusReady string = "Ready"

	// AccountAccountStatusNotReady captures enum value "NotReady"
	AccountAccountStatusNotReady string = "NotReady"

	// AccountAccountStatusLeased captures enum value "Leased"
	AccountAccountStatusLeased string = "Leased"
)

// prop value enum
func (m *Account) validateAccountStatusEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, accountTypeAccountStatusPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Account) validateAccountStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.AccountStatus) { // not required
		return nil
	}

	// value enum
	if err := m.validateAccountStatusEnum("accountStatus", "body", m.AccountStatus); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Account) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Account) UnmarshalBinary(b []byte) error {
	var res Account
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}