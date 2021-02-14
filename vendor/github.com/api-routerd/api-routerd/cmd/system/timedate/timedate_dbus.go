// SPDX-License-Identifier: Apache-2.0

package timedate

import (
	"fmt"

	"github.com/api-routerd/api-routerd/cmd/share"

	"github.com/godbus/dbus"
)

const (
	dbusInterface = "org.freedesktop.timedate1"
	dbusPath      = "/org/freedesktop/timedate1"
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

// SetTimeDate set timdate properties
func (c *Conn) SetTimeDate(property string, value string) error {
	var err error

	if property == "SetNTP" {
		b, err := share.ParseBool(value)
		if err != nil {
			return err
		}

		err = c.object.Call(dbusInterface+"."+property, 0, b, false).Err
	} else {
		err = c.object.Call(dbusInterface+"."+property, 0, value, false).Err
	}

	return err
}

// GetTimeDate get timedate properties
func (c *Conn) GetTimeDate(property string) (dbus.Variant, error) {
	p, err := c.object.GetProperty(dbusInterface + "." + property)
	if err != nil {
		return dbus.Variant{}, err
	}

	return p, nil
}
