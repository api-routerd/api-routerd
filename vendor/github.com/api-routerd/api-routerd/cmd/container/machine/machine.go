// SPDX-License-Identifier: Apache-2.0

package machine

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/api-routerd/api-routerd/cmd/share"
	"github.com/coreos/go-systemd/machine1"
)

// Machine Json request
type Machine struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
	Old      string `json:"old"`
	New      string `json:"new"`
}

var machineMethods *share.Set

// MethodGet retrives info from machined via dbus
func (m *Machine) MethodGet(rw http.ResponseWriter) error {
	c, err := machine1.New()
	if err != nil {
		return err
	}

	b := machineMethods.Contains(m.Path)
	if !b {
		return fmt.Errorf("Failed to call method machine: %s not found", m.Path)
	}

	switch m.Path {
	case "list-images":
		images, err := c.ListImages()
		if err != nil {
			return err
		}

		return share.JSONResponse(images, rw)

	case "list-machines":
		machines, err := c.ListMachines()
		if err != nil {
			return err
		}

		return share.JSONResponse(machines, rw)

	case "get-machine":
		machine, err := c.GetMachine(m.Property)
		if err != nil {
			return err
		}

		return share.JSONResponse(machine, rw)
	case "describe-machine":
		machine, err := c.GetMachine(m.Property)
		if err != nil {
			return err
		}

		return share.JSONResponse(machine, rw)

	case "get-image":
		image, err := c.GetImage(m.Property)
		if err != nil {
			return err
		}

		return share.JSONResponse(image, rw)
	case "get-machine-by-pid":
		pid, err := strconv.ParseInt(m.Property, 10, 32)
		if err != nil {
			return nil
		}

		p, err := c.GetMachineByPID(uint(pid))
		if err != nil {
			return err
		}

		return share.JSONResponse(p, rw)
	case "get-machine-address":
		addr, err := c.GetMachineAddresses(m.Property)
		if err != nil {
			return err
		}

		return share.JSONResponse(addr, rw)
	case "get-machine-osrelease":
		addr, err := GetMachineOSRelease(m.Property)
		if err != nil {
			return err
		}

		return share.JSONResponse(addr, rw)
	}

	return nil
}

// MethodConfigure Post methods
func (m *Machine) MethodConfigure(rw http.ResponseWriter) error {
	c, err := machine1.New()
	if err != nil {
		return err
	}

	b := machineMethods.Contains(m.Path)
	if !b {
		return fmt.Errorf("Failed to call method machine: %s not found", m.Path)
	}

	switch m.Path {
	case "terminate-machine":
		err := c.TerminateMachine(m.Property)
		if err != nil {
			return err
		}

		return nil
	case "clone-image":
		err := CloneImage(m.Old, m.New)
		if err != nil {
			return err
		}

		return nil
	case "rename-image":
		err := RenameImage(m.Old, m.New)
		if err != nil {
			return err
		}

		return nil
	case "remove-image":
		err := RemoveImage(m.Old)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

// InitMachine init machine package
func InitMachine() error {
	machineMethods = share.NewSet()

	machineMethods.Add("list-images")
	machineMethods.Add("list-machines")
	machineMethods.Add("get-machine")
	machineMethods.Add("get-image")
	machineMethods.Add("get-machine-by-pid")
	machineMethods.Add("get-machine-address")
	machineMethods.Add("describe-machine")
	machineMethods.Add("terminate-machine")
	machineMethods.Add("get-machine-osrelease")
	machineMethods.Add("rename-image")
	machineMethods.Add("remove-image")

	return nil
}
