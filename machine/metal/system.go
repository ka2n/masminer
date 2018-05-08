package metal

import (
	"os"
	"syscall"

	"github.com/shirou/gopsutil/cpu"
	"golang.org/x/sync/singleflight"
)

var (
	mu singleflight.Group
)

// ServerInfo : サーバの情報を取得します
func ServerInfo() (ip, mac, hostname string, err error) {
	type info struct {
		ip       string
		hostname string
		mac      string
	}
	ret, err, _ := mu.Do("serverinfo", func() (interface{}, error) {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, err
		}

		mac, ip = ActiveAddress()
		return info{ip: ip, hostname: hostname, mac: mac}, nil
	})
	if err != nil {
		mu.Forget("serverinfo")
		return
	}

	rr := ret.(info)
	return rr.ip, rr.mac, rr.hostname, nil
}

// SysInfo : sysinfoとCPUの情報を取得します
func SysInfo() (sys *syscall.Sysinfo_t, cpus []cpu.InfoStat, err error) {
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		return nil, nil, err
	}

	cpui, err := cpu.Info()
	if err != nil {
		return nil, nil, err
	}

	return &info, cpui, nil
}
