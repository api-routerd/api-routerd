// SPDX-License-Identifier: Apache-2.0

package firewalld

import (
	"fmt"
	"net/http"

	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

// Firewall Json request
type Firewall struct {
	Property string `json:"property"`
	Value    string `json:"value"`

	Zone      string `json:"zone"`
	Port      string `json:"port"`
	Protocol  string `json:"protocol"`
	Interface string `json:"interface"`
	Permanent bool   `json:"permanent,omitempty"`
}

var firewalldMethods *share.Set

// GetFirewalld wraps all FW get commands
func (f *Firewall) GetFirewalld(rw http.ResponseWriter) error {
	c, err := NewConn()
	if err != nil {
		return err
	}
	defer c.Close()

	b := firewalldMethods.Contains(f.Property)
	if !b {
		return fmt.Errorf("Failed to call method firewalld: %s not found", f.Property)
	}

	log.Debugf("Get Firewalld passthrough: %s", f.Property)

	fmt.Println(f)
	switch f.Property {
	case "list-services":
		services, err := c.ListServices()
		if err != nil {
			return err
		}

		return share.JSONResponse(services, rw)

	case "get-zones":
		z, err := c.GetZones()
		if err != nil {
			return err
		}

		return share.JSONResponse(z, rw)

	case "list-all-zones":
		z, err := c.ListAllZones()
		if err != nil {
			return err
		}

		return share.JSONResponse(z, rw)

	case "list-ports":
		z, err := c.ListPorts(f.Value)
		if err != nil {
			return err
		}

		return share.JSONResponse(z, rw)

	case "get-default-zone":
		z, err := c.GetDefaultZone()
		if err != nil {
			return err
		}

		return share.JSONResponse(z, rw)

	case "get-zone-settings":
		z, err := c.GetZoneSettings(f.Value)
		if err != nil {
			return err
		}

		return share.JSONResponse(z, rw)

	case "get-zone-settings-permanent":
		z, err := c.GetZoneSettingsPermanent(f.Value)
		if err != nil {
			return err
		}

		return share.JSONResponse(z, rw)

	case "get-service-settings":
		z, err := c.GetServiceSettings(f.Value)
		if err != nil {
			return err
		}

		return share.JSONResponse(z, rw)
	case "get-service-settings-permanent":
		z, err := c.GetServiceSettings(f.Value)
		if err != nil {
			return err
		}

		return share.JSONResponse(z, rw)
	}

	return nil
}

// AddFirewalld wrap all firewalld add
func (f *Firewall) AddFirewalld(rw http.ResponseWriter) error {
	c, err := NewConn()
	if err != nil {
		return err
	}
	defer c.Close()

	b := firewalldMethods.Contains(f.Property)
	if !b {
		return fmt.Errorf("Failed to call method firewalld: %s not found", f.Property)
	}

	log.Debugf("Set Firewalld passthrough: %s", f.Property)

	switch f.Property {
	case "add-port":
		var r string

		if f.Permanent == true {
			r, err = c.AddPortPermanent(f.Zone, f.Port, f.Protocol)
		} else {
			r, err = c.AddPort(f.Zone, f.Port, f.Protocol)
		}

		if err != nil {
			return err
		}

		return share.JSONResponse(r, rw)
	case "add-protocol":
		var r string

		if f.Permanent == true {
			r, err = c.AddProtocolPermanent(f.Zone, f.Protocol)
		} else {
			r, err = c.AddProtocol(f.Zone, f.Protocol)
		}

		if err != nil {
			return err
		}

		return share.JSONResponse(r, rw)
	case "add-interface":
		var r string

		if f.Permanent == true {
			r, err = c.AddInterfacePermanent(f.Zone, f.Interface)
		} else {
			r, err = c.AddInterface(f.Zone, f.Interface)
		}

		if err != nil {
			return err
		}

		return share.JSONResponse(r, rw)
	}

	return nil
}

// DeleteFirewalld wrap all delete commands
func (f *Firewall) DeleteFirewalld(rw http.ResponseWriter) error {
	c, err := NewConn()
	if err != nil {
		return err
	}
	defer c.Close()

	b := firewalldMethods.Contains(f.Property)
	if !b {
		return fmt.Errorf("Failed to call method firewalld: %s not found", f.Property)
	}

	log.Debugf("Delete Firewalld passthrough: %s", f.Property)

	switch f.Property {
	case "remove-port":
		var r string

		if f.Permanent == true {
			r, err = c.RemovePortPermanent(f.Zone, f.Port, f.Protocol)
		} else {
			r, err = c.RemovePort(f.Zone, f.Port, f.Protocol)
		}

		if err != nil {
			return err
		}

		return share.JSONResponse(r, rw)
	case "remove-protocol":
		var r string

		if f.Permanent == true {
			r, err = c.RemoveProtocolPermanent(f.Zone, f.Protocol)
		} else {
			r, err = c.RemoveProtocol(f.Zone, f.Protocol)
		}

		if err != nil {
			return err
		}

		return share.JSONResponse(r, rw)
	case "remove-interface":
		var r string

		if f.Permanent == true {
			r, err = c.RemoveInterfacePermanent(f.Zone, f.Interface)
		} else {
			r, err = c.RemoveInterface(f.Zone, f.Interface)
		}

		if err != nil {
			return err
		}

		return share.JSONResponse(r, rw)
	}

	return nil
}

// Init Init the FW module
func Init() error {
	firewalldMethods = share.NewSet()

	firewalldMethods.Add("list-services")
	firewalldMethods.Add("get-zones")
	firewalldMethods.Add("list-all-zones")
	firewalldMethods.Add("list-ports")
	firewalldMethods.Add("get-default-zone")
	firewalldMethods.Add("get-zone-settings")
	firewalldMethods.Add("get-zone-settings-permanent")
	firewalldMethods.Add("get-service-settings")
	firewalldMethods.Add("get-service-settings-permanent")
	firewalldMethods.Add("add-port")
	firewalldMethods.Add("remove-port")
	firewalldMethods.Add("add-protocol")
	firewalldMethods.Add("remove-protocol")
	firewalldMethods.Add("add-interface")
	firewalldMethods.Add("remove-interface")

	return nil
}
