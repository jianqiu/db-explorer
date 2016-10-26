package models

import (
	"encoding/json"
)

type VMRequestFilter struct {
	PublicVlan  int32 `json:"public_vlan"`
	PrivateVlan int32 `json:"private_vlan"`
	Cpu         int32 `json:"cpu"`
	MemoryMb    int32 `json:"memory_mb"`
}

func (req *VMRequestFilter) Unmarshal(data []byte) error {
	return json.Unmarshal(data, req)
}

func (req *VMRequestFilter) Validate() error {
	return nil
}

func (request *VMByCidRequest) Validate() error {
	var validationError ValidationError

	if request.Cid == 0 {
		validationError = validationError.Append(ErrInvalidField{"cid"})
	}

	if !validationError.Empty() {
		return validationError
	}

	return nil
}
