package controllers

import (
	"code.cloudfoundry.org/lager"
	"github.com/jianqiu/vm-pool-server/db"
	"github.com/jianqiu/vm-pool-server/models"
)

type VirtualGuestController struct {
	db            db.VirtualGuestDB
}

func NewVirtualGuestController(
	db db.VirtualGuestDB,
) *VirtualGuestController {
	return &VirtualGuestController{
		db:            db,
	}
}

func (h *VirtualGuestController) VirtualGuests(logger lager.Logger, publicVlan, privateVlan, cpu, memory_mb int32) ([]*models.VirtualGuest, error) {
	logger = logger.Session("vms")

	filter := models.VirtualGuestFilter{
		CPU: cpu,
		Memory_mb: memory_mb,
		PublicVlan: publicVlan,
		PrivateVlan: privateVlan,
	}

	return h.db.VirtualGuests(logger, filter)
}

func (h *VirtualGuestController) VirtualGuestByCid(logger lager.Logger, cid int32) (*models.VirtualGuest, error) {
	logger = logger.Session("vm-by-guid")
	return h.db.VirtualGuestByCID(logger, cid)
}