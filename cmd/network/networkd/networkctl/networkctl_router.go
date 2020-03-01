// SPDX-License-Identifier: Apache-2.0

package networkctl

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routerNetworkctlGet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	verb := vars["verb"]
	link := vars["link"]

	n := Networkctl{
		Link: link,
		Verb: verb,
	}

	switch r.Method {
	case "GET":

		err := n.NetworkctlGet(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// RegisterRouterNetworkctl register with mux
func RegisterRouterNetworkctl(r *mux.Router) {

	n := r.PathPrefix("/networkctl").Subrouter().StrictSlash(false)
	n.HandleFunc("", routerNetworkctlGet)
	n.HandleFunc("/{verb}", routerNetworkctlGet)
	n.HandleFunc("/{verb}/{link}", routerNetworkctlGet)
}
