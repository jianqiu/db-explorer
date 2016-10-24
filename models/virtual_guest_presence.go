package models

import (
	"fmt"
	"regexp"
)

var taskGuidPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type VirtualGuestFilter struct {
	CID int32
	IP string
	PublicVlan string
	PrivateVlan string
}

func (vg *VirtualGuest) Validate() error {
	var validationError ValidationError

	if vg.Cid == "" {
		validationError = validationError.Append(ErrInvalidField{"cid"})
	}

	if !taskGuidPattern.MatchString(vg.Cid) {
		validationError = validationError.Append(ErrInvalidField{"task_guid"})
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
	case VirtualGuest_Using:
		valid = from == VirtualGuest_Deleted
	case VirtualGuest_Deleted:
		valid = (from == VirtualGuest_Deleted || from == VirtualGuest_Using)
	case VirtualGuest_Unavailable:
		valid = from == VirtualGuest_Deleted
	}

	if !valid {
		return NewError(
			Error_InvalidStateTransition,
			fmt.Sprintf("Cannot transition from %s to %s", from.String(), to.String()),
		)
	}

	return nil
}
