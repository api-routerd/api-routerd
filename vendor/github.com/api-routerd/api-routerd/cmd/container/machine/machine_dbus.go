// SPDX-License-Identifier: Apache-2.0

package machine

import (
	"fmt"

	"github.com/api-routerd/api-routerd/cmd/share"
	"github.com/godbus/dbus"
)

const (
	dbusInterface = "org.freedesktop.machine1.Manager"
	dbusPath      = "/org/freedesktop/machine1"
)

// Conn connection object
type Conn struct {
	conn   *dbus.Conn
	object dbus.BusObject
}

// NewConn opens a new dbus connection
func NewConn() (*Conn, error) {
	c := new(Conn)

	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		return nil, err
	}

	c.conn = conn
	c.object = c.conn.Object("org.freedesktop.machine1", dbus.ObjectPath(dbusPath))

	return c, nil
}

// Close close a dbus connection
func (c *Conn) Close() {
	c.conn.Close()
}

// GetMachineOSRelease retrive the OSRelease info
func GetMachineOSRelease(machine string) (dbus.ObjectPath, error) {
	conn, err := NewConn()
	if err != nil {
		return "", fmt.Errorf("Failed to get systemd bus connection: %v", err)
	}
	defer conn.Close()

	r := conn.object.Call(fmt.Sprintf("%s.%s", dbusInterface, "GetMachineOSRelease"), 0, machine)
	if r != nil {
		fmt.Println(r)
		return "", fmt.Errorf("Failed to get machine release information: %v", r.Err)
	}

	path, typeErr := r.Body[0].(dbus.ObjectPath)
	if !typeErr {
		return "", fmt.Errorf("unable to convert dbus response '%v' to dbus.ObjectPath", r.Body[0])
	}

	return path, nil
}

// CloneImage clones a image
func CloneImage(image string, newImage string) error {
	conn, err := NewConn()
	if err != nil {
		return fmt.Errorf("Failed to get systemd bus connection: %v", err)
	}
	defer conn.Close()

	r := conn.object.Call(fmt.Sprintf("%s.%s", dbusInterface, "CloneImage"), 0, image, newImage, false)
	if r != nil {
		fmt.Println(r)
		return fmt.Errorf("Failed to clone image: %v", r.Err)
	}

	return nil
}

// RenameImage rename a image
func RenameImage(image string, newImage string) error {
	conn, err := NewConn()
	if err != nil {
		return fmt.Errorf("Failed to get systemd bus connection: %s", err)
	}
	defer conn.Close()

	r := conn.object.Call(fmt.Sprintf("%s.%s", dbusInterface, "RenameImage"), 0, image, newImage)
	if r != nil {
		fmt.Println(r)
		return fmt.Errorf("Failed to Rename Image : %v", r.Err)
	}

	return nil
}

// RemoveImage remove a image
func RemoveImage(image string) error {
	conn, err := NewConn()
	if err != nil {
		return fmt.Errorf("Failed to get systemd bus connection: %s", err)
	}
	defer conn.Close()

	r := conn.object.Call(fmt.Sprintf("%s.%s", dbusInterface, "RemoveImage"), 0, image)
	if r != nil {
		fmt.Println(r)
		return fmt.Errorf("Failed to Remove Image : %v", r.Err)
	}

	return nil
}
