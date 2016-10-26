package sqldb

import (
	"database/sql"
	"strings"

	"github.com/jianqiu/vm-pool-server/models"
	"code.cloudfoundry.org/lager"
)

func (db *SQLDB) VirtualGuests(logger lager.Logger, filter models.VMFilter) ([]*models.VM, error) {
	logger = logger.Session("virtualguests", lager.Data{"filter": filter})
	logger.Debug("starting")
	defer logger.Debug("complete")

	wheres := []string{}
	values := []interface{}{}

	if filter.CID != 0 {
		wheres = append(wheres, "cid = ?")
		values = append(values, filter.CID)
	}

	if filter.PrivateVlan != 0 {
		wheres = append(wheres, "private_vlan = ?")
		values = append(values, filter.PrivateVlan)
	}

	if filter.PublicVlan != 0 {
		wheres = append(wheres, "public_vlan = ?")
		values = append(values, filter.PublicVlan)
	}

	rows, err := db.all(logger, db.db, virtualGuests,
		virtualGuestColumns, LockRow,
		strings.Join(wheres, " AND "), values...,
	)
	if err != nil {
		logger.Error("failed-query", err)
		return nil, db.convertSQLError(err)
	}
	defer rows.Close()

	results := []*models.VM{}
	for rows.Next() {
		vm, err := db.fetchVirtualGuest(logger, rows, db.db)
		if err != nil {
			logger.Error("failed-fetch", err)
			return nil, err
		}
		results = append(results, vm)
	}

	if rows.Err() != nil {
		logger.Error("failed-getting-next-row", rows.Err())
		return nil, db.convertSQLError(rows.Err())
	}

	return results, nil
}

func (db *SQLDB) VirtualGuestByCID(logger lager.Logger, cid int32) (*models.VM, error) {
	logger = logger.Session("virtual-by-cid", lager.Data{"cid": cid})
	logger.Debug("starting")
	defer logger.Debug("complete")

	row := db.one(logger, db.db, virtualGuests,
		virtualGuestColumns, NoLockRow,
		"cid = ?", cid,
	)
	return db.fetchVirtualGuest(logger, row, db.db)
}

func (db *SQLDB) VirtualGuestByIP(logger lager.Logger, ip string) (*models.VM, error) {
	logger = logger.Session("virtual-by-ip", lager.Data{"ip": ip})
	logger.Debug("starting")
	defer logger.Debug("complete")

	row := db.one(logger, db.db, virtualGuests,
		virtualGuestColumns, NoLockRow,
		"ip = ?", ip,
	)
	return db.fetchVirtualGuest(logger, row, db.db)
}

func (db *SQLDB) InsertVirtualGuestToPool(logger lager.Logger, virtualGuest *models.VM) error {
	logger = logger.Session("insert-virtual-guest-to-pool", lager.Data{"cid": virtualGuest.Cid})
	logger.Info("starting")
	defer logger.Info("complete")

	now := db.clock.Now().UnixNano()

	_, err := db.insert(logger, db.db, virtualGuests,
		SQLAttributes{
			"cid":               virtualGuest.Cid,
			"hostname":          virtualGuest.Hostname,
			"ip":		     virtualGuest.Ip,
			"cpu":		     virtualGuest.Cpu,
			"memory_mb":         virtualGuest.MemoryMb,
			"public_vlan":	     virtualGuest.PublicVlan,
			"private_vlan":      virtualGuest.PrivateVlan,
			"create_at":          now,
			"updated_at":         now,
			"deployment_name":    virtualGuest.DeploymentName,
			"state":              "free",
		},
	)
	if err != nil {
		logger.Error("failed-inserting-virtual-guest", err)
		return db.convertSQLError(err)
	}

	return nil
}

func (db *SQLDB) ChangeVirtualGuestToProvision(logger lager.Logger, cid int32) (bool, error) {
	logger = logger.Session("update-virtual-guest-to-in-use", lager.Data{"cid": cid})

	var started bool

	err := db.transact(logger, func(logger lager.Logger, tx *sql.Tx) error {
		task, err := db.fetchTaskForUpdate(logger, cid, tx)
		if err != nil {
			logger.Error("failed-locking-virtual-guest", err)
			return err
		}

		if err = task.ValidateTransitionTo(models.StateProvision); err != nil {
			logger.Error("failed-to-transition-task-to-running", err)
			return err
		}

		logger.Info("starting")
		defer logger.Info("complete")
		now := db.clock.Now().UnixNano()
		_, err = db.update(logger, tx, virtualGuests,
			SQLAttributes{
				"state":      "provisioning",
				"updated_at": now,
			},
			"cid = ?", cid,
		)
		if err != nil {
			return db.convertSQLError(err)
		}

		started = true
		return nil
	})

	return started, err
}

func (db *SQLDB) ChangeVirtualGuestToUse(logger lager.Logger, cid int32) (bool, error) {
	logger = logger.Session("update-virtual-guest-to-in-use", lager.Data{"cid": cid})

	var started bool

	err := db.transact(logger, func(logger lager.Logger, tx *sql.Tx) error {
		task, err := db.fetchTaskForUpdate(logger, cid, tx)
		if err != nil {
			logger.Error("failed-locking-virtual-guest", err)
			return err
		}

		if err = task.ValidateTransitionTo(models.StateInUse); err != nil {
			logger.Error("failed-to-transition-task-to-running", err)
			return err
		}

		logger.Info("starting")
		defer logger.Info("complete")
		now := db.clock.Now().UnixNano()
		_, err = db.update(logger, tx, virtualGuests,
			SQLAttributes{
				"state":      "using",
				"updated_at": now,
			},
			"cid = ?", cid,
		)
		if err != nil {
			return db.convertSQLError(err)
		}

		started = true
		return nil
	})

	return started, err
}

func (db *SQLDB) ChangeVirtualGuestToFree(logger lager.Logger, cid int32) (bool, error) {
	logger = logger.Session("update-virtual-guest-to-deleted", lager.Data{"cid": cid})

	var started bool

	err := db.transact(logger, func(logger lager.Logger, tx *sql.Tx) error {
		task, err := db.fetchTaskForUpdate(logger, cid, tx)
		if err != nil {
			logger.Error("failed-locking-virtual-guest", err)
			return err
		}

		if err = task.ValidateTransitionTo(models.StateFree); err != nil {
			logger.Error("failed-to-transition-task-to-deleted", err)
			return err
		}

		logger.Info("starting")
		defer logger.Info("complete")
		now := db.clock.Now().UnixNano()
		_, err = db.update(logger, tx, virtualGuests,
			SQLAttributes{
				"state":      "free",
				"updated_at": now,
			},
			"cid = ?", cid,
		)
		if err != nil {
			return db.convertSQLError(err)
		}

		started = true
		return nil
	})

	return started, err
}

func (db *SQLDB) DeleteVirtualGuestFromPool(logger lager.Logger, cid int32) error {
	logger = logger.Session("delete-virtual-guest-from-pool", lager.Data{"cid": cid})
	logger.Info("starting")
	defer logger.Info("complete")

	return db.transact(logger, func(logger lager.Logger, tx *sql.Tx) error {
		task, err := db.fetchTaskForUpdate(logger, cid, tx)
		if err != nil {
			logger.Error("failed-locking-virtual-guest", err)
			return err
		}

		if task.State != models.StateFree {
			err = models.NewTaskTransitionError(task.State, models.StateUnknown)
			logger.Error("invalid-state-transition", err)
			return err
		}

		_, err = db.delete(logger, tx, virtualGuests, "cid = ?", cid)
		if err != nil {
			logger.Error("failed-deleting-virtual-guest", err)
			return db.convertSQLError(err)
		}

		return nil
	})
}

func (db *SQLDB) fetchTaskForUpdate(logger lager.Logger, cid int32, tx *sql.Tx) (*models.VM, error) {
	row := db.one(logger, tx, virtualGuests,
		virtualGuestColumns, LockRow,
		"cid = ?", cid,
	)
	return db.fetchVirtualGuest(logger, row, tx)
}

func (db *SQLDB) fetchVirtualGuest(logger lager.Logger, scanner RowScanner, tx Queryable) (*models.VM, error) {
	var hostname, ip, deployment_name, state string
	var cpu, memory_mb, cid, public_vlan, private_vlan int32
	err := scanner.Scan(
		&cid,
		&hostname,
		&ip,
		&cpu,
		&memory_mb,
		&private_vlan,
		&public_vlan,
		&deployment_name,
		&state,
	)
	if err != nil {
		logger.Error("failed-scanning-row", err)
		return nil, models.ErrResourceNotFound
	}

	virtualGuest := &models.VM{
		Cid:              cid,
		Hostname:         hostname,
		Ip:               ip,
		Cpu:    	  cpu,
		MemoryMb: 	  memory_mb,
		PrivateVlan:      private_vlan,
		PublicVlan:       public_vlan,
		DeploymentName:   deployment_name,
	}
	switch state {
	case "free":
		virtualGuest.State = models.StateFree
	case "provisioning":
		virtualGuest.State = models.StateProvision
	case "using":
		virtualGuest.State = models.StateInUse
	default:
		virtualGuest.State = models.StateUnknown
	}

	return virtualGuest, nil
}
