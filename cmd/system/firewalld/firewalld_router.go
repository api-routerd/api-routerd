// SPDX-License-Identifier: Apache-2.0

package firewalld

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func routerGetFirewalld(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	property := vars["property"]
	value := vars["value"]

	switch r.Method {
	case "GET":

		firewall := Firewall{
			Property: property,
			Value:    value,
		}

		err := firewall.GetFirewalld(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func routerConfigureFirewalld(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	property := vars["property"]

	firewall := new(Firewall)
	err := json.NewDecoder(r.Body).Decode(&firewall)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	firewall.Property = property

	switch r.Method {
	case "POST":

		err = firewall.AddFirewalld(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		break
	case "DELETE":

		err = firewall.DeleteFirewalld(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// RegisterRouterFirewalld register with mux
func RegisterRouterFirewalld(router *mux.Router) {
	f := router.PathPrefix("/firewalld").Subrouter().StrictSlash(false)

	f.HandleFunc("/get/{property}", routerGetFirewalld)
	f.HandleFunc("/get/{property}/{value}", routerGetFirewalld)
	f.HandleFunc("/set/{property}", routerConfigureFirewalld)
	f.HandleFunc("/delete/{property}", routerConfigureFirewalld)
}
