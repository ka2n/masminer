package antminer

// MinerStats : miner status
type MinerStats struct {
	Summary MinerStatsSummary
	Pools   []MinerStatsPool
	Devs    MinerStatsDevs
}

// MinerStatsSummary : summary
type MinerStatsSummary struct {
	Elapsed            string
	GHS5s              string
	GHSAvarage         string
	Foundblocks        string
	Getworks           string
	Accepted           string
	Rejected           string
	HardwareErrors     string
	Utility            string
	Discarded          string
	Stale              string
	Localwork          string
	WorkUtility        string
	DifficultyAccepted string
	DifficultyRejected string
	DifficultyStale    string
	Bestshare          string
}

// MinerStatsPool : pool status
type MinerStatsPool struct {
	Index               string
	URL                 string
	User                string
	Status              string
	StratumActive       string
	Priority            string
	Getworks            string
	Accepted            string
	Rejected            string
	Discarded           string
	Stale               string
	Diff                string
	Diff1Shares         string
	DifficultyAccepted  string
	DifficultyRejected  string
	DifficultyStale     string
	LastShareDifficulty string
	LastShareTime       string
}

// MinerStatsDevs : miner status
type MinerStatsDevs struct {
	Fans   []string
	Chains []MinerStatsChain
}

// MinerStatsChain : miner status per chain
type MinerStatsChain struct {
	Index    string
	Acn      string
	Freq     string
	Rate     string
	Hw       string
	TempPCB  string
	TempChip string
	Status   string
}
