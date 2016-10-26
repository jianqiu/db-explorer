package handlers

import (
	"net/http"

	"github.com/jianqiu/vm-pool-server/models"
	"code.cloudfoundry.org/lager"
	"github.com/gogo/protobuf/proto"
)

//go:generate counterfeiter -o fake_controllers/fake_virtual_guest_controller.go . VirtualGuestController

type VirtualGuestController interface {
	VirtualGuests(logger lager.Logger, publicVlan, privateVlan, cpu, memory_mb int32) ([]*models.VirtualGuest, error)
	VirtualGuestByCid(logger lager.Logger, cid int32) (*models.VirtualGuest, error)
}

type VirtualGuestHandler struct {
	controller VirtualGuestController
	exitChan   chan<- struct{}
}

type MessageValidator interface {
	proto.Message
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

	request := &models.VMsRequest{}
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

func (h *VirtualGuestHandler) VirtualGuestByCid(logger lager.Logger, w http.ResponseWriter, req *http.Request) {
	var err error
	logger = logger.Session("virtual-guest-by-cid")

	request := &models.VMByCidRequest{}
	response := &models.VMResponse{}

	defer func() { exitIfUnrecoverable(logger, h.exitChan, response.Error) }()
	defer func() { writeResponse(w, response) }()

	err = parseRequest(logger, req, request)
	if err != nil {
		logger.Error("failed-parsing-request", err)
		response.Error = models.ConvertError(err)
		return
	}

	response.Vm, err = h.controller.VirtualGuestByCid(logger, request.Cid)
	response.Error = models.ConvertError(err)
}
