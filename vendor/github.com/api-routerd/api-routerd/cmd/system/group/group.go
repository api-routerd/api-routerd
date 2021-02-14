// SPDX-License-Identifier: Apache-2.0

package group

import (
	"fmt"
	"os/exec"
	"os/user"

	log "github.com/sirupsen/logrus"
)

// Group Json Commands
type Group struct {
	Gid     string `json:"gid"`
	Name    string `json:"name"`
	NewName string `json:"new_name"`
}

// GroupAdd Add group
func (r *Group) GroupAdd() error {
	g, err := user.LookupGroup(r.Name)
	if err != nil {
		_, ok := err.(user.UnknownGroupError)
		if !ok {
			return err
		}
	}

	if g != nil {
		return fmt.Errorf("Failed to add group. Group '%s' already exists", r.Name)
	}

	id, err := user.LookupGroupId(r.Gid)
	if err != nil {
		_, ok := err.(user.UnknownGroupIdError)
		if !ok {
			return err
		}
	}

	if id != nil {
		return fmt.Errorf("Failed to add group '%s': Gid '%s' exists", r.Name, r.Gid)
	}

	path, err := exec.LookPath("groupadd")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, r.Name, "-g", r.Gid)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to add group %s: %s", r.Name, stdout)
		return fmt.Errorf("Failed to add group '%s': %s", r.Name, stdout)
	}

	return nil
}

// GroupDel delete a group
func (r *Group) GroupDel() error {
	g, err := user.LookupGroup(r.Name)
	if err != nil {
		_, ok := err.(user.UnknownGroupError)
		if !ok {
			return err
		}
	}

	if g == nil {
		return fmt.Errorf("Failed to delete group '%s'. Group does not exists", r.Name)
	}

	path, err := exec.LookPath("groupdel")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, r.Name)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to delete group %s: %s", r.Name, stdout)
		return fmt.Errorf("Failed to delete group '%s': %s", r.Name, stdout)
	}

	return nil
}

// GroupModify modify a group
func (r *Group) GroupModify() error {
	g, err := user.LookupGroup(r.Name)
	if err != nil {
		_, ok := err.(user.UnknownGroupError)
		if !ok {
			return err
		}
	}

	if g == nil {
		return fmt.Errorf("Failed to Modify group '%s'. Group does not exists", r.Name)
	}

	g, err = user.LookupGroup(r.NewName)
	if err != nil {
		_, ok := err.(user.UnknownGroupError)
		if !ok {
			return err
		}
	}

	if g != nil {
		return fmt.Errorf("Failed to Modify group '%s'. New Group '%s' already exists", r.Name, r.NewName)
	}

	path, err := exec.LookPath("groupmod")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "-n", r.NewName, r.Name)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to modify group %s: %s", r.Name, stdout)
		return fmt.Errorf("Failed to modify group '%s': %s", r.Name, stdout)
	}

	return nil
}
