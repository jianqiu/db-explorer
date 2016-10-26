package handlers

import (
	"net/http"

	"github.com/jianqiu/vm-pool-server/models"
	"code.cloudfoundry.org/lager"
	"github.com/tedsuo/rata"
	"strconv"
)

//go:generate counterfeiter -o fake_controllers/fake_virtual_guest_controller.go . VirtualGuestController

type VirtualGuestController interface {
	AllVirtualGuests(logger lager.Logger) ([]*models.VM, error)
	VirtualGuests(logger lager.Logger, publicVlan, privateVlan, cpu, memory_mb int32) ([]*models.VM, error)
	CreateVM(logger lager.Logger, vmkDefinition *models.VM) error
	DeleteVM(logger lager.Logger, cid int32) error
	UpdateVM(logger lager.Logger, cid int32, updateData *models.StateValue) error
	VirtualGuestByCid(logger lager.Logger, cid int32) (*models.VM, error)
}

type VirtualGuestHandler struct {
	controller VirtualGuestController
	exitChan   chan<- struct{}
}

type MessageValidator interface {
	Validate() error
	Unmarshal(data []byte) error
}

func NewVirtualGuestHandler(
controller VirtualGuestController,
exitChan chan<- struct{},
) *VirtualGuestHandler {
	return &VirtualGuestHandler{
		controller: controller,
		exitChan:   exitChan,
	}
}

func (h *VirtualGuestHandler) AllVirtualGuests(logger lager.Logger, w http.ResponseWriter, req *http.Request) {
	var err error
	logger = logger.Session("virtualguests")

	response := &models.VMsResponse{}

	defer func() { exitIfUnrecoverable(logger, h.exitChan, response.Error) }()
	defer func() { writeResponse(w, response) }()

	if err != nil {
		logger.Error("failed-parsing-request", err)
		response.Error = models.ConvertError(err)
		return
	}

	response.Vms, err = h.controller.AllVirtualGuests(logger)
	response.Error = models.ConvertError(err)
}

func (h *VirtualGuestHandler) VirtualGuest(logger lager.Logger, w http.ResponseWriter, req *http.Request) {
	var err error
	logger = logger.Session("virtualguests")

	response := &models.VMResponse{}

	vmId := rata.Param(req, "cid")
	if len(vmId) == 0 {
		logger.Error("failed-parsing-request", err)
		response.Error = models.ConvertError(err)
		return
	}

	defer func() { exitIfUnrecoverable(logger, h.exitChan, response.Error) }()
	defer func() { writeResponse(w, response) }()

	cid, err:= strconv.ParseInt(vmId, 10, 32)
	if err != nil {
		logger.Error("failed-parsing-request", err)
		response.Error = models.ConvertError(err)
		return
	}

	response.Vm, err = h.controller.VirtualGuestByCid(logger,int32(cid))
	response.Error = models.ConvertError(err)
}

func (h *VirtualGuestHandler) VirtualGuests(logger lager.Logger, w http.ResponseWriter, req *http.Request) {
	var err error
	logger = logger.Session("virtualguests")

	request := &models.VMRequestFilter{}
	response := &models.VMsResponse{}

	defer func() { exitIfUnrecoverable(logger, h.exitChan, response.Error) }()
	defer func() { writeResponse(w, response) }()

	err = parseRequest(logger, req, request)
	if err != nil {
		logger.Error("failed-parsing-request", err)
		response.Error = models.ConvertError(err)
		return
	}

	response.Vms, err = h.controller.VirtualGuests(logger, request.PublicVlan, request.PrivateVlan, request.Cpu, request.MemoryMb)
	response.Error = models.ConvertError(err)
}

func (h *VirtualGuestHandler) CreateVM(logger lager.Logger, w http.ResponseWriter, req *http.Request) {
	var err error
	logger = logger.Session("create_vm")

	request := &models.VM{}
	response := &models.VMLifecycleResponse{}

	defer func() { exitIfUnrecoverable(logger, h.exitChan, response.Error) }()
	defer func() { writeResponse(w, response) }()

	err = parseRequest(logger, req, request)
	if err != nil {
		logger.Error("failed-parsing-request", err)
		response.Error = models.ConvertError(err)
		return
	}

	err = h.controller.CreateVM(logger, request)
	response.Error = models.ConvertError(err)
}

func (h *VirtualGuestHandler) DeleteVM(logger lager.Logger, w http.ResponseWriter, req *http.Request) {
	var err error
	logger = logger.Session("delete-vm")

	response := &models.VMLifecycleResponse{}

	vmId := rata.Param(req, "cid")
	if len(vmId) == 0 {
		logger.Error("failed-parsing-request", err)
		response.Error = models.ConvertError(err)
		return
	}

	defer func() { exitIfUnrecoverable(logger, h.exitChan, response.Error) }()
	defer func() { writeResponse(w, response) }()

	cid, err:= strconv.ParseInt(vmId, 10, 32)
	if err != nil {
		logger.Error("failed-parsing-request", err)
		response.Error = models.ConvertError(err)
		return
	}

	err = h.controller.DeleteVM(logger, int32(cid))
	response.Error = models.ConvertError(err)
}

func (h *VirtualGuestHandler) UpdateVM(logger lager.Logger, w http.ResponseWriter, req *http.Request) {
	var err error
	logger = logger.Session("update-vm")

	response := &models.VMLifecycleResponse{}

	vmId := rata.Param(req, "cid")
	if len(vmId) == 0 {
		logger.Error("failed-parsing-request", err)
		response.Error = models.ConvertError(err)
		return
	}

	updateData := &models.StateValue{}

	defer func() { exitIfUnrecoverable(logger, h.exitChan, response.Error) }()
	defer func() { writeResponse(w, response) }()

	err = parseRequest(logger, req, updateData)
	if err != nil {
		logger.Error("failed-parsing-request", err)
		response.Error = models.ConvertError(err)
		return
	}

	cid, err:= strconv.ParseInt(vmId, 10, 32)
	if err != nil {
		logger.Error("failed-parsing-request", err)
		response.Error = models.ConvertError(err)
		return
	}

	err = h.controller.UpdateVM(logger, int32(cid), updateData)
	response.Error = models.ConvertError(err)
}