// SPDX-License-Identifier: Apache-2.0

package user

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"

	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

const (
	userFile = "/run/api-routerd-users"
)

//User Json request
type User struct {
	UID           string   `json:"uid"`
	Gid           string   `json:"gid"`
	Groups        []string `json:"groups"`
	Comment       string   `json:"comment"`
	HomeDirectory string   `json:"home_dir"`
	Shell         string   `json:"shell"`
	Username      string   `json:"username"`
	Password      string   `json:"password"`
}

//Add add user
func (r *User) Add() error {
	u, err := user.Lookup(r.Username)
	if err != nil {
		_, ok := err.(user.UnknownUserError)
		if !ok {
			return err
		}
	}

	if u != nil {
		return fmt.Errorf("Failed to add user '%s' already exists", r.Username)
	}

	if r.UID != "" {
		id, err := user.LookupId(r.UID)
		if err != nil {
			_, ok := err.(user.UnknownUserIdError)
			if !ok {
				return err
			}
		}

		if id != nil {
			return fmt.Errorf("Failed to add user '%s': Gid '%s' exists", r.Username, r.Gid)
		}
	}

	//<Username>:<Password>:<UID>:<GID>:<User Info>:<Home Dir>:<Default Shell>
	line := r.Username + ":" + r.Password + ":" + r.UID + ":" + r.Gid + ":" + r.Comment + ":" + r.HomeDirectory + ":" + r.Shell

	err = share.WriteOneLineFile(userFile, line)
	if err != nil {
		return err
	}

	path, err := exec.LookPath("newusers")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, userFile)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to add user %s: %s", r.Username, stdout)
		return fmt.Errorf("Failed to add user '%s': %s", r.Username, stdout)
	}

	return os.Remove(userFile)
}

//Del delete user
func (r *User) Del() error {
	g, err := user.Lookup(r.Username)
	if err != nil {
		_, ok := err.(user.UnknownUserError)
		if !ok {
			return err
		}
	}

	if g == nil {
		return fmt.Errorf("Failed to delete user '%s'. User does not exists", r.Username)
	}

	path, err := exec.LookPath("userdel")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, r.Username)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to delete user %s: %s", r.Username, stdout)
		return fmt.Errorf("Failed to delete user '%s': %s", r.Username, stdout)
	}

	return nil
}

//Modify user
func (r *User) Modify() error {
	g, err := user.Lookup(r.Username)
	if err != nil {
		_, ok := err.(user.UnknownUserError)
		if !ok {
			return err
		}
	}

	if g == nil {
		return fmt.Errorf("Failed to Modify user '%s'. User does not exists", r.Username)
	}

	path, err := exec.LookPath("usermod")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "-G", r.Groups[0], r.Username)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to modify user %s: %s", r.Username, stdout)
		return fmt.Errorf("Failed to modify user '%s': %s", r.Username, stdout)
	}

	return nil
}
