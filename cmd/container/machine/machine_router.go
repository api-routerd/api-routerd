// SPDX-License-Identifier: Apache-2.0

package machine

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func routerMachineGet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["command"]
	property := vars["property"]

	switch r.Method {
	case "GET":

		m := Machine{
			Path:     path,
			Property: property,
		}

		err := m.MethodGet(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func routerMachineConfigure(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["command"]
	property := vars["property"]

	switch r.Method {
	case "POST":

		m := new(Machine)
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		m.Path = path
		m.Property = property

		err = m.MethodConfigure(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

//RegisterRouterMachine register with mux
func RegisterRouterMachine(n *mux.Router) {
	m := n.PathPrefix("/machine").Subrouter().StrictSlash(false)

	m.HandleFunc("/list/{command}", routerMachineGet)
	m.HandleFunc("/get/{command}/{property}", routerMachineGet)
	m.HandleFunc("/configure/{command}/{property}", routerMachineConfigure)
}
