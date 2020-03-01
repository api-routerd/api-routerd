// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"net/http"
	"strings"

	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

const (
	sudoersPath   = "/etc/sudoers"
	sshConfigFile = "/etc/ssh/sshd_config"
)

// SudoersConf Json response
type SudoersConf struct {
	Sudoers []string `json:"sudoers"`
}

// SSHdConf sshd json response
type SSHdConf struct {
	Value string `json:"value"`
}

// GetSudoers read sudoers file
func GetSudoers(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(sudoersPath)
	if err != nil {
		return err
	}

	var sudoers []string
	for _, line := range lines {
		if strings.Contains(line, "%") || strings.Contains(line, "Defaults") || strings.Contains(line, "root") {
			continue
		}

		line = strings.Replace(line, "\t", " ", 5)
		sudoers = append(sudoers, line)
	}

	sshdconf := SudoersConf{Sudoers: sudoers}

	return share.JSONResponse(sshdconf, rw)
}

// SSHConfFileRead read sshd configuration file
func SSHConfFileRead(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(sshConfigFile)
	if err != nil {
		log.Errorf("Failed to read: %s, %v", sshConfigFile, err)
		return err
	}

	sshdConfMap := make(map[string]string)
	for _, line := range lines {
		fields := strings.Fields(line)
		paramName := fields[0]
		sshdConfMap[paramName] = fields[1]
	}

	return share.JSONResponse(sshdConfMap, rw)
}
