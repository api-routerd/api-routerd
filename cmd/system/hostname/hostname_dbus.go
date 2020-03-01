// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"fmt"

	"github.com/api-routerd/api-routerd/cmd/share"

	"github.com/godbus/dbus"
)

const (
	dbusInterface = "org.freedesktop.hostname1"
	dbusPath      = "/org/freedesktop/hostname1"
)

// Conn dbus connection
type Conn struct {
	conn   *dbus.Conn
	object dbus.BusObject
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

// SetHostName set hostname properties
func (c *Conn) SetHostName(property string, value string) error {
	err := c.object.Call(dbusInterface+"."+property, 0, value, false).Err
	if err != nil {
		return fmt.Errorf("Failed to set hostname: %v", err)
	}

	return nil
}

// GetHostName get hostname properties
func (c *Conn) GetHostName(property string) (string, error) {
	p, err := c.object.GetProperty(dbusInterface + "." + property)
	if err != nil {
		return "", err
	}

	v, b := p.Value().(string)
	if !b {
		return "", fmt.Errorf("Empty value received: %s", property)
	}

	return v, nil
}
