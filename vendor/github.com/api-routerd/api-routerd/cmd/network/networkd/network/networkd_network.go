// SPDX-License-Identifier: Apache-2.0

package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

const (
	networkdUnitPath = "/etc/systemd/network"
)

// Address JSON message
type Address struct {
	Address string `json:"Address"`
	Peer    string `json:"Peer"`
	Label   string `json:"Label"`
	Scope   string `json:"Scope"`
}

// Route JSON message
type Route struct {
	Gateway         string `json:"Gateway"`
	GatewayOnlink   string `json:"GatewayOnlink"`
	Destination     string `json:"Destination"`
	Source          string `json:"Source"`
	PreferredSource string `json:"PreferredSource"`
	Table           string `json:"Table"`
}

// RoutingPolicyRule JSON message
type RoutingPolicyRule struct {
	TypeOfService     string `json:"TypeOfService"`
	From              string `json:"From"`
	To                string `json:"To"`
	FirewallMark      string `json:"FirewallMark"`
	Table             string `json:"Table"`
	Priority          string `json:"Priority"`
	IncomingInterface string `json:"IncomingInterface"`
	OutgoingInterface string `json:"OutgoingInterface"`
	SourcePort        string `json:"SourcePort"`
	DestinationPort   string `json:"DestinationPort"`
	IPProtocol        string `json:"IPProtocol"`
	InvertRule        string `json:"InvertRule"`
}

// DHCPSection JSON message
type DHCPSection struct {
	UseDNS             string `json:"UseDNS"`
	UseNTP             string `json:"UseNTP"`
	UseMTU             string `json:"UseMTU"`
	Anonymize          string `json:"Anonymize"`
	SendHostname       string `json:"SendHostname"`
	UseHostname        string `json:"UseHostname"`
	Hostname           string `json:"Hostname"`
	UseDomains         string `json:"UseDomains"`
	UseRoutes          string `json:"UseRoutes"`
	UseTimezone        string `json:"UseTimezone"`
	CriticalConnection string `json:"CriticalConnection"`
	ClientIdentifier   string `json:"ClientIdentifier"`
}

// Network JSON message
type Network struct {
	ConfFile string `json:"ConfFile"`

	Match             interface{} `json:"Match"`
	Addresses         interface{} `json:"Addresses"`
	Routes            interface{} `json:"Routes"`
	RoutingPolicyRule interface{} `json:"RoutingPolicyRule"`
	DHCPSection       interface{} `json:"DHCPSection"`

	Gateway             string `json:"Gateway"`
	DHCP                string `json:"DHCP"`
	DNS                 string `json:"DNS"`
	Domains             string `json:"Domains"`
	NTP                 string `json:"NTP"`
	IPv6AcceptRA        string `json:"IPv6AcceptRA"`
	LinkLocalAddressing string `json:"LinkLocalAddressing"`
	LLDP                string `json:"LLDP"`
	EmitLLDP            string `json:"EmitLLDP"`

	Bridge  string `json:"Bridge"`
	Bond    string `json:"Bond"`
	VRF     string `json:"VRF"`
	VLAN    string `json:"VLAN"`
	MACVLAN string `json:"MACVLAN"`
	VXLAN   string `json:"VXLAN"`
	Tunnel  string `json:"Tunnel"`
}

func (network *Network) createMatchSectionConfig() string {
	conf := "[Match]\n"

	switch v := network.Match.(type) {
	case []interface{}:
		for _, b := range v {
			var mac string
			var driver string
			var name string

			if b.(map[string]interface{})["MAC"] != nil {
				mac = strings.TrimSpace(b.(map[string]interface{})["MAC"].(string))
			}

			if b.(map[string]interface{})["Driver"] != nil {
				driver = strings.TrimSpace(b.(map[string]interface{})["Driver"].(string))
			}

			if b.(map[string]interface{})["Name"] != nil {
				name = strings.TrimSpace(b.(map[string]interface{})["Name"].(string))
			}

			if mac != "" {
				mac := fmt.Sprintf("MACAddress=%s\n", mac)
				conf += mac
			}

			if driver != "" {
				driver := fmt.Sprintf("Driver=%s\n", driver)
				conf += driver
			}

			if name != "" {
				if network.ConfFile == "" {
					network.ConfFile = name
				}

				name := fmt.Sprintf("Name=%s\n", name)
				conf += name
			}
		}
		break
	}

	fmt.Println(conf)
	return conf
}

func (network *Network) createRouteSectionConfig() string {
	var routeConf string

	switch v := network.Routes.(type) {
	case []interface{}:
		for _, b := range v {
			var preferredSource string
			var gatewayOnLink string
			var destination string
			var gateway string
			var source string
			var table string

			if b.(map[string]interface{})["Gateway"] != nil {
				gateway = strings.TrimSpace(b.(map[string]interface{})["Gateway"].(string))
			}

			if b.(map[string]interface{})["GatewayOnlink"] != nil {
				gatewayOnLink = strings.TrimSpace(b.(map[string]interface{})["GatewayOnlink"].(string))
			}

			if b.(map[string]interface{})["Destination"] != nil {
				destination = strings.TrimSpace(b.(map[string]interface{})["Destination"].(string))
			}

			if b.(map[string]interface{})["Source"] != nil {
				source = strings.TrimSpace(b.(map[string]interface{})["Source"].(string))
			}

			if b.(map[string]interface{})["PreferredSource"] != nil {
				preferredSource = strings.TrimSpace(b.(map[string]interface{})["PreferredSource"].(string))
			}

			if b.(map[string]interface{})["Table"] != nil {
				table = strings.TrimSpace(b.(map[string]interface{})["Table"].(string))
			}

			routeConf += "\n[Route]\n"

			if len(gateway) != 0 {
				ip := net.ParseIP(gateway)
				if ip != nil {
					routeConf += "Gateway=" + gateway + "\n"
				} else {
					log.Error("Failed to parse Gateway: ", gateway)
				}
			}

			if len(gatewayOnLink) != 0 {
				onlink := strings.TrimSpace(gatewayOnLink)
				b, r := share.ParseBool(onlink)
				if r != nil {
					log.Error("Failed to parse GatewayOnlink: ", r, gatewayOnLink)
				} else {
					if b == true {
						routeConf += "GatewayOnlink=yes\n"
					} else {
						routeConf += "GatewayOnlink=no\n"
					}
				}
			}

			if len(destination) != 0 {
				ip := net.ParseIP(destination)
				if ip != nil {
					routeConf += "Destination=" + destination + "\n"
				} else {
					log.Error("Failed to parse Destination: ", destination)
				}
			}

			if len(source) != 0 {
				ip := net.ParseIP(source)
				if ip != nil {
					routeConf += "Source=" + source + "\n"
				} else {
					log.Error("Failed to parse Source: ", source)
				}
			}

			if len(preferredSource) != 0 {
				ip := net.ParseIP(preferredSource)
				if ip != nil {
					routeConf += "PreferredSource=" + preferredSource + "\n"
				} else {
					log.Error("Failed to parse PreferredSource: ", preferredSource)
				}
			}

			if len(table) != 0 {
				routeConf += "Table=" + table + "\n"
			}
		}
		break
	}

	return routeConf
}

func (network *Network) createAddressSectionConfig() string {
	var addressConf string

	switch v := network.Addresses.(type) {
	case []interface{}:
		for _, b := range v {
			var address string
			var peer string
			var scope string
			var label string

			if b.(map[string]interface{})["Address"] != nil {
				address = strings.TrimSpace(b.(map[string]interface{})["Address"].(string))
			}

			if b.(map[string]interface{})["Peer"] != nil {
				peer = strings.TrimSpace(b.(map[string]interface{})["Peer"].(string))
			}

			if b.(map[string]interface{})["Scope"] != nil {
				scope = strings.TrimSpace(b.(map[string]interface{})["Scope"].(string))
			}

			if b.(map[string]interface{})["Label"] != nil {
				label = strings.TrimSpace(b.(map[string]interface{})["Label"].(string))
			}

			if len(address) != 0 {
				addressConf += "\n[Address]\nAddress="

				ip := net.ParseIP(address)
				if ip != nil {
					addressConf += address + "\n"
				} else {
					log.Error("Failed to parse address: ", address)
				}

				if len(peer) != 0 {
					ip = net.ParseIP(peer)
					if ip != nil {
						addressConf += "Peer=" + peer + "\n"
					} else {
						log.Error("Failed to parse peer address: ", peer)
					}
				}

				if len(scope) != 0 {
					addressConf += "Scope=" + scope + "\n"
				}

				if len(label) != 0 {
					addressConf += "Label=" + label + "\n"
				}
			}
		}
		break
	}

	return addressConf
}

func (network *Network) createRoutingPolicyRuleSectionConfig() string {
	ruleConf := "\n[RoutingPolicyRule]\n"

	switch v := network.RoutingPolicyRule.(type) {
	case []interface{}:
		for _, b := range v {
			var IncomingInterface string
			var OutgoingInterface string
			var DestinationPort string
			var TypeOfService string
			var FirewallMark string
			var SourcePort string
			var IPProtocol string
			var InvertRule string
			var Priority string
			var Table string
			var From string
			var To string

			if b.(map[string]interface{})["TypeOfService"] != nil {
				TypeOfService = strings.TrimSpace(b.(map[string]interface{})["TypeOfService"].(string))

				if TypeOfService != "" {
					ruleConf += "TypeOfService=" + TypeOfService + "\n"
				}
			}

			if b.(map[string]interface{})["From"] != nil {
				From = strings.TrimSpace(b.(map[string]interface{})["From"].(string))

				if From != "" {
					ruleConf += "From=" + From + "\n"
				}
			}

			if b.(map[string]interface{})["To"] != nil {
				To = strings.TrimSpace(b.(map[string]interface{})["To"].(string))

				if To != "" {
					ruleConf += "To=" + To + "\n"
				}
			}

			if b.(map[string]interface{})["FirewallMark"] != nil {
				FirewallMark = strings.TrimSpace(b.(map[string]interface{})["FirewallMark"].(string))

				if FirewallMark != "" {
					ruleConf += "FirewallMark=" + FirewallMark + "\n"
				}
			}

			if b.(map[string]interface{})["Table"] != nil {
				Table = strings.TrimSpace(b.(map[string]interface{})["Table"].(string))

				if Table != "" {
					ruleConf += "Table=" + Table + "\n"
				}
			}

			if b.(map[string]interface{})["Priority"] != nil {
				Priority = strings.TrimSpace(b.(map[string]interface{})["Priority"].(string))

				if Priority != "" {
					ruleConf += "Priority=" + Priority + "\n"
				}
			}

			if b.(map[string]interface{})["IncomingInterface"] != nil {
				IncomingInterface = strings.TrimSpace(b.(map[string]interface{})["IncomingInterface"].(string))

				if IncomingInterface != "" {
					ruleConf += "IncomingInterface=" + IncomingInterface + "\n"
				}
			}

			if b.(map[string]interface{})["OutgoingInterface"] != nil {
				OutgoingInterface = strings.TrimSpace(b.(map[string]interface{})["OutgoingInterface"].(string))

				if OutgoingInterface != "" {
					ruleConf += "OutgoingInterface=" + OutgoingInterface + "\n"
				}

			}

			if b.(map[string]interface{})["SourcePort"] != nil {
				SourcePort = strings.TrimSpace(b.(map[string]interface{})["SourcePort"].(string))

				if SourcePort != "" {
					ruleConf += "SourcePort=" + SourcePort + "\n"
				}
			}

			if b.(map[string]interface{})["DestinationPort"] != nil {
				DestinationPort = strings.TrimSpace(b.(map[string]interface{})["DestinationPort"].(string))

				if DestinationPort != "" {
					ruleConf += "DestinationPort=" + DestinationPort + "\n"
				}
			}

			if b.(map[string]interface{})["IPProtocol"] != nil {
				IPProtocol = strings.TrimSpace(b.(map[string]interface{})["IPProtocol"].(string))

				if IPProtocol != "" {
					ruleConf += "IPProtocol=" + IPProtocol + "\n"
				}
			}

			if b.(map[string]interface{})["InvertRule"] != nil {
				InvertRule = strings.TrimSpace(b.(map[string]interface{})["InvertRule"].(string))

				if InvertRule != "" {
					ruleConf += "InvertRule=" + InvertRule + "\n"
				}
			}
		}

		break
	}
	return ruleConf
}

func (network *Network) createDHCPSectionConfig() string {
	dhcpConf := "\n[DHCP]\n"

	switch v := network.DHCPSection.(type) {
	case []interface{}:
		for _, b := range v {
			var UseDNS string
			var UseNTP string
			var UseMTU string
			var Anonymize string
			var SendHostname string
			var UseHostname string
			var Hostname string
			var UseDomains string
			var UseRoutes string
			var UseTimezone string
			var CriticalConnection string
			var ClientIdentifier string

			if b.(map[string]interface{})["UseDNS"] != nil {
				UseDNS = strings.TrimSpace(b.(map[string]interface{})["UseDNS"].(string))

				if UseDNS != "" {
					dhcpConf += "UseDNS=" + UseDNS + "\n"
				}
			}

			if b.(map[string]interface{})["UseNTP"] != nil {
				UseNTP = strings.TrimSpace(b.(map[string]interface{})["UseNTP"].(string))

				if UseNTP != "" {
					dhcpConf += "UseNTP=" + UseNTP + "\n"
				}
			}

			if b.(map[string]interface{})["UseMTU"] != nil {
				UseMTU = strings.TrimSpace(b.(map[string]interface{})["UseMTU"].(string))

				if UseMTU != "" {
					dhcpConf += "UseMTU=" + UseMTU + "\n"
				}
			}

			if b.(map[string]interface{})["Anonymize"] != nil {
				Anonymize = strings.TrimSpace(b.(map[string]interface{})["Anonymize"].(string))

				if Anonymize != "" {
					dhcpConf += "Anonymize=" + Anonymize + "\n"
				}
			}

			if b.(map[string]interface{})["SendHostname"] != nil {
				SendHostname = strings.TrimSpace(b.(map[string]interface{})["SendHostname"].(string))

				if SendHostname != "" {
					dhcpConf += "SendHostname=" + SendHostname + "\n"
				}
			}

			if b.(map[string]interface{})["Hostname"] != nil {
				Hostname = strings.TrimSpace(b.(map[string]interface{})["Hostname"].(string))

				if Hostname != "" {
					dhcpConf += "Hostname=" + Hostname + "\n"
				}
			}

			if b.(map[string]interface{})["UseHostname"] != nil {
				UseHostname = strings.TrimSpace(b.(map[string]interface{})["UseHostname"].(string))

				if UseHostname != "" {
					dhcpConf += "UseHostname=" + UseHostname + "\n"
				}
			}

			if b.(map[string]interface{})["UseDomains"] != nil {
				UseDomains = strings.TrimSpace(b.(map[string]interface{})["UseDomains"].(string))

				if UseDomains != "" {
					dhcpConf += "UseDomains=" + UseDomains + "\n"
				}
			}

			if b.(map[string]interface{})["UseRoutes"] != nil {
				UseRoutes = strings.TrimSpace(b.(map[string]interface{})["UseRoutes"].(string))

				if UseRoutes != "" {
					dhcpConf += "UseRoutes=" + UseRoutes + "\n"
				}

			}

			if b.(map[string]interface{})["UseTimezone"] != nil {
				UseTimezone = strings.TrimSpace(b.(map[string]interface{})["UseTimezone"].(string))

				if UseTimezone != "" {
					dhcpConf += "UseTimezone=" + UseTimezone + "\n"
				}
			}

			if b.(map[string]interface{})["CriticalConnection"] != nil {
				CriticalConnection = strings.TrimSpace(b.(map[string]interface{})["CriticalConnection"].(string))

				if CriticalConnection != "" {
					dhcpConf += "CriticalConnection=" + CriticalConnection + "\n"
				}
			}

			if b.(map[string]interface{})["ClientIdentifier"] != nil {
				ClientIdentifier = strings.TrimSpace(b.(map[string]interface{})["ClientIdentifier"].(string))

				if ClientIdentifier != "" {
					dhcpConf += "ClientIdentifier=" + ClientIdentifier + "\n"
				}
			}
		}

		break
	}
	return dhcpConf
}

// createNetworkSectionConfig generate config
func (network *Network) createNetworkSectionConfig() string {
	conf := "[Network]\n"

	if network.DHCP != "" {
		dhcpConf := "DHCP="

		dhcp := strings.TrimSpace(network.DHCP)
		b, r := share.ParseBool(dhcp)
		if r != nil {
			switch dhcp {
			case "ipv4", "ipv6":
				dhcpConf += dhcp
				break
			default:
				log.Error("Failed to parse DHCP: ", r, network.DHCP)
			}
		} else {
			if b == true {
				dhcpConf += "yes"
			} else {
				dhcpConf += "no"
			}
		}
		conf += dhcpConf + "\n"
	}

	if network.Gateway != "" {
		gatewayConf := "Gateway="

		gw := strings.TrimSpace(network.Gateway)

		ip := net.ParseIP(gw)
		if ip != nil {
			gatewayConf += gw
			conf += gatewayConf + "\n"
		} else {
			log.Error("Failed to parse Gateway Address: ", network.Gateway)
		}
	}

	if network.DNS != "" {
		conf += "DNS=" + network.DNS
	}

	if network.Domains != "" {
		conf += "Domains=" + network.Domains + "\n"
	}

	if network.NTP != "" {
		conf += "NTP=" + network.NTP + "\n"
	}

	if network.IPv6AcceptRA != "" {
		IPv6AcceptRAConf := "IPv6AcceptRA="

		IPv6RA := strings.TrimSpace(network.IPv6AcceptRA)
		b, err := share.ParseBool(IPv6RA)
		if err != nil {
			log.Error("Failed to parse IPv6AcceptRA: ", err, network.IPv6AcceptRA)
		} else {
			if b == true {
				IPv6AcceptRAConf += "yes"
			} else {
				IPv6AcceptRAConf += "no"
			}
		}
		conf += IPv6AcceptRAConf + "\n"
	}

	if network.LinkLocalAddressing != "" {
		LinkLocalAddressingConf := "LinkLocalAddressing="

		IPv6RA := strings.TrimSpace(network.LinkLocalAddressing)
		b, err := share.ParseBool(IPv6RA)
		if err != nil {
			log.Error("Failed to parse LinkLocalAddressing: ", err, network.LinkLocalAddressing)
		} else {
			if b == true {
				LinkLocalAddressingConf += "yes"
			} else {
				LinkLocalAddressingConf += "no"
			}
		}
		conf += LinkLocalAddressingConf + "\n"
	}

	if network.LLDP != "" {
		LLDPConf := "LLDP="

		LLDP := strings.TrimSpace(network.LLDP)
		b, err := share.ParseBool(LLDP)
		if err != nil {
			log.Error("Failed to parse LLDP: ", err, network.LLDP)
		} else {
			if b == true {
				LLDPConf += "yes"
			} else {
				LLDPConf += "no"
			}
		}
		conf += LLDPConf + "\n"
	}

	if network.EmitLLDP != "" {
		EmitLLDPConf := "EmitLLDP="

		EmitLLDP := strings.TrimSpace(network.EmitLLDP)
		b, err := share.ParseBool(EmitLLDP)
		if err != nil {
			log.Error("Failed to parse EmitLLDP: ", err, network.EmitLLDP)
		} else {
			if b == true {
				EmitLLDPConf += "yes"
			} else {
				EmitLLDPConf += "no"
			}
		}
		conf += EmitLLDPConf + "\n"
	}

	if network.NTP != "" {
		conf += "NTP=" + network.NTP + "\n"
	}

	if network.Bridge != "" {
		conf += "Bridge=" + network.Bridge + "\n"
	}

	if network.Bond != "" {
		conf += "Bond=" + network.Bond + "\n"
	}

	if network.VLAN != "" {
		conf += "VLAN=" + network.VLAN + "\n"
	}

	return conf
}

func parseJSONfromHTTPReq(req *http.Request) error {
	var configs map[string]interface{}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("Failed to parse HTTP request: %v", err)
		return err
	}

	json.Unmarshal([]byte(body), &configs)

	network := new(Network)
	json.Unmarshal([]byte(body), &network)

	matchConfig := network.createMatchSectionConfig()
	networkConfig := network.createNetworkSectionConfig()
	addressConfig := network.createAddressSectionConfig()
	routeConfig := network.createRouteSectionConfig()
	ruleConfig := network.createRoutingPolicyRuleSectionConfig()
	dhcpConfig := network.createDHCPSectionConfig()

	config := []string{matchConfig, networkConfig, addressConfig, routeConfig, ruleConfig, dhcpConfig}

	unitName := fmt.Sprintf("25-%s.network", network.ConfFile)
	unitPath := filepath.Join(networkdUnitPath, unitName)

	return share.WriteFullFile(unitPath, config)
}

// CreateFile generate .network
func CreateFile(rw http.ResponseWriter, req *http.Request) {
	parseJSONfromHTTPReq(req)
}
