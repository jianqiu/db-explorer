package db

import (
	"github.com/jianqiu/vm-pool-server/models"
	"code.cloudfoundry.org/lager"
)

//go:generate counterfeiter . VirtualGuestDB
type VirtualGuestDB interface {
	VirtualGuests(logger lager.Logger, filter models.VirtualGuestFilter) ([]*models.VirtualGuest, error)
	VirtualGuestByCID(logger lager.Logger, cid string) (*models.VirtualGuest, error)
	VirtualGuestByIP(logger lager.Logger, ip string) (*models.VirtualGuest, error)

	InsertVirtualGuestToPool(logger lager.Logger,virtualGuest *models.VirtualGuest) (error)
	ChangeVirtualGuestToUse(logger lager.Logger, cid int32) (bool, error)
	ChangeVirtualGuestToDeleted(logger lager.Logger, cid int32) (bool, error)
	DeleteVirtualGuestFromPool(logger lager.Logger, cid int32) (bool, error)
}

