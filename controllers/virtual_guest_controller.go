package controllers

import (
	. "github.com/vm-pool-server"
	"github.com/vm-pool-server/db"
	"github.com/vm-pool-server/models"
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

func (h *VirtualGuestController) VirtualGuests(logger lager.Logger, filter models.VirtualGuestFilter) ([]*models.VirtualGuest, error) {
	logger = logger.Session("virtualguests")

	return h.db.VirtualGuests(logger, filter)
}

func (h *VirtualGuestController) VirtualGuestByCid(logger lager.Logger, cid string) (*models.VirtualGuest, error) {
	logger = logger.Session("virutal-guest-by-guid")

	return h.db.VirtualGuestByCID(logger, cid)
}

func (h *VirtualGuestController) RemoveVirtualGuestFromPool(logger lager.Logger, cid int32) error {
	logger = logger.Session("delete-virtual-guest")

	return h.db.DeleteVirtualGuestFromPool(logger, cid)
}
