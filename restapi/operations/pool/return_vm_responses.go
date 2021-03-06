package pool

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/jianqiu/vm-pool-server/models"
)

/*ReturnVMOK Returning a VM into the pool succeeded

swagger:response returnVmOK
*/
type ReturnVMOK struct {
}

// NewReturnVMOK creates ReturnVMOK with default headers values
func NewReturnVMOK() *ReturnVMOK {
	return &ReturnVMOK{}
}

// WriteResponse to the client
func (o *ReturnVMOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
}

/*ReturnVMDefault unexpected error

swagger:response returnVmDefault
*/
type ReturnVMDefault struct {
	_statusCode int

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// NewReturnVMDefault creates ReturnVMDefault with default headers values
func NewReturnVMDefault(code int) *ReturnVMDefault {
	if code <= 0 {
		code = 500
	}

	return &ReturnVMDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the return VM default response
func (o *ReturnVMDefault) WithStatusCode(code int) *ReturnVMDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the return VM default response
func (o *ReturnVMDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the return VM default response
func (o *ReturnVMDefault) WithPayload(payload *models.Error) *ReturnVMDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the return VM default response
func (o *ReturnVMDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ReturnVMDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
