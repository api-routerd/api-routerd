// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Info Json request
type Info struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

func routerGetProcNetDev(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetNetDev(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcVersion(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetVersion(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcPlatformInformation(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetPlatformInformation(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcVirtualization(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetVirtualization(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcUserStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetUserStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcTemperatureStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetTemperatureStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcNetStat(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol := vars["protocol"]

	switch r.Method {
	case "GET":
		err := GetNetStat(rw, protocol)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcPidNetStat(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol := vars["protocol"]
	pid := vars["pid"]

	switch r.Method {
	case "GET":
		err := GetNetStatPid(rw, protocol, pid)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcInterfaceStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetInterfaceStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func routerGetProcProtoCountersStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetProtoCountersStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcGetSwapMemoryStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetSwapMemoryStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcVirtualMemoryStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetVirtualMemoryStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcCPUInfo(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetCPUInfo(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcCPUTimeStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetCPUTimeStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcAvgStat(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetAvgStat(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func configureProcSysVM(rw http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)
	vm := VM{
		Property: vars["path"],
	}

	switch r.Method {
	case "GET":
		err = vm.GetVM(rw)
		break
	case "PUT":

		v := new(Info)

		err = json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		vm.Value = v.Value
		err = vm.SetVM(rw)
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func configureProcSysNet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := SysNet{Path: vars["path"], Property: vars["conf"], Link: vars["link"]}

	switch r.Method {
	case "GET":
		err := proc.GetSysNet(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	case "PUT":

		v := new(Info)
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		proc.Value = v.Value
		err = proc.SetSysNet(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func routerGetProcMisc(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetMisc(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcNetArp(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetNetArp(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcModules(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetModules(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetProcProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := vars["pid"]
	property := vars["property"]

	switch r.Method {
	case "GET":

		err := GetProcessInfo(rw, pid, property)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetPartitions(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetPartitions(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetIOCounters(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetIOCounters(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetDiskUsage(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetDiskUsage(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

// RegisterRouterProc register with mux
func RegisterRouterProc(router *mux.Router) {
	n := router.PathPrefix("/proc").Subrouter().StrictSlash(false)

	n.HandleFunc("/avgstat", routerGetProcAvgStat)
	n.HandleFunc("/cpuinfo", routerGetProcCPUInfo)
	n.HandleFunc("/cputimestat", routerGetProcCPUTimeStat)
	n.HandleFunc("/diskusage", routerGetDiskUsage)
	n.HandleFunc("/interface-stat", routerGetProcInterfaceStat)
	n.HandleFunc("/iocounters", routerGetIOCounters)
	n.HandleFunc("/misc", routerGetProcMisc)
	n.HandleFunc("/modules", routerGetProcModules)
	n.HandleFunc("/net/arp", routerGetProcNetArp)
	n.HandleFunc("/netdev", routerGetProcNetDev)
	n.HandleFunc("/netstat/{protocol}", routerGetProcNetStat)
	n.HandleFunc("/partitions", routerGetPartitions)
	n.HandleFunc("/platform", routerGetProcPlatformInformation)
	n.HandleFunc("/process/{pid}/{property}/", routerGetProcProcess)
	n.HandleFunc("/proto-counter-stat", routerGetProcProtoCountersStat)
	n.HandleFunc("/proto-pid-stat/{pid}/{protocol}", routerGetProcPidNetStat)
	n.HandleFunc("/swap-memory", routerGetProcGetSwapMemoryStat)
	n.HandleFunc("/sys/net/{path}/{link}/{conf}", configureProcSysNet)
	n.HandleFunc("/sys/vm/{path}", configureProcSysVM)
	n.HandleFunc("/temperaturestat", routerGetProcTemperatureStat)
	n.HandleFunc("/userstat", routerGetProcUserStat)
	n.HandleFunc("/version", routerGetProcVersion)
	n.HandleFunc("/virtual-memory", routerGetProcVirtualMemoryStat)
	n.HandleFunc("/virtualization", routerGetProcVirtualization)
}
