// SPDX-License-Identifier: Apache-2.0

package kmod

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func routerGetModules(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := LsMod(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func routerModProbe(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		kmod := new(KMod)
		err := json.NewDecoder(r.Body).Decode(&kmod)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = kmod.ModProbe()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func routerRmMod(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":

		kmod := new(KMod)
		err := json.NewDecoder(r.Body).Decode(&kmod)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = kmod.RmMod()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

// RegisterRouterKMod register with mux
func RegisterRouterKMod(router *mux.Router) {
	s := router.PathPrefix("/kmod").Subrouter().StrictSlash(false)

	s.HandleFunc("/lsmod", routerGetModules)
	s.HandleFunc("/modprobe", routerModProbe)
	s.HandleFunc("/rmmod", routerRmMod)
}
