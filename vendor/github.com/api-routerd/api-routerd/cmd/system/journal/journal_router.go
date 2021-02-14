// SPDX-License-Identifier: Apache-2.0

package journal

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routerConfigureJournalConf(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break

	case "POST":
		err := UpdateConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

// RegisterRouterJournal register with mux
func RegisterRouterJournal(router *mux.Router) {
	j := router.PathPrefix("/journal").Subrouter().StrictSlash(false)

	// conf
	j.HandleFunc("/conf", routerConfigureJournalConf)
	j.HandleFunc("/conf/update", routerConfigureJournalConf)
}
