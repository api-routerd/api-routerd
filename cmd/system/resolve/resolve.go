// SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/api-routerd/api-routerd/cmd/share"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

const (
	resolvedConfPath = "/etc/systemd/resolved.conf"
)

// DNSConfig Json request and response
type DNSConfig struct {
	DNS         []string `json:"dns"`
	FallbackDNS []string `json:"fallback_dns"`
}

func (d *DNSConfig) writeConfig() error {
	f, err := os.OpenFile(resolvedConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	conf := "[Resolve]\n"

	dnsConf := "DNS="
	for _, s := range d.DNS {
		dnsConf += s + " "
	}
	conf += dnsConf + "\n"

	fallbackDNS := "FallbackDNS="
	for _, s := range d.FallbackDNS {
		fallbackDNS += s + " "
	}
	conf += fallbackDNS + "\n"

	fmt.Fprintln(w, conf)
	w.Flush()

	return nil
}

func readConf() (*DNSConfig, error) {
	cfg, err := ini.Load(resolvedConfPath)
	if err != nil {
		return nil, err
	}

	conf := new(DNSConfig)

	dns := cfg.Section("Resolve").Key("DNS").String()
	fallbackDNS := cfg.Section("Resolve").Key("FallbackDNS").String()

	conf.DNS = strings.Fields(dns)
	conf.FallbackDNS = strings.Fields(fallbackDNS)

	return conf, nil
}

// GetConf read conf and send response
func GetConf(rw http.ResponseWriter) error {
	conf, err := readConf()
	if err != nil {
		return err
	}

	return share.JSONResponse(conf, rw)
}

// UpdateConf update conf
func UpdateConf(rw http.ResponseWriter, r *http.Request) error {
	dns := DNSConfig{
		DNS:         []string{""},
		FallbackDNS: []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to parse HTTP request: %v", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &dns)
	if err != nil {
		log.Errorf("Failed to Decode HTTP request to json: %v", err)
		return err
	}

	conf, err := readConf()
	if err != nil {
		return err
	}

	// update DNS
	for _, s := range dns.DNS {
		if share.StringContains(conf.DNS, s) {
			continue
		}

		conf.DNS = append(conf.DNS, s)
	}

	// update fallback
	for _, s := range dns.FallbackDNS {
		if share.StringContains(conf.FallbackDNS, s) {
			continue
		}
		conf.FallbackDNS = append(conf.FallbackDNS, s)
	}

	err = conf.writeConfig()
	if err != nil {
		log.Errorf("Failed Write to resolv conf: %v", err)
		return err
	}

	return share.JSONResponse(conf, rw)
}

// DeleteConf remove conf from file
func DeleteConf(rw http.ResponseWriter, r *http.Request) error {
	dns := DNSConfig{
		DNS:         []string{""},
		FallbackDNS: []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to parse HTTP request: %v", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &dns)
	if err != nil {
		log.Errorf("Failed to Decode HTTP request to json: %v", err)
		return err
	}

	conf, err := readConf()
	if err != nil {
		return err
	}

	// update DNS
	for _, s := range dns.DNS {
		if !share.StringContains(conf.DNS, s) {
			continue
		}

		conf.DNS, _ = share.StringDeleteSlice(conf.DNS, s)
	}

	// update Fallback
	for _, s := range dns.FallbackDNS {
		if !share.StringContains(conf.FallbackDNS, s) {
			continue
		}

		conf.FallbackDNS, _ = share.StringDeleteSlice(conf.FallbackDNS, s)
	}

	err = conf.writeConfig()
	if err != nil {
		log.Errorf("Failed Write to resolv conf: %v", err)
		return err
	}

	return share.JSONResponse(conf, rw)
}
