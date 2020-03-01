// SPDX-License-Identifier: Apache-2.0

package login

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func routerLoginMethodGet(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		vars := mux.Vars(r)
		path := vars["path"]

		login := &Login{
			Path: path,
		}

		err := login.LoginMethodGet(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		break
	}
}

func routerLoginMethodPost(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		login := new(Login)
		err := json.NewDecoder(r.Body).Decode(&login)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)
		login.Path = vars["path"]

		err = login.LoginMethodPost(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		rw.WriteHeader(http.StatusOK)

		break
	}
}

// RegisterRouterLogin register with mux
func RegisterRouterLogin(router *mux.Router) {
	s := router.PathPrefix("/login").Subrouter().StrictSlash(false)

	s.HandleFunc("/get/{path}", routerLoginMethodGet)
	s.HandleFunc("/post/{path}", routerLoginMethodPost)
}
