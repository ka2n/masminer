package machine

// MinerSetting is a miner setting
type MinerSetting struct {
	Pools   []Pool
	Options map[string]string
}
