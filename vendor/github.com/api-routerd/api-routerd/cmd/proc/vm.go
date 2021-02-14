// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"net/http"
	"path"

	"github.com/api-routerd/api-routerd/cmd/share"
)

const (
	vmPath = "/proc/sys/vm"
)

// VM JSON message
type VM struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

// GetVM read proc vm property
func (r *VM) GetVM(rw http.ResponseWriter) error {
	line, err := share.ReadOneLineFile(path.Join(vmPath, r.Property))
	if err != nil {
		return err
	}

	vm := VM{
		Property: r.Property,
		Value:    line,
	}

	return share.JSONResponse(vm, rw)
}

// SetVM write a value to VM
func (r *VM) SetVM(rw http.ResponseWriter) error {
	err := share.WriteOneLineFile(path.Join(vmPath, r.Property), r.Value)
	if err != nil {
		return err
	}

	line, err := share.ReadOneLineFile(path.Join(vmPath, r.Property))
	if err != nil {
		return err
	}

	vm := VM{
		Property: r.Property,
		Value:    line,
	}

	return share.JSONResponse(vm, rw)
}
