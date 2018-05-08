package baikal

import (
	"github.com/ka2n/masminer/inspect"
)

type SystemInfo struct {
	MACAddr       string
	Hostname      string
	KernelVersion string

	ProductType    inspect.MinerType
	ProductVersion string

	FileSystemVersion string

	MinerDescription string
	MinerVersion     string
	APIVersion       string
}
