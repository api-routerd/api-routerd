// SPDX-License-Identifier: Apache-2.0

package timedate

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func routerGetTimeDate(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	property := vars["property"]

	switch r.Method {
	case "GET":

		timedate := new(TimeDate)
		if property != "" {
			timedate.Property = property
		}

		err := GetTimeDate(rw, property)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerSetTimeDate(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":

		timedate := new(TimeDate)
		err := json.NewDecoder(r.Body).Decode(&timedate)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = timedate.SetTimeDate()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

// RegisterRouterTimeDate register with mux
func RegisterRouterTimeDate(router *mux.Router) {
	t := router.PathPrefix("/timedate").Subrouter().StrictSlash(false)

	t.HandleFunc("", routerGetTimeDate)
	t.HandleFunc("/get/{property}", routerGetTimeDate)
	t.HandleFunc("/set", routerSetTimeDate)
}
