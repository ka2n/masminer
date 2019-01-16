package baikal

import (
	"github.com/ka2n/masminer/machine"
)

// SystemInfo : Generic system info
type SystemInfo struct {
	IPAddr        string
	MACAddr       string
	Hostname      string
	KernelVersion string

	UptimeSeconds string

	ProductType      machine.Model
	ProductVersion   string
	MinerDescription string
	MinerVersion     string
	APIVersion       string
}
