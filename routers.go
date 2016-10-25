package vps

import "github.com/tedsuo/rata"

const (
// Ping
	PingRoute = "Ping"

// VMs
	VMsRoute = "Vms_r2"
	VMByCidRoute = "VMByGuid_r2"
	StartTaskRoute     = "StartTask"
	CancelTaskRoute    = "CancelTask"
	FailTaskRoute      = "FailTask"
	CompleteTaskRoute  = "CompleteTask"
	ResolvingTaskRoute = "ResolvingTask"
	DeleteTaskRoute    = "DeleteTask"
)

var Routes = rata.Routes{
	// Vms
	{Path: "/v1/vms/list.r2", Method: "POST", Name: VMsRoute},
	{Path: "/v1/vms/get_by_vm_cid.r2", Method: "POST", Name: VMByCidRoute},

	// VM Lifecycle
	{Path: "/v1/tasks/start", Method: "POST", Name: StartTaskRoute},
	{Path: "/v1/tasks/cancel", Method: "POST", Name: CancelTaskRoute},
	{Path: "/v1/tasks/fail", Method: "POST", Name: FailTaskRoute},
	{Path: "/v1/tasks/complete", Method: "POST", Name: CompleteTaskRoute},
	{Path: "/v1/tasks/resolving", Method: "POST", Name: ResolvingTaskRoute},
	{Path: "/v1/tasks/delete", Method: "POST", Name: DeleteTaskRoute},
}
