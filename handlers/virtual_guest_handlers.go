package handlers

import (
	"net/http"

	"github.com/jianqiu/vm-pool-server/models"
	"code.cloudfoundry.org/lager"
)

//go:generate counterfeiter -o fake_controllers/fake_virtual_guest_controller.go . VirtualGuestController

type VirtualGuestController interface {
	VirtualGuests(logger lager.Logger, publicVlan, privateVlan, cpu, memory_mb int32) ([]*models.VM, error)
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