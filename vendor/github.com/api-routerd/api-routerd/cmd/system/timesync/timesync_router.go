// SPDX-License-Identifier: Apache-2.0

package timesync

import (
	"net/http"

	"github.com/gorilla/mux"
)

func configureSystemdTimeSyncd(rw http.ResponseWriter, r *http.Request) {
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

// RegisterRouterTimeSync register with mux
func RegisterRouterTimeSync(router *mux.Router) {
	t := router.PathPrefix("/systemdtimesyncd").Subrouter().StrictSlash(false)

	t.HandleFunc("/get", configureSystemdTimeSyncd)
	t.HandleFunc("/add", configureSystemdTimeSyncd)
	t.HandleFunc("/delete", configureSystemdTimeSyncd)
}
