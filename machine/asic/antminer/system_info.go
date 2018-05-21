package antminer

import "github.com/ka2n/masminer/machine"

// SystemInfo : Generic system info
type SystemInfo struct {
	IPAddr            string
	MACAddr           string
	Hostname          string
	Model             machine.Model
	KernelVersion     string
	FileSystemVersion string
	MinerType         string
	MinerVersion      string
	HardwareVersions  []string
}
