package models

import (
	"fmt"
	"strings"
	"encoding/json"
)

// VM filter request

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

// VMs response

type VMsResponse struct {
	Error *Error          `json:"error,omitempty"`
	Vms   []*VM 	      `json:"vms,omitempty"`
}

func (m *VMsResponse) Reset()  { *m = VMsResponse{} }
func (*VMsResponse) ProtoMessage() {}
func (this *VMsResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&VMsResponse{`,
		`Error:` + strings.Replace(fmt.Sprintf("%v", this.Error), "Error", "Error", 1) + `,`,
		`Vms:` + strings.Replace(fmt.Sprintf("%v", this.Vms), "VirtualGuest", "VirtualGuest", 1) + `,`,
		`}`,
	}, "")
	return s
}