// SPDX-License-Identifier: Apache-2.0

package sysctl

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func routerSysctlGet(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := Get(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func routerSysctlUpdate(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST", "PUT":
		s := new(Sysctl)
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.Update()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func routerSysctlDelete(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		s := new(Sysctl)
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.Delete()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

//RegisterRouterSysctl register with mux
func RegisterRouterSysctl(router *mux.Router) {
	s := router.PathPrefix("/sysctl").Subrouter().StrictSlash(false)

	s.HandleFunc("/get", routerSysctlGet)
	s.HandleFunc("/add", routerSysctlUpdate)
	s.HandleFunc("/modify", routerSysctlUpdate)
	s.HandleFunc("/delete", routerSysctlDelete)
}
