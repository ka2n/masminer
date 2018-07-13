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
	Model           Model
	HardwareVersion string
	FirmwareVersion string
	MinerType       string
	MinerVersion    string
	Algos           []string
	BootTime        string
}
