package dayun

import (
	"github.com/ka2n/masminer/machine"
)

// SystemInfo : Generic system info
type SystemInfo struct {
	IPAddr           string
	MACAddr          string
	Hostname         string
	KernelVersion    string
	DashboardVersion string

	UptimeSeconds string

	ProductType      machine.Model
	MinerDescription string
	MinerVersion     string
	APIVersion       string
}
