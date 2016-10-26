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

// VM response

type VMResponse struct {
	Error *Error          `json:"error,omitempty"`
	Vm   *VM 	      `json:"vm,omitempty"`
}

func (m *VMResponse) Reset()  { *m = VMResponse{} }
func (*VMResponse) ProtoMessage() {}
func (this *VMResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&VMResponse{`,
		`Error:` + strings.Replace(fmt.Sprintf("%v", this.Error), "Error", "Error", 1) + `,`,
		`Vm:` + strings.Replace(fmt.Sprintf("%v", this.Vm), "Vm", "Vm", 1) + `,`,
		`}`,
	}, "")
	return s
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
		`Vms:` + strings.Replace(fmt.Sprintf("%v", this.Vms), "Vms", "Vms", 1) + `,`,
		`}`,
	}, "")
	return s
}

//  VM life cycle response

type VMLifecycleResponse struct {
	Error *Error `json:"error,omitempty"`
}

func (m *VMLifecycleResponse) Reset()      { *m = VMLifecycleResponse{} }
func (*VMLifecycleResponse) ProtoMessage() {}

func (m *VMLifecycleResponse) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (this *VMLifecycleResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&VMLifecycleResponse{`,
		`Error:` + strings.Replace(fmt.Sprintf("%v", this.Error), "Error", "Error", 1) + `,`,
		`}`,
	}, "")
	return s
}