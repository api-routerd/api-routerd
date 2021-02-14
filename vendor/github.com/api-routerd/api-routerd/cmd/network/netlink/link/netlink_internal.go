// SPDX-License-Identifier: Apache-2.0

package link

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func (r *Link) setMasterBridge() error {
	bridge, err := netlink.LinkByName(r.Link)
	if err != nil {
		log.Errorf("Failed to find bridge link %s: %v", r.Link, err)
		return err
	}

	br, b := bridge.(*netlink.Bridge)
	if !b {
		log.Errorf("Link '%s'is not a bridge: %v", r.Link, err)
		return fmt.Errorf("Link is not a bridge")
	}

	for _, n := range r.Enslave {
		link, err := netlink.LinkByName(n)
		if err != nil {
			log.Errorf("Failed to find slave link %s: %v", n, err)
			continue
		}

		err = netlink.LinkSetMaster(link, br)
		if err != nil {
			log.Errorf("Failed to set link %s master device %s: %v", n, r.Link, err)
		}
	}

	return nil
}

func (r *Link) createBridge() error {
	_, err := netlink.LinkByName(r.Link)
	if err == nil {
		log.Infof("Bridge link %s exists. Using the bridge", r.Link)
	} else {

		bridge := &netlink.Bridge{
			LinkAttrs: netlink.LinkAttrs{
				Name: r.Link,
			},
		}
		err = netlink.LinkAdd(bridge)
		if err != nil {
			log.Errorf("Failed to create bridge %s: %v", r.Link, err)
			return err
		}

		log.Debugf("Successfully create bridge link: %s", r.Link)
	}

	return r.setMasterBridge()
}

func (r *Link) setMasterBond() error {
	bond, err := netlink.LinkByName(r.Link)
	if err != nil {
		log.Errorf("Failed to find bond link %s: %v", r.Link, err)
		return err
	}

	for _, n := range r.Enslave {
		link, err := netlink.LinkByName(n)
		if err != nil {
			log.Errorf("Failed to find slave link %s: %v", n, err)
			continue
		}

		err = netlink.LinkSetBondSlave(link, &netlink.Bond{LinkAttrs: *bond.Attrs()})
		if err != nil {
			log.Errorf("Failed to set link %s master device %s: %v", n, r.Link, err)
		}
	}

	return nil
}

func (r *Link) createBond() error {
	_, err := netlink.LinkByName(r.Link)
	if err == nil {
		log.Infof("Bond link %s exists. Using the bond", r.Link)
	} else {

		bond := netlink.NewLinkBond(
			netlink.LinkAttrs{
				Name: r.Link,
			},
		)

		bond.Mode = netlink.StringToBondModeMap[r.Mode]
		err = netlink.LinkAdd(bond)
		if err != nil {
			log.Errorf("Failed to create bond %s: %v", r.Link, err)
			return err
		}

		log.Debugf("Successfully create bond link: %s", r.Link)
	}

	return r.setMasterBond()
}

func setUp(link string) error {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %v", link, err)
		return err
	}

	err = netlink.LinkSetUp(l)
	if err != nil {
		log.Errorf("Failed to set link %s up: %v", l, err)
		return err
	}

	return nil
}

func setDown(link string) error {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %v", link, err)
		return err
	}

	err = netlink.LinkSetDown(l)
	if err != nil {
		log.Errorf("Failed to set link down %s: %v", l, err)
		return err
	}

	return nil
}

func setMTU(link string, mtu int) error {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %v", link, err)
		return err
	}

	err = netlink.LinkSetMTU(l, mtu)
	if err != nil {
		log.Errorf("Failed to set link %s MTU %d: %v", link, mtu, err)
		return err
	}

	return nil
}
