package antminer

import (
	"github.com/ka2n/masminer/machine"
)

// MinerSetting : /config/cgminer.conf
type MinerSetting struct {
	Options map[string]string
	Pools   []PoolSetting `json:"pools"`
}

// PoolSetting : A part of /config/cgminer.conf
type PoolSetting struct {
	URL  string `json:"url"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func (m *MinerSetting) LoadCommonMinerSetting(s machine.MinerSetting) error {
	var mv MinerSetting
	opt := make(map[string]string)
	opt["api-allow"] = defaultAPIAllow
	opt["api-groups"] = defaultAPIGroups
	opt["api-listen"] = defaultAPIListen
	opt["api-network"] = defaultAPINetwork
	opt["bitmain-use-vil"] = defaultBitmainUseVil

	for k, v := range s.Options {
		opt[k] = v
	}
	mv.Options = opt

	pools := make([]PoolSetting, len(s.Pools))
	for i, ip := range s.Pools {
		var p PoolSetting
		p.URL = ip.URL
		p.User = ip.User
		p.Pass = ip.Pass
		pools[i] = p
	}
	mv.Pools = pools
	*m = mv
	return nil
}

func (m MinerSetting) CommonMinerSetting() machine.MinerSetting {
	var s machine.MinerSetting

	s.Options = make(map[string]string)
	for k, v := range m.Options {
		s.Options[k] = v
	}

	s.Pools = make([]machine.Pool, len(m.Pools))
	for i, p := range m.Pools {
		var pool machine.Pool
		pool.URL = p.URL
		pool.User = p.User
		pool.Pass = p.Pass
		s.Pools[i] = pool
	}
	return s
}
