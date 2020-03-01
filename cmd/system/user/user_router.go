// SPDX-License-Identifier: Apache-2.0

package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func routerAdd(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		u := new(User)
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = u.Add()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func routerModify(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":

		u := new(User)
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = u.Modify()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

func routerDel(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":

		u := new(User)
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = u.Del()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		break
	}

	rw.WriteHeader(http.StatusOK)
}

//RegisterRouterUser register with mux
func RegisterRouterUser(router *mux.Router) {
	s := router.PathPrefix("/user").Subrouter().StrictSlash(false)

	s.HandleFunc("/add", routerAdd)
	s.HandleFunc("/delete", routerDel)
	s.HandleFunc("/modify", routerModify)
}
