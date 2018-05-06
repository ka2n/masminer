package baikal

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ka2n/masminer/inspect"
	"golang.org/x/crypto/ssh"
)

// MinerSetting : parsed /opt/scripta/etc/miner.options.json
type MinerSetting map[string]string

type minerSettingRaw []minerSettingRawRow

type minerSettingRawRow struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (m *MinerSetting) UnmarshalJSON(in []byte) error {
	var r minerSettingRaw
	mv := make(MinerSetting)

	if err := json.Unmarshal(in, &r); err != nil {
		return err
	}

	for _, kv := range r {
		mv[kv.Key] = kv.Value
	}

	*m = mv
	return nil
}

func (m MinerSetting) MarshalJSON() ([]byte, error) {
	r := make(minerSettingRaw, 0, len(m))
	for k, v := range m {
		r = append(r, minerSettingRawRow{k, v})
	}
	return json.Marshal(&r)
}

type PoolSetting struct {
	URL        string `json:"url"`
	Pass       string `json:"pass"`
	Priority   string `json:"priority"`
	Algo       string `json:"algo"`
	Extranonce string `json:"extranonce"`
	User       string `json:"user"`
}

// MinerOptions : /opt/scripta/etc/miner.options.json
type MinerOptions []MinerOptionsEntry

type MinerOptionsEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func getMinerSetting(client *ssh.Client) (setting MinerSetting, err error) {
	output, err := outputRemoteShell(client, `cat `+minerOptionsPath)
	if err != nil {
		return setting, err
	}

	if err := json.Unmarshal(output, &setting); err != nil {
		return setting, err
	}

	return setting, nil
}

func getMinerPools(client *ssh.Client) (pools []PoolSetting, err error) {
	output, err := outputRemoteShell(client, `cat `+minerPoolsPath)
	if err != nil {
		return pools, err
	}

	if err := json.Unmarshal(output, &pools); err != nil {
		return pools, err
	}

	return pools, nil
}

func mergeMinerAndPoolSetting(m MinerSetting, ps []PoolSetting) map[string]interface{} {
	minerConf := make(map[string]interface{})
	for k, v := range m {
		minerConf[k] = v
	}
	sort.Slice(ps, func(i, j int) bool {
		return ps[i].Priority < ps[j].Priority
	})
	minerConf["pools"] = ps
	return minerConf
}

func writeMinerAndPoolSetting(client *ssh.Client, m MinerSetting, ps []PoolSetting) error {
	mc := mergeMinerAndPoolSetting(m, ps)
	const indent = "    "

	mcb, err := json.MarshalIndent(mc, "", indent)
	if err != nil {
		return err
	}

	mb, err := json.MarshalIndent(m, "", indent)
	if err != nil {
		return err
	}

	pb, err := json.MarshalIndent(ps, "", indent)
	if err != nil {
		return err
	}

	err = runRemoteShell(client, fmt.Sprintf(`
set -ex
MINER_CONF_PATH=%s
MINER_POOLS_PATH=%s
MINER_OPTIONS_PATH=%s
cat <<'EOF' > /tmp/miner.conf
%s
EOF
cat <<'EOF' > /tmp/miner.pools.json
%s
EOF
cat <<'EOF' > /tmp/miner.options.json
%s
EOF
sudo cp /tmp/miner.conf $MINER_CONF_PATH
sudo cp /tmp/miner.pools.json $MINER_POOLS_PATH
sudo cp /tmp/miner.options.json $MINER_OPTIONS_PATH
sudo chown www-data $MINER_CONF_PATH
sudo chown www-data $MINER_POOLS_PATH
sudo chown www-data $MINER_OPTIONS_PATH
`,
		minerConfPath, minerPoolsPath, minerOptionsPath,
		string(mcb), string(pb), string(mb),
	))
	if err != nil {
		return err
	}

	_, err = outputMinerRPC(client, "restart", "")
	return err
}

func getCommonMinerSetting(m MinerSetting, ps []PoolSetting) inspect.MinerSetting {
	var s inspect.MinerSetting
	s.Options = m
	pools := make([]inspect.Pool, len(ps))
	for i, p := range ps {
		var pool inspect.Pool
		pool.URL = p.URL
		pool.User = p.User
		pool.Pass = p.Pass
		pool.Algo = p.Algo
		pool.Options = make(map[string]string)
		pool.Options["priority"] = p.Priority
		pool.Options["extranonce"] = p.Extranonce
		pools[i] = pool
	}
	s.Pools = pools
	return s
}

func loadCommonMinerSetting(s inspect.MinerSetting, m *MinerSetting, ps *[]PoolSetting) error {
	mv := make(MinerSetting)

	mv["api-port"] = defaultAPIPort
	mv["api-allow"] = defaultAPIAllow

	for k, v := range s.Options {
		mv[k] = v
	}
	*m = mv

	psv := make([]PoolSetting, len(s.Pools))
	for i, ip := range s.Pools {
		var p PoolSetting
		p.URL = ip.URL
		p.User = ip.User
		p.Pass = ip.Pass
		p.Algo = ip.Algo
		p.Priority = ip.Options["priority"]
		p.Extranonce = ip.Options["extranonce"]
		psv[i] = p
	}
	*ps = psv
	return nil
}
