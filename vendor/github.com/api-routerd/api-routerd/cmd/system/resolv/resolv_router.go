// SPDX-License-Identifier: Apache-2.0

package resolv

import (
	"net/http"

	"github.com/gorilla/mux"
)

func configureResolv(rw http.ResponseWriter, r *http.Request) {
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

// RegisterRouterResolv register with mux
func RegisterRouterResolv(router *mux.Router) {
	r := router.PathPrefix("/resolv").Subrouter().StrictSlash(false)

	r.HandleFunc("/get", configureResolv)
	r.HandleFunc("/add", configureResolv)
	r.HandleFunc("/delete", configureResolv)
}
