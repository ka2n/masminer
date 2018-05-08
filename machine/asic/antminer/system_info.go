package antminer

import "github.com/ka2n/masminer/inspect"

// SystemInfo : Generic system info
type SystemInfo struct {
	MACAddr           string
	Hostname          string
	ProductType       inspect.MinerType
	SystemMode        string
	KernelVersion     string
	FileSystemVersion string
	CGMinerVersion    string
	HardwareVersions  []string
}
