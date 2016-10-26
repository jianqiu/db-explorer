package controllers

import (
	"code.cloudfoundry.org/lager"
	vps "github.com/jianqiu/vm-pool-server"
	"github.com/jianqiu/vm-pool-server/db"
	"github.com/jianqiu/vm-pool-server/models"
)

type VirtualGuestController struct {
	db            db.VirtualGuestDB
	serviceClient vps.ServiceClient
}

func NewVirtualGuestController(
	db db.VirtualGuestDB,
	serviceClient vps.ServiceClient,
) *VirtualGuestController {
	return &VirtualGuestController{
		db:            db,
		serviceClient: serviceClient,
	}
}

func (h *VirtualGuestController) AllVirtualGuests(logger lager.Logger) ([]*models.VM, error) {
	logger = logger.Session("vms")

	filter := models.VMFilter{}

	return h.db.VirtualGuests(logger, filter)
}

func (h *VirtualGuestController) VirtualGuests(logger lager.Logger, publicVlan, privateVlan, cpu, memory_mb int32) ([]*models.VM, error) {
	logger = logger.Session("vms")

	filter := models.VMFilter{
		CPU: cpu,
		Memory_mb: memory_mb,
		PublicVlan: publicVlan,
		PrivateVlan: privateVlan,
	}

	return h.db.VirtualGuests(logger, filter)
}

func (h *VirtualGuestController) CreateVM(logger lager.Logger, vmDefinition *models.VM) error {
	var err error
	logger = logger.Session("create-vm")

	err = h.db.InsertVirtualGuestToPool(logger, vmDefinition)
	if err != nil {
		return err
	}

	return nil
}

func (h *VirtualGuestController) DeleteVM(logger lager.Logger, cid int32) error {
	logger = logger.Session("delete-vm")

	return h.db.DeleteVirtualGuestFromPool(logger, cid)
}

func (h *VirtualGuestController) UpdateVM(logger lager.Logger, cid int32, updateData *models.StateValue) error {
	var err error
	logger = logger.Session("update-vm")

	switch updateData.State {
	case models.StateInUse:
		err = h.db.ChangeVirtualGuestToUse(logger, cid)
	case models.StateFree:
		err = h.db.ChangeVirtualGuestToFree(logger, cid)
	case models.StateProvision:
		err = h.db.ChangeVirtualGuestToProvision(logger, cid)
	}

	if err != nil {
		return err
	}

	return nil
}

func (h *VirtualGuestController) VirtualGuestByCid(logger lager.Logger, cid int32) (*models.VM, error) {
	logger = logger.Session("vm-by-guid")
	return h.db.VirtualGuestByCID(logger, cid)
}
