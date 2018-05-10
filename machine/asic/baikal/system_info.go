package baikal

import (
	"github.com/ka2n/masminer/machine"
)

type SystemInfo struct {
	IPAddr        string
	MACAddr       string
	Hostname      string
	KernelVersion string

	ProductType    machine.MinerType
	ProductVersion string

	FileSystemVersion string

	MinerDescription string
	MinerVersion     string
	APIVersion       string
}
