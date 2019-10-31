// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewPostLeasesIDAuthParams creates a new PostLeasesIDAuthParams object
// with the default values initialized.
func NewPostLeasesIDAuthParams() *PostLeasesIDAuthParams {
	var ()
	return &PostLeasesIDAuthParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostLeasesIDAuthParamsWithTimeout creates a new PostLeasesIDAuthParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostLeasesIDAuthParamsWithTimeout(timeout time.Duration) *PostLeasesIDAuthParams {
	var ()
	return &PostLeasesIDAuthParams{

		timeout: timeout,
	}
}

// NewPostLeasesIDAuthParamsWithContext creates a new PostLeasesIDAuthParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostLeasesIDAuthParamsWithContext(ctx context.Context) *PostLeasesIDAuthParams {
	var ()
	return &PostLeasesIDAuthParams{

		Context: ctx,
	}
}

// NewPostLeasesIDAuthParamsWithHTTPClient creates a new PostLeasesIDAuthParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostLeasesIDAuthParamsWithHTTPClient(client *http.Client) *PostLeasesIDAuthParams {
	var ()
	return &PostLeasesIDAuthParams{
		HTTPClient: client,
	}
}

/*PostLeasesIDAuthParams contains all the parameters to send to the API endpoint
for the post leases ID auth operation typically these are written to a http.Request
*/
type PostLeasesIDAuthParams struct {

	/*ID
	  Id for lease

	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post leases ID auth params
func (o *PostLeasesIDAuthParams) WithTimeout(timeout time.Duration) *PostLeasesIDAuthParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post leases ID auth params
func (o *PostLeasesIDAuthParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post leases ID auth params
func (o *PostLeasesIDAuthParams) WithContext(ctx context.Context) *PostLeasesIDAuthParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post leases ID auth params
func (o *PostLeasesIDAuthParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post leases ID auth params
func (o *PostLeasesIDAuthParams) WithHTTPClient(client *http.Client) *PostLeasesIDAuthParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post leases ID auth params
func (o *PostLeasesIDAuthParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the post leases ID auth params
func (o *PostLeasesIDAuthParams) WithID(id string) *PostLeasesIDAuthParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the post leases ID auth params
func (o *PostLeasesIDAuthParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *PostLeasesIDAuthParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
