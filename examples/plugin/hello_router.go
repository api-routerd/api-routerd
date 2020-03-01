// SPDX-License-Identifier: Apache-2.0

package hello

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routerSayHello(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		g := new(Hello)

		vars := mux.Vars(r)
		text := vars["text"]

		g.Text = text

		err := g.SayHello(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

// RegisterRouterSayHello register with mux
func RegisterRouterSayHello(router *mux.Router) {
	s := router.PathPrefix("/hello").Subrouter().StrictSlash(false)

	s.HandleFunc("/sayhello/{text}", routerSayHello)
}
