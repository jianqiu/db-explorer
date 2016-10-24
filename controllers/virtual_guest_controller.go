package controllers

import (
	"code.cloudfoundry.org/lager"
	"github.com/jianqiu/vm-pool-server"
	"github.com/jianqiu/vm-pool-server/db"
	"github.com/jianqiu/vm-pool-server/models"
)

type VirtualGuestController struct {
	db            db.VirtualGuestDB
	serviceClient ServiceClient
}

func NewVirtualGuestController(
	db db.VirtualGuestDB,
	serviceClient ServiceClient,
) *VirtualGuestController {
	return &VirtualGuestController{
		db:            db,
		serviceClient: serviceClient,
	}
}

func (h *VirtualGuestController) VirtualGuests(logger lager.Logger, filter models.VirtualGuestFilter) ([]*models.VirtualGuest, error) {
	logger = logger.Session("virtualguests")

	return h.db.VirtualGuests(logger, filter)
}

func (h *VirtualGuestController) RequestVirtualGuestByCID(logger lager.Logger, cid string) (*models.VirtualGuest, error) {
	logger = logger.Session("virutal-guest-by-guid")

	return h.db.VirtualGuestByCID(logger, cid)
}

func (h *VirtualGuestController) ReturnVirtualGuestIntoPool(logger lager.Logger, cid int32) error {
	logger = logger.Session("delete-virtual-guest")

	return h.db.DeleteVirtualGuestFromPool(logger, cid)
}
