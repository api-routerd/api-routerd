// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"fmt"
	"net/http"

	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

// Hostname commands
type Hostname struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

var hostNameMethods *share.Set
var hostNameInfo = map[string]string{}

// SetHostname set hostname via dbus
func (h *Hostname) SetHostname() error {
	conn, err := NewConn()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	b := hostNameMethods.Contains(h.Property)
	if !b {
		return fmt.Errorf("Failed to set hostname property: %s not found", h.Property)
	}

	return conn.SetHostName(h.Property, h.Value)
}

// GetHostname retrieves properties from hostnamed via dbus
func GetHostname(rw http.ResponseWriter, property string) error {
	conn, err := NewConn()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	for k := range hostNameInfo {
		p, err := conn.GetHostName(k)
		if err != nil {
			log.Errorf("Failed to get org.freedesktop.hostname1.%s", k)
			continue
		}

		hostNameInfo[k] = p
	}

	if property == "" {
		return share.JSONResponse(hostNameInfo, rw)
	}

	host := Hostname{
		Property: property,
		Value:    hostNameInfo[property],
	}

	return share.JSONResponse(host, rw)
}

// InitHostName init hostname package
func InitHostname() error {
	hostNameMethods = share.NewSet()

	hostNameMethods.Add("SetHostname")
	hostNameMethods.Add("SetStaticHostname")
	hostNameMethods.Add("SetPrettyHostname")
	hostNameMethods.Add("SetIconName")
	hostNameMethods.Add("SetChassis")
	hostNameMethods.Add("SetDeployment")
	hostNameMethods.Add("SetLocation")

	hostNameInfo["Hostname"] = ""
	hostNameInfo["StaticHostname"] = ""
	hostNameInfo["PrettyHostname"] = ""
	hostNameInfo["IconName"] = ""
	hostNameInfo["Chassis"] = ""
	hostNameInfo["Deployment"] = ""
	hostNameInfo["Location"] = ""
	hostNameInfo["KernelName"] = ""
	hostNameInfo["KernelRelease"] = ""
	hostNameInfo["KernelVersion"] = ""
	hostNameInfo["OperatingSystemPrettyName"] = ""
	hostNameInfo["OperatingSystemCPEName"] = ""
	hostNameInfo["HomeURL"] = ""

	return nil
}
