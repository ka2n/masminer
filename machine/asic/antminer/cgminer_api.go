package antminer

import (
	"bytes"
	"strconv"
)

func parseCGMinerStats(in []byte) []map[string]string {
	in = bytes.TrimSpace(in)
	segs := bytes.Split(in, []byte{'|'})
	lprops := make([]map[string]string, 0, len(segs))

	for _, seg := range segs {
		if len(seg) == 0 {
			continue
		}
		rows := bytes.Split(seg, []byte{','})
		if len(rows) == 0 {
			continue
		}
		props := make(map[string]string, len(rows))
		for _, row := range rows {
			kv := bytes.SplitN(row, []byte{'='}, 2)
			if len(kv) == 2 {
				props[string(kv[0])] = string(kv[1])
			}
		}
		lprops = append(lprops, props)
	}

	return lprops
}

func parseCGMinerVersion(in []byte) string {
	lprops := parseCGMinerStats(bytes.TrimSpace(in))
	props := lprops[1]
	return props["CGMiner"]
}

func parseHWVersionsFromCGMinerStats(in []byte) ([]string, error) {
	lprops := parseCGMinerStats(in)
	props := lprops[2]
	if _, ok := props["hwv1"]; ok {
		vs := []string{}
		i := 1
		for {
			key := "hwv" + strconv.Itoa(i)
			v, ok := props[key]
			if !ok {
				break
			}
			vs = append(vs, v)
			i++
		}
		return vs, nil
	}
	return []string{lprops[1]["Miner"]}, nil
}

func parseDevsFromCGMinerStats(in []byte) (MinerStatsDevs, error) {
	var devs MinerStatsDevs
	lprops := parseCGMinerStats(in)
	props := lprops[2]

	devs.Fans = []string{}
	i := 1
	for {
		fan, ok := props["fan"+strconv.Itoa(i)]
		if !ok {
			break
		}
		devs.Fans = append(devs.Fans, fan)
		i++
	}

	freq := props["frequency"]

	devs.Chains = []MinerStatsChain{}
	i = 1
	for {
		var chain MinerStatsChain
		var ok bool
		is := strconv.Itoa(i)

		chain.Index = is

		chain.Freq = freq

		chain.TempPCB, ok = props["temp"+is]
		if !ok {
			break
		}
		chain.TempChip, ok = props["temp2_"+is]
		if !ok {
			break
		}
		chain.Acn, ok = props["chain_acn"+is]
		if !ok {
			break
		}
		chain.Hw, ok = props["chain_hw"+is]
		if !ok {
			break
		}
		chain.Rate, ok = props["chain_rate"+is]
		if !ok {
			break
		}
		chain.Status, ok = props["chain_acs"+is]
		if !ok {
			break
		}

		devs.Chains = append(devs.Chains, chain)
		i++
	}

	return devs, nil
}

func parsePoolsFromCGMinerPools(in []byte) ([]MinerStatsPool, error) {
	lprops := parseCGMinerStats(in)
	pools := make([]MinerStatsPool, len(lprops[1:]))
	for i, props := range lprops[1:] {
		var pool MinerStatsPool
		pool.Index = props["POOL"]
		pool.URL = props["URL"]
		pool.User = props["User"]
		pool.Status = props["Status"]
		pool.StratumActive = props["Stratum Active"]
		pool.Priority = props["Priority"]
		pool.Getworks = props["Getworks"]
		pool.Accepted = props["Accepted"]
		pool.Rejected = props["Rejected"]
		pool.Discarded = props["Discarded"]
		pool.Stale = props["Stale"]
		pool.Diff = props["Diff"]
		pool.Diff1Shares = props["Diff1 Shares"]
		pool.DifficultyAccepted = props["Difficulty Accepted"]
		pool.DifficultyRejected = props["Difficulty Rejected"]
		pool.DifficultyStale = props["Difficulty Stale"]
		pool.LastShareDifficulty = props["Last Share Difficulty"]
		pool.LastShareTime = props["Last Share Time"]
		pools[i] = pool
	}
	return pools, nil
}
