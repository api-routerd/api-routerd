// SPDX-License-Identifier: Apache-2.0

package link

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

// Link JSON message
type Link struct {
	ConfFile string      `json:"ConfFile"`
	Match    interface{} `json:"Match"`

	Description                string `json:"Description"`
	Alias                      string `json:"Alias"`
	MACAddressPolicy           string `json:"MACAddressPolicy"`
	MACAddress                 string `json:"MACAddress"`
	NamePolicy                 string `json:"NamePolicy"`
	Name                       string `json:"Name"`
	MTUBytes                   string `json:"MTUBytes"`
	BitsPerSecond              string `json:"BitsPerSecond"`
	Duplex                     string `json:"Duplex"`
	AutoNegotiation            string `json:"AutoNegotiation"`
	WakeOnLan                  string `json:"WakeOnLan"`
	Port                       string `json:"Port"`
	TCPSegmentationOffload     string `json:"TCPSegmentationOffload"`
	TCP6SegmentationOffload    string `json:"TCP6SegmentationOffload"`
	GenericSegmentationOffload string `json:"GenericSegmentationOffload"`
	GenericReceiveOffload      string `json:"GenericReceiveOffload"`
	LargeReceiveOffload        string `json:"LargeReceiveOffload"`
	RxChannels                 string `json:"RxChannels"`
	TxChannels                 string `json:"TxChannels"`
	OtherChannels              string `json:"OtherChannels"`
	CombinedChannels           string `json:"CombinedChannels"`
}

func (link *Link) createMatchSectionConfig() string {
	conf := "[Match]\n"

	switch v := link.Match.(type) {
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
				if link.ConfFile == "" {
					link.ConfFile = name
				}

				name := fmt.Sprintf("Name=%s\n", name)
				conf += name
			}
		}
		break
	}

	return conf
}

func (link *Link) createLinkSectionConfig() string {
	conf := "[Link]\n"

	if link.Description != "" {
		conf += "Description=" + link.Description + "\n"
	}

	if link.Alias != "" {
		conf += "Alias=" + link.Alias + "\n"
	}

	if link.MACAddressPolicy != "" {
		conf += "MACAddressPolicy=" + link.MACAddressPolicy + "\n"
	}

	if link.MACAddress != "" {
		conf += "MACAddress=" + link.MACAddress + "\n"
	}

	if link.NamePolicy != "" {
		conf += "NamePolicy=" + link.NamePolicy + "\n"
	}

	if link.Name != "" {
		conf += "Name=" + link.Name + "\n"
	}

	if link.MTUBytes != "" {
		conf += "MTUBytes=" + link.MTUBytes + "\n"
	}

	if link.BitsPerSecond != "" {
		conf += "BitsPerSecond=" + link.BitsPerSecond + "\n"
	}

	if link.Duplex != "" {
		conf += "Duplex=" + link.Duplex + "\n"
	}

	if link.AutoNegotiation != "" {
		conf += "AutoNegotiation=" + link.AutoNegotiation + "\n"
	}

	if link.WakeOnLan != "" {
		conf += "WakeOnLan=" + link.WakeOnLan + "\n"
	}

	if link.Port != "" {
		conf += "Port=" + link.Port + "\n"
	}

	if link.TCPSegmentationOffload != "" {
		conf += "TCPSegmentationOffload=" + link.TCPSegmentationOffload + "\n"
	}

	if link.TCP6SegmentationOffload != "" {
		conf += "TCP6SegmentationOffload=" + link.TCP6SegmentationOffload + "\n"
	}

	if link.GenericSegmentationOffload != "" {
		conf += "GenericSegmentationOffload=" + link.GenericSegmentationOffload + "\n"
	}

	if link.GenericReceiveOffload != "" {
		conf += "GenericReceiveOffload=" + link.GenericReceiveOffload + "\n"
	}

	if link.LargeReceiveOffload != "" {
		conf += "LargeReceiveOffload=" + link.LargeReceiveOffload + "\n"
	}

	if link.RxChannels != "" {
		conf += "RxChannels=" + link.RxChannels + "\n"
	}

	if link.TxChannels != "" {
		conf += "TxChannels=" + link.TxChannels + "\n"
	}

	if link.OtherChannels != "" {
		conf += "OtherChannels=" + link.OtherChannels + "\n"
	}

	if link.CombinedChannels != "" {
		conf += "CombinedChannels=" + link.CombinedChannels + "\n"
	}

	return conf
}

func parseJSONFromHTTPReq(req *http.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("Failed to parse HTTP request: %v ", err)
		return err
	}

	link := new(Link)
	json.Unmarshal([]byte(body), &link)

	matchConfig := link.createMatchSectionConfig()
	linkConfig := link.createLinkSectionConfig()

	config := []string{matchConfig, linkConfig}

	unitName := fmt.Sprintf("00-%s.link", link.Name)
	unitPath := filepath.Join(networkdUnitPath, unitName)

	return share.WriteFullFile(unitPath, config)
}

// CreateFile generate .link file
func CreateFile(rw http.ResponseWriter, req *http.Request) {
	parseJSONFromHTTPReq(req)
}
