// SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"net/http"

	"github.com/gorilla/mux"
)

func configureSystemdResolved(rw http.ResponseWriter, r *http.Request) {
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

// RegisterRouterResolve register with mux
func RegisterRouterResolve(router *mux.Router) {
	r := router.PathPrefix("/systemdresolved").Subrouter().StrictSlash(false)

	// systemd-resolved
	r.HandleFunc("/get", configureSystemdResolved)
	r.HandleFunc("/add", configureSystemdResolved)
	r.HandleFunc("/delete", configureSystemdResolved)
}
