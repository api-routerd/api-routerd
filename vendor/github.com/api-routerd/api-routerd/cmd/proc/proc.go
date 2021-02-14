// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/api-routerd/api-routerd/cmd/share"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

const (
	procMiscPath    = "/proc/misc"
	procNetArpPath  = "/proc/net/arp"
	procModulesPath = "/proc/modules"
)

// NetARP Json response
type NetARP struct {
	IPAddress string `json:"ip_address"`
	HWType    string `json:"hw_type"`
	Flags     string `json:"flags"`
	HWAddress string `json:"hw_address"`
	Mask      string `json:"mask"`
	Device    string `json:"device"`
}

// Modules Json response
type Modules struct {
	Module     string `json:"module"`
	MemorySize string `json:"memory_size"`
	Instances  string `json:"instances"`
	Dependent  string `json:"dependent"`
	State      string `json:"state"`
}

// GetVersion get system version
func GetVersion(rw http.ResponseWriter) error {
	infostat, err := host.Info()
	if err != nil {
		return err
	}

	return share.JSONResponse(infostat, rw)
}

// GetPlatformInformation read platform info
func GetPlatformInformation(rw http.ResponseWriter) error {
	platform, family, version, err := host.PlatformInformation()
	if err != nil {
		return err
	}

	p := struct {
		Platform string
		Family   string
		Version  string
	}{
		platform,
		family,
		version,
	}

	return share.JSONResponse(p, rw)
}

// GetVirtualization read virt info
func GetVirtualization(rw http.ResponseWriter) error {
	system, role, err := host.Virtualization()
	if err != nil {
		return err
	}

	v := struct {
		System string
		Role   string
	}{
		system,
		role,
	}

	return share.JSONResponse(v, rw)
}

// GetUserStat active users
func GetUserStat(rw http.ResponseWriter) error {
	userstat, err := host.Users()
	if err != nil {
		return err
	}

	return share.JSONResponse(userstat, rw)
}

// GetTemperatureStat read temp of HW
func GetTemperatureStat(rw http.ResponseWriter) error {
	tempstat, err := host.SensorsTemperatures()
	if err != nil {
		return err
	}

	return share.JSONResponse(tempstat, rw)
}

// GetNetStat read netstat from proc tcp/udp/sctp
func GetNetStat(rw http.ResponseWriter, protocol string) error {
	conn, err := net.Connections(protocol)
	if err != nil {
		return err
	}

	return share.JSONResponse(conn, rw)
}

// GetNetStatPid connection by the pid
func GetNetStatPid(rw http.ResponseWriter, protocol string, process string) error {
	pid, err := strconv.ParseInt(process, 10, 32)
	if err != nil || protocol == "" || pid == 0 {
		return errors.New("Can't parse request")
	}

	conn, err := net.ConnectionsPid(protocol, int32(pid))
	if err != nil {
		return err
	}

	return share.JSONResponse(conn, rw)
}

// GetProtoCountersStat protocol specific statitics
func GetProtoCountersStat(rw http.ResponseWriter) error {
	protocols := []string{"ip", "icmp", "icmpmsg", "tcp", "udp", "udplite"}

	proto, err := net.ProtoCounters(protocols)
	if err != nil {
		return err
	}

	return share.JSONResponse(proto, rw)
}

// GetNetDev network device stat
func GetNetDev(rw http.ResponseWriter) error {
	netdev, err := net.IOCounters(true)
	if err != nil {
		return err
	}

	return share.JSONResponse(netdev, rw)
}

// GetInterfaceStat network interface stat
func GetInterfaceStat(rw http.ResponseWriter) error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	return share.JSONResponse(interfaces, rw)
}

// GetSwapMemoryStat swap memory info
func GetSwapMemoryStat(rw http.ResponseWriter) error {
	swap, err := mem.SwapMemory()
	if err != nil {
		return err
	}

	return share.JSONResponse(swap, rw)
}

// GetVirtualMemoryStat VM information
func GetVirtualMemoryStat(rw http.ResponseWriter) error {
	virt, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	return share.JSONResponse(virt, rw)
}

// GetCPUInfo CPU info
func GetCPUInfo(rw http.ResponseWriter) error {
	cpus, err := cpu.Info()
	if err != nil {
		return err
	}

	return share.JSONResponse(cpus, rw)
}

// GetCPUTimeStat read CPU time
func GetCPUTimeStat(rw http.ResponseWriter) error {
	cpus, err := cpu.Times(true)
	if err != nil {
		return err
	}

	return share.JSONResponse(cpus, rw)
}

// GetAvgStat read avg stat
func GetAvgStat(rw http.ResponseWriter) error {
	avgstat, r := load.Avg()
	if r != nil {
		return r
	}

	return share.JSONResponse(avgstat, rw)
}

// GetPartitions read all partitions
func GetPartitions(rw http.ResponseWriter) error {
	p, r := disk.Partitions(true)
	if r != nil {
		return r
	}

	return share.JSONResponse(p, rw)
}

// GetIOCounters read IO counters
func GetIOCounters(rw http.ResponseWriter) error {
	i, r := disk.IOCounters()
	if r != nil {
		return r
	}

	return share.JSONResponse(i, rw)
}

// GetDiskUsage read disk usage
func GetDiskUsage(rw http.ResponseWriter) error {
	u, r := disk.Usage("/")
	if r != nil {
		return r
	}

	return share.JSONResponse(u, rw)
}

// GetMisc read /proc/misc
func GetMisc(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(procMiscPath)
	if err != nil {
		log.Fatalf("Failed to read: %s", procMiscPath)
		return errors.New("Failed to read misc")
	}

	miscMap := make(map[int]string)
	for _, line := range lines {
		fields := strings.Fields(line)

		deviceNum, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		miscMap[deviceNum] = fields[1]
	}

	return share.JSONResponse(miscMap, rw)
}

// GetNetArp get ARP info
func GetNetArp(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(procNetArpPath)
	if err != nil {
		log.Fatalf("Failed to read: %s", procNetArpPath)
		return errors.New("Failed to read /proc/net/arp")
	}

	netarp := make([]NetARP, len(lines)-1)
	for i, line := range lines {
		if i == 0 {
			continue
		}

		fields := strings.Fields(line)

		arp := NetARP{}
		arp.IPAddress = fields[0]
		arp.HWType = fields[1]
		arp.Flags = fields[2]
		arp.HWAddress = fields[3]
		arp.Mask = fields[4]
		arp.Device = fields[5]
		netarp[i-1] = arp
	}

	return share.JSONResponse(netarp, rw)
}

// GetModules Get all installed modules
func GetModules(rw http.ResponseWriter) error {
	lines, err := share.ReadFullFile(procModulesPath)
	if err != nil {
		log.Fatalf("Failed to read: %s", procModulesPath)
		return errors.New("Failed to read /proc/modules")
	}

	modules := make([]Modules, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)

		module := Modules{}

		for j, field := range fields {
			switch j {
			case 0:
				module.Module = field
				break
			case 1:
				module.MemorySize = field
				break
			case 2:
				module.Instances = field
				break
			case 3:
				module.Dependent = field
				break
			case 4:
				module.State = field
				break
			}
		}

		modules[i] = module
	}

	return share.JSONResponse(modules, rw)
}

// GetProcessInfo get process information from proc
func GetProcessInfo(rw http.ResponseWriter, proc string, property string) error {
	pid, err := strconv.ParseInt(proc, 10, 32)
	if err != nil {
		return err
	}

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return err
	}

	switch property {
	case "pid-connections":
		conn, err := p.Connections()
		if err != nil {
			return err
		}

		return share.JSONResponse(conn, rw)

	case "pid-rlimit":
		rlimit, err := p.Rlimit()
		if err != nil {
			return err
		}

		return share.JSONResponse(rlimit, rw)

	case "pid-rlimit-usage":
		rlimit, err := p.RlimitUsage(true)
		if err != nil {
			return err
		}

		return share.JSONResponse(rlimit, rw)

	case "pid-status":
		s, err := p.Status()
		if err != nil {
			return err
		}

		return share.JSONResponse(s, rw)

	case "pid-username":
		u, err := p.Username()
		if err != nil {
			return err
		}

		return share.JSONResponse(u, rw)

	case "pid-open-files":
		f, err := p.OpenFiles()
		if err != nil {
			return err
		}

		return share.JSONResponse(f, rw)

	case "pid-fds":
		f, err := p.NumFDs()
		if err != nil {
			return err
		}

		return share.JSONResponse(f, rw)

	case "pid-name":
		n, err := p.Name()
		if err != nil {
			return err
		}

		return share.JSONResponse(n, rw)

	case "pid-memory-percent":
		m, err := p.MemoryPercent()
		if err != nil {
			return err
		}

		return share.JSONResponse(m, rw)

	case "pid-memory-maps":
		m, err := p.MemoryMaps(true)
		if err != nil {
			return err
		}

		return share.JSONResponse(m, rw)

	case "pid-memory-info":
		m, err := p.MemoryInfo()
		if err != nil {
			return err
		}

		return share.JSONResponse(m, rw)

	case "pid-io-counters":
		m, err := p.IOCounters()
		if err != nil {
			return err
		}

		return share.JSONResponse(m, rw)
	}

	return nil
}
