// SPDX-License-Identifier: Apache-2.0

package networkctl

import (
	"os/exec"
	"strings"

	"github.com/api-routerd/api-routerd/cmd/share"
)

// Response JSON response for networkctl
type Response struct {
	Index       string `json:"index"`
	Link        string `json:"link"`
	Type        string `json:"type"`
	Operational string `json:"operational"`
	Setup       string `json:"setup"`
}

// ResponseLLDP LLDP response
type ResponseLLDP struct {
	Link       string `json:"link"`
	ChassisID  string `json:"chassis_id"`
	System     string `json:"system"`
	Capability string `json:"capability"`
	Port       string `json:"port"`
	PortDdesc  string `json:"port_description"`
}

// ResponseStatus one link status response
type ResponseStatus struct {
	Link        string   `json:"link"`
	Index       string   `json:"index"`
	LinkFile    string   `json:"link_file"`
	NetworkFile string   `json:"network_file"`
	Type        string   `json:"type"`
	State       string   `json:"state"`
	Driver      string   `json:"driver"`
	HWAddress   string   `json:"hw_address"`
	Address     []string `json:"address"`
	ConnectedTo []string `json:"connected_to"`
}

func getPathNetworkctl() (string, error) {
	err := share.CheckBinaryExists("networkctl")
	if err != nil {
		return "", err
	}

	path, err := exec.LookPath("networkctl")
	if err != nil {
		return "", err
	}

	return path, nil
}

// ExecuteNetworkctl execute networkctl same as status
func ExecuteNetworkctl() ([]Response, error) {
	path, err := getPathNetworkctl()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(path)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(stdout), "\n")

	r := make([]Response, len(lines)-4)
	for i, line := range lines {
		if i == 0 || i >= len(lines)-3 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		n := Response{
			Index:       fields[0],
			Link:        fields[1],
			Type:        fields[2],
			Operational: fields[3],
			Setup:       fields[4],
		}

		r[i-1] = n
	}

	return r, nil
}

// ExecuteNetworkctlLLDP execute networkctl same as status
func ExecuteNetworkctlLLDP() ([]ResponseLLDP, error) {
	path, err := getPathNetworkctl()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(path, "lldp")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(stdout), "\n")

	r := make([]ResponseLLDP, len(lines)-6)
	for i, line := range lines {
		if i == 0 || i >= len(lines)-4 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		n := ResponseLLDP{
			Link:       fields[0],
			ChassisID:  fields[1],
			System:     fields[2],
			Capability: fields[3],
			Port:       fields[4],
			PortDdesc:  fields[5],
		}

		r[i-1] = n
	}

	return r, nil
}

// ExecuteNetworkctlStatus execute networkctl same as status
func ExecuteNetworkctlStatus(link string) (ResponseStatus, error) {
	r := ResponseStatus{}

	path, err := getPathNetworkctl()
	if err != nil {
		return r, err
	}

	cmd := exec.Command(path, "status", link)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return r, err
	}

	lines := strings.Split(string(stdout), "\n")

	address := false
	lldp := false
	for i, line := range lines {
		if len(strings.TrimSpace(line)) <= 0 {
			continue
		}

		fields := strings.SplitN(line, ":", 2)
		switch i {
		case 0:
			r.Index = strings.TrimSpace(fields[0])
			r.Link = strings.TrimSpace(fields[1])
			break
		case 1:
			r.LinkFile = strings.TrimSpace(fields[1])
			break
		case 2:
			r.NetworkFile = strings.TrimSpace(fields[1])
			break
		case 3:
			r.Type = strings.TrimSpace(fields[1])
			break
		case 4:
			r.State = strings.TrimSpace(fields[1])
			break
		case 5:
			r.Driver = strings.TrimSpace(fields[1])
			break
		case 6:
			r.HWAddress = strings.TrimSpace(fields[1])
			break
		default:

			if strings.Contains(line, "Address") {
				if len(fields) == 2 {
					r.Address = append(r.Address, strings.TrimSpace(fields[1]))
					address = true
					continue
				}
			}

			if strings.Contains(line, "Connected To") {
				if len(fields) == 2 {
					r.ConnectedTo = append(r.ConnectedTo, strings.TrimSpace(fields[1]))
					address = false
					lldp = true
					continue
				}
			}

			if address == true {
				r.Address = append(r.Address, strings.TrimSpace(line))
			} else if lldp == true {
				r.ConnectedTo = append(r.ConnectedTo, strings.TrimSpace(line))
			}
		}
	}

	return r, nil
}
