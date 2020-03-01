// SPDX-License-Identifier: Apache-2.0

package firewalld

import (
	"fmt"

	"github.com/api-routerd/api-routerd/cmd/share"

	"github.com/godbus/dbus"
)

const (
	dbusInterface = "org.fedoraproject.FirewallD1"
	dbusPath      = "/org/fedoraproject/FirewallD1"

	dbusInterfaceConfig = "org.fedoraproject.FirewallD1.config"
	dbusPathConfig      = "/org/fedoraproject/FirewallD1/config"
)

// Conn dbus connection
type Conn struct {
	conn   *dbus.Conn
	object dbus.BusObject
}

// Zone FW Zone object
type Zone struct {
	Version     string   `json:"version"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Services    []string `json:"services"`
	Interfaces  []string `json:"interfaces"`
}

// Service FW Service object
type Service struct {
	Version      string            `json:"version"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Ports        [][]interface{}   `json:"ports"`
	Destinations map[string]string `json:"destinations"`
	Protocols    []string          `json:"protocols"`
	SourcePorts  [][]interface{}   `json:"source_ports"`
}

// NewConn get a new connection
func NewConn() (*Conn, error) {
	c := new(Conn)

	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to system bus: %v", err)
	}

	c.conn = conn
	c.object = conn.Object(dbusInterface, dbus.ObjectPath(dbusPath))

	return c, nil
}

// Close close the connection
func (c *Conn) Close() {
	c.conn.Close()
}

func (c *Conn) getZonePathbyName(zone string) (string, error) {
	var r string

	c.object = c.conn.Object(dbusInterface, dbus.ObjectPath(dbusPathConfig))
	err := c.object.Call(dbusInterface+".config.getZoneByName", 0, zone).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// GetZones gets all zones from runtime
func (c *Conn) GetZones() ([]string, error) {
	var z []string

	err := c.object.Call(dbusInterface+".zone.getZones", 0).Store(&z)
	if err != nil {
		return nil, err
	}

	return z, nil
}

// ListAllZones lists all zones from permanent config
func (c *Conn) ListAllZones() ([]string, error) {
	var z []string

	c.object = c.conn.Object(dbusInterface, dbus.ObjectPath(dbusPathConfig))
	err := c.object.Call(dbusInterfaceConfig+".listZones", 0).Store(&z)
	if err != nil {
		return nil, err
	}

	return z, nil
}

// ListServices lists all services from runtime
func (c *Conn) ListServices() ([]string, error) {
	var services []string

	err := c.object.Call(dbusInterface+".listServices", 0).Store(&services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

// GetDefaultZone gets the default runtime zone
func (c *Conn) GetDefaultZone() (string, error) {
	var zone string

	err := c.object.Call(dbusInterface+".getDefaultZone", 0).Store(&zone)
	if err != nil {
		return zone, err
	}

	return zone, nil
}

// ListPorts lists all ports from runtime
func (c *Conn) ListPorts(zone string) ([][]string, error) {
	var ports [][]string

	err := c.object.Call(dbusInterface+".zone.getPorts", 0, zone).Store(&ports)
	if err != nil {
		return nil, err
	}

	return ports, nil
}

// GetZoneSettings get all zone settings from runtime
func (c *Conn) GetZoneSettings(zone string) (*Zone, error) {
	out := []interface{}{}

	err := c.object.Call(dbusInterface+".getZoneSettings", 0, zone).Store(&out)
	if err != nil {
		return nil, err
	}

	z := new(Zone)
	for i, el := range out {
		switch i {
		case 1:
			z.Name = el.(string)
			break
		case 2:
			z.Description = el.(string)
			break
		case 5:
			z.Services = el.([]string)
			break
		case 10:
			z.Interfaces = el.([]string)
		}
	}

	return z, nil
}

// GetZoneSettingsPermanent get all zone settings from permament config
func (c *Conn) GetZoneSettingsPermanent(zone string) (*Zone, error) {
	out := []interface{}{}

	r, err := c.getZonePathbyName(zone)
	if err != nil {
		return nil, err
	}

	c.object = c.conn.Object(dbusInterface, dbus.ObjectPath(r))
	err = c.object.Call(dbusInterfaceConfig+".zone.getSettings", 0).Store(&out)
	if err != nil {
		return nil, err
	}

	z := new(Zone)
	for i, el := range out {
		switch i {
		case 1:
			z.Name = el.(string)
			break
		case 2:
			z.Description = el.(string)
			break
		case 5:
			z.Services = el.([]string)
			break
		case 10:
			z.Interfaces = el.([]string)
		}
	}

	return z, nil
}

// GetServiceSettings get service setting from runtime
func (c *Conn) GetServiceSettings(zone string) (*Service, error) {
	out := []interface{}{}

	err := c.object.Call(dbusInterface+".getServiceSettings", 0, zone).Store(&out)
	if err != nil {
		return nil, err
	}

	s := new(Service)
	for i, el := range out {
		switch i {
		case 1:
			s.Name = el.(string)
			break
		case 2:
			s.Description = el.(string)
			break
		case 3:
			s.Ports = el.([][]interface{})
			break
		case 5:
			s.Destinations = el.(map[string]string)
			break
		case 6:
			s.Protocols = el.([]string)
			break
		case 7:
			s.SourcePorts = el.([][]interface{})
			break
		}
	}

	return s, nil
}

// GetServiceSettingsPermanent gets all settings from permanent config
func (c *Conn) GetServiceSettingsPermanent(zone string) (*Service, error) {
	out := []interface{}{}

	r, err := c.getZonePathbyName(zone)
	if err != nil {
		return nil, err
	}

	c.object = c.conn.Object(dbusInterface, dbus.ObjectPath(r))
	err = c.object.Call(dbusInterfaceConfig+".getServiceSettings", 0).Store(&out)
	if err != nil {
		return nil, err
	}

	s := new(Service)
	for i, el := range out {
		switch i {
		case 1:
			s.Name = el.(string)
			break
		case 2:
			s.Description = el.(string)
			break
		case 3:
			s.Ports = el.([][]interface{})
			break
		case 5:
			s.Destinations = el.(map[string]string)
			break
		case 6:
			s.Protocols = el.([]string)
			break
		case 7:
			s.SourcePorts = el.([][]interface{})
			break
		}
	}

	return s, nil
}

//AddPort add a port to runtime
func (c *Conn) AddPort(zone string, port string, protocol string) (string, error) {
	var r string
	var err error

	err = c.object.Call(dbusInterface+".zone.addPort", 0, zone, port, protocol, 0).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// RemovePort remove a port from runtime
func (c *Conn) RemovePort(zone string, port string, protocol string) (string, error) {
	var r string

	err := c.object.Call(dbusInterface+".zone.removePort", 0, zone, port, protocol).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// AddPortPermanent add a port to permanent config
func (c *Conn) AddPortPermanent(zone string, port string, protocol string) (string, error) {
	r, err := c.getZonePathbyName(zone)
	if err != nil {
		return r, err
	}

	c.object = c.conn.Object(dbusInterface, dbus.ObjectPath(r))
	err = c.object.Call(dbusInterface+".config.zone.addPort", 0, port, protocol).Err
	if err != nil {
		return r, err
	}

	return r, nil
}

// RemovePortPermanent remove a port from permanent config
func (c *Conn) RemovePortPermanent(zone string, port string, protocol string) (string, error) {
	r, err := c.getZonePathbyName(zone)
	if err != nil {
		return r, err
	}

	c.object = c.conn.Object(dbusInterface, dbus.ObjectPath(r))
	err = c.object.Call(dbusInterface+".config.zone.removePort", 0, port, protocol).Err
	if err != nil {
		return r, err
	}

	return r, nil
}

// AddProtocol add a prortocol to runtime
func (c *Conn) AddProtocol(zone string, protocol string) (string, error) {
	var r string
	var err error

	err = c.object.Call(dbusInterface+".zone.addProtocol", 0, zone, protocol, 0).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// RemoveProtocol remove a protocol from runtime
func (c *Conn) RemoveProtocol(zone string, protocol string) (string, error) {
	var r string

	err := c.object.Call(dbusInterface+".zone.removeProtocol", 0, zone, protocol).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// AddProtocolPermanent add a protocol to permanent config
func (c *Conn) AddProtocolPermanent(zone string, protocol string) (string, error) {
	r, err := c.getZonePathbyName(zone)
	if err != nil {
		return r, err
	}

	c.object = c.conn.Object(dbusInterface, dbus.ObjectPath(r))
	err = c.object.Call(dbusInterfaceConfig+".zone.addProtocol", 0, protocol).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// RemoveProtocolPermanent removes a protcol from a zone in permanent configuration
func (c *Conn) RemoveProtocolPermanent(zone string, protocol string) (string, error) {
	r, err := c.getZonePathbyName(zone)
	if err != nil {
		return r, err
	}

	c.object = c.conn.Object(dbusInterface, dbus.ObjectPath(r))
	err = c.object.Call(dbusInterfaceConfig+".zone.removeProtocol", 0, protocol).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// AddInterface Add a interface to a zone in runtime configuration
func (c *Conn) AddInterface(zone string, intf string) (string, error) {
	var r string
	var err error

	err = c.object.Call(dbusInterface+".zone.addInterface", 0, zone, intf).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// RemoveInterface Remove a interface from a zone in runtime configuration
func (c *Conn) RemoveInterface(zone string, intf string) (string, error) {
	var r string

	err := c.object.Call(dbusInterface+".zone.removeInterface", 0, zone, intf).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// AddInterfacePermanent Add a interface to a zone in permanent configuration
func (c *Conn) AddInterfacePermanent(zone string, intf string) (string, error) {
	r, err := c.getZonePathbyName(zone)
	if err != nil {
		return r, err
	}

	c.object = c.conn.Object(dbusInterface, dbus.ObjectPath(r))
	err = c.object.Call(dbusInterfaceConfig+".zone.addInterface", 0, intf).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// RemoveInterfacePermanent Remove a interface from a zone in permanent configuration
func (c *Conn) RemoveInterfacePermanent(zone string, intf string) (string, error) {
	r, err := c.getZonePathbyName(zone)
	if err != nil {
		return r, err
	}

	c.object = c.conn.Object(dbusInterface, dbus.ObjectPath(r))
	err = c.object.Call(dbusInterfaceConfig+".zone.removeInterface", 0, intf).Store(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}
