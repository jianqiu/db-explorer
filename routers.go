package vps

import "github.com/tedsuo/rata"

const (
// Ping
	PingRoute = "Ping"

// VMs
	VMsRoute = "Vms_r2"
)

var Routes = rata.Routes{
	// Vms
	{Path: "/v1/vms/list.r2", Method: "POST", Name: VMsRoute},
}
