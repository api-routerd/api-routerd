// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

const (
	networkdUnitPath = "/etc/systemd/network"
)

// Match section
type Match struct {
	MAC    string `json:"MAC"`
	Driver string `json:"Driver"`
	Name   string `json:"Name"`
}

// InitNetworkd init networkd module
func InitNetworkd() error {
	err := share.CreateDirectory(networkdUnitPath, 0777)
	if err != nil {
		log.Errorf("Failed create network unit path %s: %v", networkdUnitPath, err)
		return err
	}

	return nil
}
