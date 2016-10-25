package controllers

import (
	. "github.com/jianqiu/vm-pool-server"
	"github.com/jianqiu/vm-pool-server/db"
	"github.com/jianqiu/vm-pool-server/models"
	"code.cloudfoundry.org/lager"
)

type VirtualGuestController struct {
	db                   db.VirtualGuestDB
	serviceClient        ServiceClient
}

func NewVirtualGuestController(
db db.VirtualGuestDB,
serviceClient ServiceClient,
) *VirtualGuestController {
	return &VirtualGuestController{
		db:                   db,
		serviceClient:        serviceClient,
	}
}

func (h *VirtualGuestController) VirtualGuests(logger lager.Logger, publicVlan, privateVlan, cpu, memory_mb int32) ([]*models.VirtualGuest, error) {
	logger = logger.Session("vms")

	filter := models.VirtualGuestFilter{
		PublicVlan: publicVlan,
		PrivateVlan: privateVlan,
	}

	return h.db.VirtualGuests(logger, filter)
}

func (h *VirtualGuestController) VirtualGuestByCid(logger lager.Logger, cid string) (*models.VirtualGuest, error) {
	logger = logger.Session("vm-by-guid")

	return h.db.VirtualGuestByCID(logger, cid)
}

func (h *VirtualGuestController) RemoveVirtualGuestFromPool(logger lager.Logger, cid int32) error {
	logger = logger.Session("delete-vm")

	return h.db.DeleteVirtualGuestFromPool(logger, cid)
}
