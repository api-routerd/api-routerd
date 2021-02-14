// SPDX-License-Identifier: Apache-2.0

package ethtool

import (
	"errors"
	"net/http"
	"unsafe"

	"github.com/api-routerd/api-routerd/cmd/share"

	"github.com/safchain/ethtool"
	log "github.com/sirupsen/logrus"
)

// Ethtool JSON message
type Ethtool struct {
	Action   string `json:"action"`
	Link     string `json:"link"`
	Property string `json:"property"`
	Value    string `json:"Value"`
}

// GetEthTool collect info via ethtool ioctl
func (r *Ethtool) GetEthTool(rw http.ResponseWriter) error {
	link := share.LinkExists(r.Link)
	if !link {
		log.Errorf("Failed to get link: %s", r.Link)
		return errors.New("Link not found")
	}

	e, err := ethtool.NewEthtool()
	if err != nil {
		log.Errorf("Failed to init ethtool for link %v: %s", err, r.Link)
		return err
	}
	defer e.Close()

	switch r.Action {
	case "get-link-stat":
		stats, err := e.Stats(r.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool statitics for link %v: %s", err, r.Link)
			return err
		}

		return share.JSONResponse(stats, rw)

	case "get-link-features":

		features, err := e.Features(r.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool features for link %v: %s", err, r.Link)
			return err
		}

		return share.JSONResponse(features, rw)

	case "get-link-bus":

		bus, err := e.BusInfo(r.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool bus for link %v: %s", err, r.Link)
			return err
		}

		b := struct {
			Bus string
		}{
			bus,
		}

		return share.JSONResponse(b, rw)

	case "get-link-driver-name":

		driver, err := e.DriverName(r.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool driver name for link %v: %s", err, r.Link)
			return err
		}

		d := struct {
			Driver string
		}{
			driver,
		}

		return share.JSONResponse(d, rw)

	case "get-link-driver-info":

		e, err := NewEthTool()
		if err != nil {
			log.Errorf("Failed to init ethtool for link %v: %s", err, r.Link)
			return err
		}
		defer e.Close()

		drvinfo := DrvInfo{
			Cmd: EthtoolGDRVInfo,
		}

		err = e.ioctl(r.Link, uintptr(unsafe.Pointer(&drvinfo)))
		if err != nil {
			return err
		}

		return share.JSONResponse(drvinfo, rw)

	case "get-link-permaddr":

		permaddr, err := e.PermAddr(r.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool Perm Addr for link %v: %s", err, r.Link)
			return err
		}

		p := struct {
			PermAddr string
		}{
			permaddr,
		}

		return share.JSONResponse(p, rw)

	case "get-link-eeprom":

		eeprom, err := e.ModuleEepromHex(r.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool eeprom for link %v: %s", err, r.Link)
			return err
		}

		e := struct {
			ModuleEeprom string
		}{
			eeprom,
		}

		return share.JSONResponse(e, rw)

	case "get-link-msglvl":

		msglvl, err := e.MsglvlGet(r.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool msglvl for link %v: %s", err, r.Link)
			return err
		}

		g := struct {
			ModuleMsglv uint32
		}{
			msglvl,
		}

		return share.JSONResponse(g, rw)

	case "get-link-mapped":

		mapped, err := e.CmdGetMapped(r.Link)
		if err != nil {
			log.Errorf("Failed to get ethtool msglvl for link %v: %s", err, r.Link)
			return err
		}

		return share.JSONResponse(mapped, rw)
	}

	return nil
}

// SetEthTool set ethtool info
func (r *Ethtool) SetEthTool(rw http.ResponseWriter) error {
	link := share.LinkExists(r.Link)
	if !link {
		log.Errorf("Failed to get link: %s", r.Link)
		return errors.New("Link not found")
	}

	e, err := ethtool.NewEthtool()
	if err != nil {
		log.Errorf("Failed to init ethtool for link %v: %s", err, r.Link)
		return err
	}
	defer e.Close()

	switch r.Action {
	case "set-link-feature":

		feature := make(map[string]bool)

		b, err := share.ParseBool(r.Value)
		feature[r.Property] = b

		err = e.Change(r.Link, feature)
		if err != nil {
			return err
		}
	}

	return nil
}
