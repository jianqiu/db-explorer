package vps

import "github.com/tedsuo/rata"

const (
// Ping
	PingRoute = "Ping"

// VMs
	AllVMsRoute ="all_vms"
	VMsRoute = "list_vms"
	VMRoute  = "get_vm"
	VMUpdateRoute ="update_vm"
	VMCreateRoute ="create_vm"
	VMDeleteRoute ="delete_vm"
)

var Routes = rata.Routes{
	// Vms
	{Path: "/v1/vms", Method: rata.GET, Name: AllVMsRoute},
	{Path: "/v1/vms/list", Method: rata.POST, Name: VMsRoute},
	{Path: "/v1/vms/:cid", Method: rata.GET, Name: VMRoute},
	{Path: "/v1/vms/:cid", Method: rata.PUT, Name: VMUpdateRoute},
	{Path: "/v1/vms", Method: rata.POST, Name: VMCreateRoute},
	{Path: "/v1/vms/:cid", Method: rata.DELETE, Name: VMDeleteRoute},
}