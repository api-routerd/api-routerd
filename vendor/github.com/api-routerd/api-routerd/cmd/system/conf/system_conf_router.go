// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"net/http"

	"github.com/gorilla/mux"
)

func readSudoersConfig(rw http.ResponseWriter, req *http.Request) {
	err := GetSudoers(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func readSSHConfig(rw http.ResponseWriter, req *http.Request) {
	err := SSHConfFileRead(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RegisterRouterSystemConf register with mux
func RegisterRouterSystemConf(router *mux.Router) {
	c := router.PathPrefix("/conf").Subrouter().StrictSlash(false)

	// Generic system confs
	c.HandleFunc("/sudoers", readSudoersConfig)
	c.HandleFunc("/sshd", readSSHConfig)
}
