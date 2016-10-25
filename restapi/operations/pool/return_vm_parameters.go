package pool

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/jianqiu/vm-pool-server/models"
)

// NewReturnVMParams creates a new ReturnVMParams object
// with the default values initialized.
func NewReturnVMParams() ReturnVMParams {
	var ()
	return ReturnVMParams{}
}

// ReturnVMParams contains all the bound params for the return VM operation
// typically these are obtained from a http.Request
//
// swagger:parameters returnVM
type ReturnVMParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*VM ID
	  Required: true
	  In: body
	*/
	Body models.VMID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *ReturnVMParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.VMID
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("body", "body"))
			} else {
				res = append(res, errors.NewParseError("body", "body", "", err))
			}

		} else {

			if len(res) == 0 {
				o.Body = body
			}
		}

	} else {
		res = append(res, errors.Required("body", "body"))
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
