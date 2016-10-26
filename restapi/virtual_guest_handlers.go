package restapi

import (
	"code.cloudfoundry.org/cflager"
	"github.com/jianqiu/vm-pool-server/controllers"

	"github.com/jianqiu/vm-pool-server/models"
	"github.com/jianqiu/vm-pool-server/restapi/operations/pool"

	"github.com/go-openapi/runtime/middleware"
	"github.com/jianqiu/vm-pool-server/config"
)

//go:generate counterfeiter -o fake_controllers/fake_virtual_guest_controller.go . VirtualGuestController

func RequestVMHandleFunc (params pool.RequestVMParams) middleware.Responder {

	vmController := controllers.NewVirtualGuestController(config.ActiveDB)

	requestVM := pool.NewRequestVMOK()
	vm := models.VM{}

	logger, _ := cflager.New("req_handler")

	vms, _ := vmController.VirtualGuests(logger,params.Body.PublicVlan, params.Body.PrivateVlan, params.Body.CPU, params.Body.Memory)

	if len (vms) >0 {
		vm.VMID = vms[0].Cid
		vm.CPU = vms[0].Cpu
		vm.Memory = vms[0].MemoryMb
		vm.PublicVlan = vms[0].PublicVlan
		vm.PrivateVlan = vms[0].PrivateVlan
	}

	requestVM.SetPayload(&vm)
	return requestVM
}
