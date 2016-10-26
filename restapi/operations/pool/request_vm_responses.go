package pool

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/jianqiu/vm-pool-server/models"
)

/*RequestVMOK Requesting a VM from the pool succeeded

swagger:response requestVmOK
*/
type RequestVMOK struct {

	// In: body
	Payload *models.VM `json:"body,omitempty"`
}

// NewRequestVMOK creates RequestVMOK with default headers values
func NewRequestVMOK() *RequestVMOK {
	return &RequestVMOK{}
}

// WithPayload adds the payload to the request Vm o k response
func (o *RequestVMOK) WithPayload(payload *models.VM) *RequestVMOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the request Vm o k response
func (o *RequestVMOK) SetPayload(payload *models.VM) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RequestVMOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*RequestVMDefault unexpected error

swagger:response requestVmDefault
*/
type RequestVMDefault struct {
	_statusCode int

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// NewRequestVMDefault creates RequestVMDefault with default headers values
func NewRequestVMDefault(code int) *RequestVMDefault {
	if code <= 0 {
		code = 500
	}

	return &RequestVMDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the request VM default response
func (o *RequestVMDefault) WithStatusCode(code int) *RequestVMDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the request VM default response
func (o *RequestVMDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the request VM default response
func (o *RequestVMDefault) WithPayload(payload *models.Error) *RequestVMDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the request VM default response
func (o *RequestVMDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RequestVMDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}