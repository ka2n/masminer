package machine

// RemoteRig is a reference to remote rig
type RemoteRig struct {
	Name     string
	MACAddr  string
	IPAddr   string
	Hostname string
}

// RigInfo is a detail of Rig
type RigInfo struct {
	Rig             RemoteRig
	Manufacture     string
	MinerType       MinerType
	HardwareVersion string
	FirmwareVersion string
	MinerVersion    string
	Algos           []string
}
