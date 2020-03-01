// SPDX-License-Identifier: Apache-2.0

package resolv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/api-routerd/api-routerd/cmd/share"

	log "github.com/sirupsen/logrus"
)

const (
	resolvConfPath = "/etc/resolv.conf"
)

// DNSConfig Json Request
type DNSConfig struct {
	Servers []string `json:"servers"`
	Search  []string `json:"search"`
}

func (conf *DNSConfig) writeConfig() error {
	f, err := os.OpenFile(resolvConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, server := range conf.Servers {
		line := "nameserver " + server
		fmt.Fprintln(w, line)
	}
	for _, s := range conf.Search {
		line := "search " + s
		fmt.Fprintln(w, line)
	}

	w.Flush()

	return nil
}

func readConf() (*DNSConfig, error) {
	lines, err := share.ReadFullFile(resolvConfPath)
	if err != nil {
		log.Errorf("Failed to read: %s", resolvConfPath)
		return nil, err
	}

	conf := new(DNSConfig)

	for _, line := range lines {
		fields := strings.Fields(line)

		switch fields[0] {
		case "nameserver":
			conf.Servers = append(conf.Servers, fields[1])
			break
		case "search":

			for i, search := range strings.Fields(line) {
				if i == 0 {
					continue
				}

				conf.Search = append(conf.Search, search)
			}
		}
	}

	// Don't return nil in json
	if len(conf.Servers) == 0 {
		conf.Servers = []string{""}
	}

	if len(conf.Search) == 0 {
		conf.Search = []string{""}
	}

	return conf, nil
}

// GetConf read resolv.conf and send response
func GetConf(rw http.ResponseWriter) error {
	conf, err := readConf()
	if err != nil {
		return err
	}

	return share.JSONResponse(conf, rw)
}

// UpdateConf update resolv.conf
func UpdateConf(rw http.ResponseWriter, r *http.Request) error {
	dns := DNSConfig{
		Servers: []string{""},
		Search:  []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &dns)
	if err != nil {
		log.Error("Failed to Decode HTTP request to json: ", err)
		return err
	}

	conf, err := readConf()
	if err != nil {
		return err
	}

	// update nameserver
	for _, s := range dns.Servers {
		if share.StringContains(conf.Servers, s) {
			continue
		}

		conf.Servers = append(conf.Servers, s)
	}

	// update domains
	for _, s := range dns.Search {
		if share.StringContains(conf.Search, s) {
			continue
		}
		conf.Search = append(conf.Search, s)
	}

	err = conf.writeConfig()
	if err != nil {
		log.Errorf("Failed Write to resolv conf: %s", err)
		return err
	}

	return share.JSONResponse(conf, rw)
}

// DeleteConf delete conf from file
func DeleteConf(rw http.ResponseWriter, r *http.Request) error {
	dns := DNSConfig{
		Servers: []string{""},
		Search:  []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &dns)
	if err != nil {
		log.Error("Failed to Decode HTTP request to json: ", err)
		return err
	}

	conf, err := readConf()
	if err != nil {
		return err
	}

	// update nameserver
	for _, s := range dns.Servers {
		if !share.StringContains(conf.Servers, s) {
			continue
		}

		conf.Servers, _ = share.StringDeleteSlice(conf.Servers, s)
	}

	// update domains
	for _, s := range dns.Search {
		if !share.StringContains(conf.Search, s) {
			continue
		}

		conf.Search, _ = share.StringDeleteSlice(conf.Search, s)
	}

	err = conf.writeConfig()
	if err != nil {
		log.Errorf("Failed Write to resolv conf: %s", err)
		return err
	}

	return share.JSONResponse(conf, rw)
}
