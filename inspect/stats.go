package inspect

// RigStat is a status of machine
type RigStat struct {
	GHS5s      string
	GHSAvarage string
	MHS5s      string
	MHSAvarage string
	KHS5s      string
	KHSAvarage string

	Accepted       string
	Rejected       string
	HardwareErrors string
	Utility        string

	Devices []DeviceStat
	Pools   []PoolStat
}

// DeviceStat is a status of each chip
type DeviceStat struct {
	TempPCB        string
	TempChip       string
	Frequency      string
	Chips          int
	HardwareErrors string
	Hashrate       string
}

// PoolStat is a status of pool connection
type PoolStat struct {
	URL  string
	User string
	Algo string

	Status              string
	StratumActive       bool
	Priority            int
	Getworks            string
	Accepted            string
	Rejected            string
	Discarded           string
	Stale               string
	DifficultyAccepted  string
	DifficultyRejected  string
	DifficultyStale     string
	LastShareDifficulty string
	LastShareTime       string
}
