// SPDX-License-Identifier: Apache-2.0

package coredump

import (
	"net/http"

	"github.com/gorilla/mux"
)

func configureSystemdCoreDump(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		err := GetConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "POST":

		err := UpdateConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	case "DELETE":

		err := DeleteConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}
}

// RegisterRouterCoreDump register with mux
func RegisterRouterCoreDump(router *mux.Router) {
	c := router.PathPrefix("/coredump").Subrouter().StrictSlash(false)

	// coredump.conf
	c.HandleFunc("/get", configureSystemdCoreDump)
	c.HandleFunc("/add", configureSystemdCoreDump)
	c.HandleFunc("/delete", configureSystemdCoreDump)
}
