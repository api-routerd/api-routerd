// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"net/http"

	"github.com/api-routerd/api-routerd/cmd/network/networkd/link"
	"github.com/api-routerd/api-routerd/cmd/network/networkd/netdev"
	"github.com/api-routerd/api-routerd/cmd/network/networkd/network"
	"github.com/api-routerd/api-routerd/cmd/network/networkd/networkctl"

	"github.com/gorilla/mux"
)

func routerConfigureNetworkdLink(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		link.CreateFile(rw, r)
	}
}

func routerConfigureNetworkdNetDev(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		netdev.CreateFile(rw, r)
	}
}

func routerConfigureNetworkdNetwork(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		network.CreateFile(rw, r)
	}
}

// RegisterRouterNetworkd register with mux
func RegisterRouterNetworkd(router *mux.Router) {
	InitNetworkd()

	n := router.PathPrefix("/networkd").Subrouter().StrictSlash(false)

	// systemd-networkd
	n.HandleFunc("/network", routerConfigureNetworkdNetwork)
	n.HandleFunc("/netdev", routerConfigureNetworkdNetDev)
	n.HandleFunc("/link", routerConfigureNetworkdLink)

	// networkctl
	networkctl.RegisterRouterNetworkctl(n)
}
