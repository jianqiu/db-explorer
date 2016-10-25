package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/jianqiu/vm-pool-server/models"
	"github.com/jianqiu/vm-pool-server/restapi/operations"
	"github.com/jianqiu/vm-pool-server/restapi/operations/pool"
	"github.com/jianqiu/vm-pool-server/restapi/operations/vms"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name vm-pool-server --spec ../docs/vm_pool_server_api.json

func configureFlags(api *operations.VMPoolServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.VMPoolServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.VmsGetVmsHandler = vms.GetVmsHandlerFunc(func(params vms.GetVmsParams) middleware.Responder {
		return middleware.NotImplemented("operation vms.GetVms has not yet been implemented")
	})
	api.PoolRequestVMHandler = pool.RequestVMHandlerFunc(func(params pool.RequestVMParams) middleware.Responder {

		requestVM := pool.NewRequestVMOK()
		vm := models.VM{
			CPU:         4,
			Deployment:  "test-deployment",
			Hostname:    "test-hostname",
			Memory:      32768,
			PrivateIP:   "10.0.0.99",
			PrivateVlan: 123456,
			PublicVlan:  123457,
			Status:      "in_req",
			VMID:        *params.VMID,
		}
		requestVM.SetPayload(&vm)
		return requestVM
	})
	api.PoolReturnVMHandler = pool.ReturnVMHandlerFunc(func(params pool.ReturnVMParams) middleware.Responder {
		return middleware.NotImplemented("operation pool.ReturnVM has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
