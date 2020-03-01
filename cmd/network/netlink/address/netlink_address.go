// SPDX-License-Identifier: Apache-2.0

package address

import (
        "encoding/json"
        "net/http"

        "github.com/api-routerd/api-routerd/cmd/share"

        "github.com/vishvananda/netlink"
)

// Address JSON request
type Address struct {
        Action  string `json:"action"`
        Link    string `json:"link"`
        Address string `json:"address"`
        Label   string `json:"label"`
}

// DecodeJSONRequest parses JSON request to Address type
func DecodeJSONRequest(r *http.Request) (Address, error) {
        address := new(Address)

        err := json.NewDecoder(r.Body).Decode(&address)
        if err != nil {
                return *address, err
        }

        return *address, nil
}

// Add Add a new address to interface
func (a *Address) Add() error {
        link, err := netlink.LinkByName(a.Link)
        if err != nil {
                return err
        }

        addr, err := netlink.ParseAddr(a.Address)
        if err != nil {
                return err
        }

        err = netlink.AddrAdd(link, addr)
        if err != nil {
                return err
        }

        return nil
}

// Del remove a address from interface
func (a *Address) Del() error {
        link, err := netlink.LinkByName(a.Link)
        if err != nil {
                return err
        }

        addr, err := netlink.ParseAddr(a.Address)
        if err != nil {
                return err
        }

        err = netlink.AddrDel(link, addr)
        if err != nil {
                return err
        }

        return nil
}

// Get link address
func (a *Address) Get(rw http.ResponseWriter) error {
        if a.Link != "" {
                link, err := netlink.LinkByName(a.Link)
                if err != nil {
                        return err
                }

                addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
                if err != nil {
                        return err
                }

                return share.JSONResponse(addrs, rw)

        }

        addrs, err := netlink.AddrList(nil, netlink.FAMILY_ALL)
        if err != nil {
                return err
        }

        return share.JSONResponse(addrs, rw)
}
