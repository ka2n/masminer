package antminer

import "github.com/ka2n/masminer/machine"

// SystemInfo : Generic system info
type SystemInfo struct {
	MACAddr           string
	Hostname          string
	ProductType       machine.MinerType
	SystemMode        string
	KernelVersion     string
	FileSystemVersion string
	CGMinerVersion    string
	HardwareVersions  []string
}
