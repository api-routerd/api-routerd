package timesync

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
	timeSyncdConfPath = "/etc/systemd/timesyncd.conf"
)

// TimeSyncConfig Json request
type TimeSyncConfig struct {
	NTP                []string `json:"NTP"`
	FallbackNTP        []string `json:"FallbackNTP"`
	RootDistanceMaxSec string   `json:"RootDistanceMaxSec"`
	PollIntervalMinSec string   `json:"PollIntervalMinSec"`
	PollIntervalMaxSec string   `json:"PollIntervalMaxSec"`
}

func (t *TimeSyncConfig) writeConf() error {
	f, err := os.OpenFile(timeSyncdConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	conf := "[Time]\n"

	ntpConf := "NTP="
	for _, s := range t.NTP {
		ntpConf += s + " "
	}
	conf += ntpConf + "\n"

	FallbackNTP := "FallbackNTP="
	for _, s := range t.FallbackNTP {
		FallbackNTP += s + " "
	}
	conf += FallbackNTP + "\n"

	if t.RootDistanceMaxSec != "" {
		conf += "RootDistanceMaxSec=" + t.RootDistanceMaxSec + "\n"
	}

	if t.PollIntervalMinSec != "" {
		conf += "PollIntervalMinSec=" + t.PollIntervalMinSec + "\n"
	}

	if t.PollIntervalMaxSec != "" {
		conf += "PollIntervalMaxSec=" + t.PollIntervalMaxSec + "\n"
	}

	fmt.Fprintln(w, conf)
	w.Flush()

	return nil
}

func readConf() (*TimeSyncConfig, error) {
	cfg, err := ini.Load(timeSyncdConfPath)
	if err != nil {
		return nil, err
	}

	conf := new(TimeSyncConfig)

	ntp := cfg.Section("Time").Key("NTP").String()
	FallbackNTP := cfg.Section("Time").Key("FallbackNTP").String()

	conf.NTP = strings.Fields(ntp)
	conf.FallbackNTP = strings.Fields(FallbackNTP)

	return conf, nil
}

// GetConf read from file and send response
func GetConf(rw http.ResponseWriter) error {
	conf, err := readConf()
	if err != nil {
		return err
	}

	return share.JSONResponse(conf, rw)
}

// UpdateConf update timesync conf
func UpdateConf(rw http.ResponseWriter, r *http.Request) error {
	t := TimeSyncConfig{
		NTP:         []string{""},
		FallbackNTP: []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to parse HTTP request: %v", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &t)
	if err != nil {
		log.Errorf("Failed to Decode HTTP request to json: %v", err)
		return err
	}

	conf, err := readConf()
	if err != nil {
		return err
	}

	// update NTP
	for _, s := range t.NTP {
		if share.StringContains(conf.NTP, s) {
			continue
		}

		conf.NTP = append(conf.NTP, s)
	}

	// update fallback
	for _, s := range t.FallbackNTP {
		if share.StringContains(conf.FallbackNTP, s) {
			continue
		}
		conf.FallbackNTP = append(conf.FallbackNTP, s)
	}

	if t.RootDistanceMaxSec != "" {
		conf.RootDistanceMaxSec = t.RootDistanceMaxSec
	}

	if t.PollIntervalMinSec != "" {
		conf.PollIntervalMinSec = t.PollIntervalMinSec
	}

	if t.PollIntervalMaxSec != "" {
		conf.PollIntervalMaxSec = "t.PollIntervalMaxSec"
	}

	err = conf.writeConf()
	if err != nil {
		log.Errorf("Failed Write to time sync conf: %v", err)
		return err
	}

	return share.JSONResponse(conf, rw)
}

// DeleteConf remove conf from file
func DeleteConf(rw http.ResponseWriter, r *http.Request) error {
	t := TimeSyncConfig{
		NTP:         []string{""},
		FallbackNTP: []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to parse HTTP request: %v", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &t)
	if err != nil {
		log.Errorf("Failed to Decode HTTP request to json: %v", err)
		return err
	}

	conf, err := readConf()
	if err != nil {
		return err
	}

	// update NTP
	for _, s := range t.NTP {
		if !share.StringContains(conf.NTP, s) {
			continue
		}

		conf.NTP, _ = share.StringDeleteSlice(conf.NTP, s)
	}

	// update Fallback
	for _, s := range t.FallbackNTP {
		if !share.StringContains(conf.FallbackNTP, s) {
			continue
		}

		conf.FallbackNTP, _ = share.StringDeleteSlice(conf.FallbackNTP, s)
	}

	if t.RootDistanceMaxSec != "" {
		conf.RootDistanceMaxSec = t.RootDistanceMaxSec
	}

	if t.PollIntervalMinSec != "" {
		conf.PollIntervalMinSec = t.PollIntervalMinSec
	}

	if t.PollIntervalMaxSec != "" {
		conf.PollIntervalMaxSec = "t.PollIntervalMaxSec"
	}

	err = conf.writeConf()
	if err != nil {
		log.Errorf("Failed Write to time sync conf: %v", err)
		return err
	}

	return share.JSONResponse(conf, rw)
}
