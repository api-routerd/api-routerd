// SPDX-License-Identifier: Apache-2.0

package networkctl

import (
	"net/http"

	"fmt"

	"github.com/api-routerd/api-routerd/cmd/share"
)

// Networkctl JSON message
type Networkctl struct {
	Verb string `json:"verb"`
	Link string `json:"link"`
}

// NetworkctlGet collect info via networkctl
func (n *Networkctl) NetworkctlGet(rw http.ResponseWriter) error {
	link := share.LinkExists(n.Link)
	if !link {
		return fmt.Errorf("Failed to find link: %s", n.Link)
	}

	switch n.Verb {
	case "status":
		r, err := ExecuteNetworkctlStatus(n.Link)
		if err != nil {
			return err
		}

		return share.JSONResponse(r, rw)

	case "lldp":
		r, err := ExecuteNetworkctlLLDP()
		if err != nil {
			return err
		}

		return share.JSONResponse(r, rw)

	case "networkctl":
	default:
		r, err := ExecuteNetworkctl()
		if err != nil {
			return err
		}

		return share.JSONResponse(r, rw)
	}

	return nil
}
