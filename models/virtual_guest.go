package models

import (
	"fmt"
)

type VirtualGuestFilter struct {
	CID int32
	IP string
	CPU int32
	Memory_mb int32
	PublicVlan int32
	PrivateVlan int32
}

func (vg *VirtualGuest) Validate() error {
	var validationError ValidationError

	if vg.Cid == 0 {
		validationError = validationError.Append(ErrInvalidField{"cid"})
	}

	if !validationError.Empty() {
		return validationError
	}

	return nil
}

func (vg *VirtualGuest) Copy() *VirtualGuest {
	newTask := *vg
	return &newTask
}

func (t *VirtualGuest) ValidateTransitionTo(to VirtualGuest_State) error {
	var valid bool
	from := t.State
	switch to {
	case Using:
		valid = from == Deleted
	case Deleted:
		valid = (from == Deleted || from == Using)
	case Unavailable:
		valid = from == Deleted
	}

	if !valid {
		return NewError(
			Error_InvalidStateTransition,
			fmt.Sprintf("Cannot transition from %s to %s", from.String(), to.String()),
		)
	}

	return nil
}
