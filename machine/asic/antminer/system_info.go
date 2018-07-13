package antminer

import "github.com/ka2n/masminer/machine"

// SystemInfo : Generic system info
type SystemInfo struct {
	IPAddr            string
	MACAddr           string
	Hostname          string
	KernelVersion     string
	FileSystemVersion string

	BootTime string

	Model            machine.Model
	MinerType        string
	MinerVersion     string
	HardwareVersions []string
}
