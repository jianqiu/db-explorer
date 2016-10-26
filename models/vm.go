package models

import (
	"fmt"
)

type State string

const (
	StateFree         State = "free"
	StateProvision    State = "provisioning"
	StateInUse        State = "using"
	StateUnknown      State = "unknown"
)

type VM struct {
	Hostname       string             `json:"hostname"`
	Cpu            int32              `json:"cpu"`
	MemoryMb       int32              `json:"memory_mb"`
	PublicVlan     int32              `json:"public_vlan"`
	PrivateVlan    int32              `json:"private_vlan"`
	State          State	  	  `json:"state"`
	Cid            int32              `json:"cid"`
	DeploymentName string             `json:"deployment_name"`
	Ip             string             `json:"ip"`
	CreatedAt      int64              `json:"created_at"`
	UpdatedAt      int64              `json:"updated_at"`
}

type StateValue struct {
	State       State  `json:"state"`
}

type VMFilter struct {
	CID int32
	IP string
	CPU int32
	Memory_mb int32
	PublicVlan int32
	PrivateVlan int32
	State	    State
}

func (vm *VM) Validate() error {
	var validationError ValidationError

	if vm.Cid == 0 {
		validationError = validationError.Append(ErrInvalidField{"cid"})
	}

	if !validationError.Empty() {
		return validationError
	}

	return nil
}

func (vm *VM) Copy() *VM {
	newVM := *vm
	return &newVM
}

func (t *VM) ValidateTransitionTo(to State) error {
	var valid bool
	from := t.State
	switch to {
	case StateFree:
		valid = ( from == StateInUse || from == StateProvision )
	case StateProvision:
		valid = ( from == StateFree || from == StateInUse )
	case StateUnknown:
		valid = true
	case StateInUse:
		valid = from == StateProvision
	}

	if !valid {
		return NewError(
			Error_InvalidStateTransition,
			fmt.Sprintf("Cannot transition from %v to %v", from, to),
		)
	}

	return nil
}
