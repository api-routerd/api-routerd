// SPDX-License-Identifier: Apache-2.0

package link

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/api-routerd/api-routerd/cmd/share"

	"github.com/vishvananda/netlink"
)

// Link JSON message
type Link struct {
	Action  string   `json:"action"`
	Link    string   `json:"link"`
	MTU     string   `json:"mtu"`
	Kind    string   `json:"kind"`
	Mode    string   `json:"mode"`
	Enslave []string `json:"enslave"`
}

// DecodeJSONRequest decodes JSON message to type Link
func DecodeJSONRequest(r *http.Request) (Link, error) {
	link := new(Link)

	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		return *link, err
	}

	return *link, nil
}

// Set sets link status/attribute
func (link *Link) Set() error {
	switch link.Action {
	case "set-link-up":
		return setUp(link.Link)
	case "set-link-down":
		return setDown(link.Link)
	case "set-link-mtu":

		mtu, err := strconv.ParseInt(strings.TrimSpace(link.MTU), 10, 64)
		if err != nil {
			return err
		}

		return setMTU(link.Link, int(mtu))
	}

	return nil
}

// Get all link
func (link *Link) Get(rw http.ResponseWriter) error {
	if link.Link != "" {
		l, err := netlink.LinkByName(link.Link)
		if err != nil {
			return err
		}

		return share.JSONResponse(l, rw)

	}

	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

	return share.JSONResponse(links, rw)
}

//Delete remove a netdev
func (link *Link) Delete() error {
	l, err := netlink.LinkByName(link.Link)
	if err != nil {
		return err
	}

	err = netlink.LinkDel(l)
	if err != nil {
		return err
	}

	return nil
}

// Create create virtual netdevs
func (link *Link) Create() error {
	switch link.Action {
	case "add-link-bridge":
		return link.createBridge()
	case "add-link-bond":
		return link.createBond()
	}

	return nil
}
