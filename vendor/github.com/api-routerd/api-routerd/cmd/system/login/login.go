// SPDX-License-Identifier: Apache-2.0

package login

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/api-routerd/api-routerd/cmd/share"
	sd "github.com/coreos/go-systemd/login1"
)

const (
	dbusInterface = "org.freedesktop.login1.Manager"
	dbusPath      = "/org/freedesktop/login1"
)

var loginMethod = map[string]string{
	"list-sessions":     "ListSessions",
	"list-users":        "ListUsers",
	"lock-session":      "LockSession",
	"lock-sessions":     "LockSessions",
	"terminate-session": "TerminateSession",
	"terminate-user":    "TerminateUser",
}

// Login JSON message
type Login struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

// LoginMethodGet Pull properties from systemd
func (t *Login) LoginMethodGet(rw http.ResponseWriter) error {
	c, err := sd.New()
	if err != nil {
		return err
	}
	defer c.Close()

	_, k := loginMethod[t.Path]
	if !k {
		return fmt.Errorf("Failed to call method login:  %s not found", t.Path)
	}

	switch loginMethod[t.Path] {
	case "ListUsers":
		users, err := c.ListUsers()
		if err != nil {
			return err
		}

		return share.JSONResponse(users, rw)
	case "ListSessions":
		sessions, err := c.ListSessions()
		if err != nil {
			return err
		}

		return share.JSONResponse(sessions, rw)
	}

	return nil
}

// LoginMethodPost Do call login methods via dbus
func (t *Login) LoginMethodPost(rw http.ResponseWriter) error {
	c, err := sd.New()
	if err != nil {
		return err
	}
	defer c.Close()

	_, k := loginMethod[t.Path]
	if !k {
		return fmt.Errorf("Failed to call method login:  %s not found", t.Path)
	}

	switch loginMethod[t.Path] {
	case "LockSession":
		c.LockSession(t.Value)
		break
	case "LockSessions":
		c.LockSessions()
		break
	case "TerminateSession":
		c.TerminateSession(t.Value)
		break
	case "TerminateUser":
		v, err := strconv.ParseInt(t.Value, 10, 32)
		if err != nil {
			return err
		}

		c.TerminateUser(uint32(v))
		if err != nil {
			return err
		}
	}

	return nil
}
