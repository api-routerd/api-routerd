// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/api-routerd/api-routerd/cmd/share"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

const (
	systemConfPath = "/etc/systemd/system.conf"
)

var systemConfig = map[string]string{}

func writeSystemConfig() error {
	f, err := os.OpenFile(systemConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	conf := "[Manager]\n"
	for k, v := range systemConfig {
		if v != "" {
			conf += k + "=" + v
		} else {
			conf += "#" + k + "="
		}
		conf += "\n"
	}

	fmt.Fprintln(w, conf)
	w.Flush()

	return nil
}

func readSystemConf() error {
	cfg, err := ini.Load(systemConfPath)
	if err != nil {
		return err
	}

	for k := range systemConfig {
		systemConfig[k] = cfg.Section("Manager").Key(k).String()
	}

	return nil
}

// GetSystemConf read system.conf
func GetSystemConf(rw http.ResponseWriter) error {
	err := readSystemConf()
	if err != nil {
		return err
	}

	return share.JSONResponse(systemConfig, rw)
}

// UpdateSystemConf update the system.conf
func UpdateSystemConf(rw http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to parse HTTP request: %v", err)
		return err
	}

	conf := make(map[string]string)
	err = json.Unmarshal([]byte(body), &conf)
	if err != nil {
		log.Errorf("Failed to Decode HTTP request to json: %v", err)
		return err
	}

	err = readSystemConf()
	if err != nil {
		return err
	}

	for k, v := range conf {
		_, ok := systemConfig[k]
		if ok {
			systemConfig[k] = v
		}
	}

	err = writeSystemConfig()
	if err != nil {
		log.Errorf("Failed Write to system conf: %v", err)
		return err
	}

	return share.JSONResponse(systemConfig, rw)
}

// InitSystemd Init systemd conf
func InitSystemd() {
	systemConfig["LogLevel"] = ""
	systemConfig["LogTarget"] = ""
	systemConfig["LogColor"] = ""
	systemConfig["LogLocation"] = ""
	systemConfig["DumpCore"] = ""
	systemConfig["ShowStatus"] = ""
	systemConfig["CrashChangeVT"] = ""
	systemConfig["CrashShell"] = ""
	systemConfig["CrashReboot"] = ""
	systemConfig["CtrlAltDelBurstAction"] = ""
	systemConfig["CPUAffinity"] = ""
	systemConfig["JoinControllers"] = ""
	systemConfig["RuntimeWatchdogSec"] = ""
	systemConfig["ShutdownWatchdogSec"] = ""
	systemConfig["CapabilityBoundingSe"] = ""
	systemConfig["SystemCallArchitectures"] = ""
	systemConfig["TimerSlackNSec"] = ""
	systemConfig["DefaultTimerAccuracySec"] = ""
	systemConfig["DefaultStandardOutput"] = ""
	systemConfig["DefaultStandardError"] = ""
	systemConfig["DefaultTimeoutStartSec"] = ""
	systemConfig["DefaultTimeoutStopSec"] = ""
	systemConfig["DefaultRestartSec"] = ""
	systemConfig["DefaultStartLimitIntervalSec"] = ""
	systemConfig["DefaultStartLimitBurst"] = ""
	systemConfig["DefaultEnvironment"] = ""
	systemConfig["DefaultCPUAccounting"] = ""
	systemConfig["DefaultIOAccounting"] = ""
	systemConfig["DefaultIPAccounting"] = ""
	systemConfig["DefaultBlockIOAccounting"] = ""
	systemConfig["DefaultMemoryAccounting"] = ""
	systemConfig["DefaultTasksAccounting"] = ""
	systemConfig["DefaultTasksMax"] = ""
	systemConfig["DefaultLimitCPU"] = ""
	systemConfig["DefaultLimitFSIZE"] = ""
	systemConfig["DefaultLimitDATA"] = ""
	systemConfig["DefaultLimitSTACK"] = ""
	systemConfig["DefaultLimitCORE"] = ""
	systemConfig["DefaultLimitRSS"] = ""
	systemConfig["DefaultLimitNOFILE"] = ""
	systemConfig["DefaultLimitAS"] = ""
	systemConfig["DefaultLimitNPROC"] = ""
	systemConfig["DefaultLimitMEMLOCK"] = ""
	systemConfig["DefaultLimitLOCKS"] = ""
	systemConfig["DefaultLimitSIGPENDING"] = ""
	systemConfig["DefaultLimitMSGQUEUE"] = ""
	systemConfig["DefaultLimitNICE"] = ""
	systemConfig["DefaultLimitRTPRIO"] = ""
	systemConfig["DefaultLimitRTTIME"] = ""
	systemConfig["IPAddressAllow"] = ""
	systemConfig["IPAddressDeny"] = ""
}
