// SPDX-License-Identifier: Apache-2.0

package route

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"syscall"

	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

//Route message
type Route struct {
	Action  string `json:"action"`
	Link    string `json:"link"`
	Gateway string `json:"gateway"`
	OnLink  string `json:"onlink"`
}

//DecodeJSONRequest decoded route JSON message
func DecodeJSONRequest(r *http.Request) (Route, error) {
	route := new(Route)

	err := json.NewDecoder(r.Body).Decode(&route)
	if err != nil {
		return *route, nil
	}

	return *route, nil
}

//AddDefaultGateWay add a default GW
func (route *Route) AddDefaultGateWay() error {
	link, err := netlink.LinkByName(route.Link)
	if err != nil {
		log.Errorf("Failed to find link %s: %v", err, route.Link)
		return err
	}

	ipAddr, _, err := net.ParseCIDR(route.Gateway)
	if err != nil {
		log.Errorf("Failed to parse default GateWay address %s: %v", route.Gateway, err)
		return err
	}

	onlink := 0
	b, err := share.ParseBool(strings.TrimSpace(route.OnLink))
	if err != nil {
		log.Errorf("Failed to parse GatewayOnlink %s: %v", route.OnLink, err)
	} else {
		if b == true {
			onlink |= syscall.RTNH_F_ONLINK
		}
	}

	// add a gateway route
	rt := &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: link.Attrs().Index,
		Gw:        ipAddr,
		Flags:     onlink,
	}

	err = netlink.RouteAdd(rt)
	if err != nil {
		log.Errorf("Failed to add default GateWay address %s: %v", route.Gateway, err)
		return err
	}

	return nil
}

//ReplaceDefaultGateWay replace default GW with new one
func (route *Route) ReplaceDefaultGateWay() error {
	link, err := netlink.LinkByName(route.Link)
	if err != nil {
		return err
	}

	ipAddr, _, err := net.ParseCIDR(route.Gateway)
	if err != nil {
		log.Errorf("Failed to parse default GateWay address %s: %v", route.Gateway, err)
		return err
	}

	onlink := 0
	b, err := share.ParseBool(strings.TrimSpace(route.OnLink))
	if err != nil {
		log.Errorf("Failed to parse GatewayOnlink %s: %v", route.OnLink, err)
	} else {
		if b == true {
			onlink |= syscall.RTNH_F_ONLINK
		}
	}

	// add a gateway route
	rt := &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: link.Attrs().Index,
		Gw:        ipAddr,
		Flags:     onlink,
	}

	err = netlink.RouteReplace(rt)
	if err != nil {
		log.Errorf("Failed to replace default GateWay address %s: %v", route.Gateway, err)
		return err
	}

	return nil
}

//DeleteGateWay remove a gateway
func (route *Route) DeleteGateWay() error {
	link, err := netlink.LinkByName(route.Link)
	if err != nil {
		log.Errorf("Failed to delete default gateway %s: %v", link, err)
		return err
	}

	ipAddr, _, err := net.ParseCIDR(route.Gateway)
	if err != nil {
		return err
	}

	switch route.Action {
	case "del-default-gw":
		// del a gateway route
		rt := &netlink.Route{
			Scope:     netlink.SCOPE_UNIVERSE,
			LinkIndex: link.Attrs().Index,
			Gw:        ipAddr,
		}

		err = netlink.RouteDel(rt)
		if err != nil {
			log.Errorf("Failed to delete default GateWay address %s: %v", ipAddr, err)
			return err
		}
		break
	}

	return nil
}

//Get get routes
func (route *Route) Get(rw http.ResponseWriter) error {
	routes, err := netlink.RouteList(nil, 0)
	if err != nil {
		return err
	}

	return share.JSONResponse(routes, rw)
}

//Configure routes
func (route *Route) Configure() error {
	switch route.Action {
	case "add-default-gw":
		return route.AddDefaultGateWay()
	case "replace-default-gw":
		return route.ReplaceDefaultGateWay()
	}

	return nil
}
