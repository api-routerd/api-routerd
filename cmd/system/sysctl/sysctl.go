// SPDX-License-Identifier: Apache-2.0

package sysctl

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

const (
	sysctlPath = "/etc/sysctl.conf"
)

//Sysctl json request
type Sysctl struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Apply string `json:"apply"`
}

// Apply sysctl conf to system
func (s *Sysctl) apply() error {
	b, err := share.ParseBool(s.Apply)
	if err != nil {
		return fmt.Errorf("Failed to apply: %s", s.Key)
	}

	if b != true {
		return nil
	}

	path, err := exec.LookPath("sysctl")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "-p")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to load sysctl variable: %s", stdout)
		return fmt.Errorf("Failed to load sysctl variable: %s", stdout)
	}

	return nil
}

// Read sysctl config to a map
func readConfig() (map[string]string, error) {
	lines, err := share.ReadFullFile(sysctlPath)
	if err != nil {
		return nil, err
	}

	sysctl := make(map[string]string)
	for _, line := range lines {
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			log.Errorf("could not parse line %s", line)
			continue
		}

		k := strings.TrimSpace(tokens[0])
		v := strings.TrimSpace(tokens[1])
		sysctl[k] = v
	}

	return sysctl, nil
}

//WriteConfig write config to file
func writeConfig(sysctl map[string]string) error {
	var lines []string
	var line string

	for k, v := range sysctl {
		line = k + "=" + v
		lines = append(lines, line)
	}

	return share.WriteFullFile(sysctlPath, lines)
}

// Get read sysctl file
func Get(rw http.ResponseWriter) error {
	sysctl, err := readConfig()
	if err != nil {
		return err
	}

	return share.JSONResponse(sysctl, rw)
}

// Update update sysctl file
func (s *Sysctl) Update() error {
	sysctl, err := readConfig()
	if err != nil {
		return err
	}

	sysctl[s.Key] = s.Value

	err = writeConfig(sysctl)
	if err != nil {
		return err
	}

	return s.apply()
}

// Delete delete sysctl value in file
func (s *Sysctl) Delete() error {
	sysctl, err := readConfig()
	if err != nil {
		return err
	}

	_, ok := sysctl[s.Key]
	if !ok {
		return fmt.Errorf("Failed to delete sysctl parameter '%s'. Key not found", s.Key)
	}

	delete(sysctl, s.Key)
	err = writeConfig(sysctl)
	if err != nil {
		return err
	}

	return s.apply()
}
