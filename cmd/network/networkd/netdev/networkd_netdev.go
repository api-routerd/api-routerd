// SPDX-License-Identifier: Apache-2.0

package netdev

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

const (
	networkdUnitPath = "/etc/systemd/network"
)

// NetDev JSON message
type NetDev struct {
	Description string `json:"Description"`
	MACAddress  string `json:"MACAddress"`
	MTUBytes    string `json:"MTUBytes"`
	Name        string `json:"Name"`
	Kind        string `json:"Kind"`

	// Bond
	Mode               string `json:"Mode"`
	TransmitHashPolicy string `json:"TransmitHashPolicy"`

	// Vlan
	VlanID string `json:"VlanId"`

	//Bridge
	HelloTimeSec    string `json:"HelloTimeSec"`
	ForwardDelaySec string `json:"ForwardDelaySec"`
	AgeingTimeSec   string `json:"AgeingTimeSec"`

	//Tunnel
	Local              string `json:"Local"`
	Remote             string `json:"Remote"`
	TTL                string `json:"TTL"`
	DiscoverPathMTU    string `json:"DiscoverPathMTU"`
	IPv6FlowLabel      string `json:"IPv6FlowLabel"`
	EncapsulationLimit string `json:"EncapsulationLimit"`
	Key                string `json:"Key"`
	Independent        string `json:"Independent"`

	//VxLan
	ID              string `json:"Id"`
	TOS             string `json:"TOS"`
	MacLearning     string `json:"MacLearning"`
	DestinationPort string `json:"DestinationPort"`
	PortRange       string `json:"PortRange"`
	FlowLabel       string `json:"FlowLabel"`

	//Veth
	Peer           string `json:"Peer"`
	PeerMACAddress string `json:"PeerMACAddress"`
}

func (netdev *NetDev) createBondSectionConfig() string {
	conf := "\n[Bond]\n"

	if netdev.Mode != "" {
		conf += "Mode=" + strings.TrimSpace(netdev.Mode) + "\n"
	}

	if netdev.TransmitHashPolicy != "" {
		conf += "TransmitHashPolicy=" + strings.TrimSpace(netdev.TransmitHashPolicy) + "\n"
	}

	return conf
}

func (netdev *NetDev) createBridgeSectionConfig() string {
	conf := "\n[Bridge]\n"

	if netdev.HelloTimeSec != "" {
		conf += "HelloTimeSec=" + strings.TrimSpace(netdev.HelloTimeSec) + "\n"
	}

	if netdev.ForwardDelaySec != "" {
		conf += "ForwardDelaySec=" + strings.TrimSpace(netdev.ForwardDelaySec) + "\n"
	}

	if netdev.AgeingTimeSec != "" {
		conf += "AgeingTimeSec=" + strings.TrimSpace(netdev.AgeingTimeSec) + "\n"
	}

	return conf
}

func (netdev *NetDev) createTunnelSectionConfig() string {
	conf := "\n[Tunnel]\n"

	if netdev.Local != "" {
		conf += "Local=" + strings.TrimSpace(netdev.Local) + "\n"
	}

	if netdev.Remote != "" {
		conf += "Remote=" + strings.TrimSpace(netdev.Remote) + "\n"
	}

	if netdev.TTL != "" {
		conf += "TTL=" + strings.TrimSpace(netdev.TTL) + "\n"
	}

	if netdev.DiscoverPathMTU != "" {
		conf += "DiscoverPathMTU=" + strings.TrimSpace(netdev.DiscoverPathMTU) + "\n"
	}

	if netdev.IPv6FlowLabel != "" {
		conf += "IPv6FlowLabel=" + strings.TrimSpace(netdev.IPv6FlowLabel) + "\n"
	}

	if netdev.EncapsulationLimit != "" {
		conf += "EncapsulationLimit=" + strings.TrimSpace(netdev.EncapsulationLimit) + "\n"
	}

	if netdev.Key != "" {
		conf += "Key=" + strings.TrimSpace(netdev.Key) + "\n"
	}

	if netdev.Independent != "" {
		conf += "Independent=" + strings.TrimSpace(netdev.Independent) + "\n"
	}

	return conf
}

func (netdev *NetDev) createVxLanSectionConfig() string {
	conf := "\n[VXLAN]\n"

	if netdev.ID != "" {
		conf += "Id=" + strings.TrimSpace(netdev.ID) + "\n"
	}

	if netdev.Local != "" {
		conf += "Local=" + strings.TrimSpace(netdev.Local) + "\n"
	}

	if netdev.Remote != "" {
		conf += "Remote=" + strings.TrimSpace(netdev.Remote) + "\n"
	}

	if netdev.TOS != "" {
		conf += "TOS=" + strings.TrimSpace(netdev.TOS) + "\n"
	}

	if netdev.TTL != "" {
		conf += "TTL=" + strings.TrimSpace(netdev.TTL) + "\n"
	}

	if netdev.MacLearning != "" {
		conf += "MacLearning=" + strings.TrimSpace(netdev.MacLearning) + "\n"
	}

	if netdev.DestinationPort != "" {
		conf += "DestinationPort=" + strings.TrimSpace(netdev.DestinationPort) + "\n"
	}

	if netdev.PortRange != "" {
		conf += "PortRange=" + strings.TrimSpace(netdev.PortRange) + "\n"
	}

	if netdev.FlowLabel != "" {
		conf += "FlowLabel=" + strings.TrimSpace(netdev.FlowLabel) + "\n"
	}

	return conf
}

// CreateNetDevSectionConfig generate netdev config
func (netdev *NetDev) CreateNetDevSectionConfig() string {
	conf := "[NetDev]\n"

	if netdev.Name != "" {
		conf += "Name=" + strings.TrimSpace(netdev.Name) + "\n"
	}

	if netdev.Description != "" {
		conf += "Description=" + strings.TrimSpace(netdev.Description) + "\n"
	}

	if netdev.Kind != "" {
		conf += "Kind=" + strings.TrimSpace(netdev.Kind) + "\n"
	}

	if netdev.MACAddress != "" {
		conf += "MACAddress=" + strings.TrimSpace(netdev.MACAddress) + "\n"
	}

	if netdev.MTUBytes != "" {
		conf += "MTUBytes=" + strings.TrimSpace(netdev.MTUBytes) + "\n"
	}

	switch netdev.Kind {
	case "bond":

		conf += netdev.createBondSectionConfig()

		break

	case "vlan":

		conf += "\n[VLAN]\n"

		if netdev.VlanID != "" {
			conf += "Id=" + strings.TrimSpace(netdev.VlanID) + "\n"
		}

		break

	case "bridge":

		conf += netdev.createBridgeSectionConfig()

		break
	case "tunnel":

		conf += netdev.createTunnelSectionConfig()

		break
	case "veth":

		conf += "\n[Peer]\n"

		if netdev.Peer != "" {
			conf += "Name=" + strings.TrimSpace(netdev.Peer) + "\n"
		}

		if netdev.PeerMACAddress != "" {
			conf += "MACAddress=" + strings.TrimSpace(netdev.PeerMACAddress) + "\n"
		}

		break

	case "macvlan":

		conf += "\n[MACVLAN]\n"

		if netdev.Peer != "" {
			conf += "Mode=" + strings.TrimSpace(netdev.Mode) + "\n"
		}

		break
	case "macvtap":

		conf += "\n[MACVTAP\n"

		if netdev.Peer != "" {
			conf += "Mode=" + strings.TrimSpace(netdev.Mode) + "\n"
		}

		break

	case "ipvlan":

		conf += "\n[IPVLAN]\n"

		if netdev.Peer != "" {
			conf += "Mode=" + strings.TrimSpace(netdev.Mode) + "\n"
		}

	case "vxlan":

		conf += netdev.createVxLanSectionConfig()

	}

	return conf
}

func parseJSONFromHTTPReq(req *http.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("Failed to parse HTTP request: %v", err)
		return err
	}

	netdev := new(NetDev)
	json.Unmarshal([]byte(body), &netdev)

	netdevConfig := netdev.CreateNetDevSectionConfig()
	config := []string{netdevConfig}

	fmt.Println(config)

	unitName := fmt.Sprintf("25-%s.netdev", netdev.Name)
	unitPath := filepath.Join(networkdUnitPath, unitName)

	return share.WriteFullFile(unitPath, config)
}

// CreateFile generate .netdev
func CreateFile(rw http.ResponseWriter, req *http.Request) {
	parseJSONFromHTTPReq(req)
}
