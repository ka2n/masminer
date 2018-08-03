package baikal

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetStats() (stat MinerStats, err error) {
	return c.GetStatsContext(nil)
}

func (c *Client) GetStatsContext(ctx context.Context) (stat MinerStats, err error) {
	ret, err := outputMinerRPC(ctx, c.ssh, "devs+pools+stats+summary", "")
	if err != nil {
		return stat, err
	}

	var resp struct {
		SGMultipleCMDResponse
		Devs    []SGDevsResponse    `json:"devs"`
		Pools   []SGPoolsResponse   `json:"pools"`
		Stats   []SGStatsResponse   `json:"stats"`
		Summary []SGSummaryResponse `json:"summary"`
	}
	err = json.Unmarshal(ret, &resp)
	if err != nil {
		return stat, err
	}

	if len(resp.Devs) == 0 {
		return stat, fmt.Errorf("does not return DEVS")
	}

	if len(resp.Summary) == 0 || len(resp.Summary[0].Summary) == 0 {
		return stat, fmt.Errorf("does not return SUMMARY")
	}

	if len(resp.Stats) == 0 {
		return stat, fmt.Errorf("does not return STATS")
	}

	if len(resp.Pools) == 0 {
		return stat, fmt.Errorf("does not return POOLS")
	}

	stat.Summary = resp.Summary[0].Summary[0]
	stat.Devs = resp.Devs[0].Devs
	stat.Stats = resp.Stats[0].Stats
	stat.Pools = resp.Pools[0].Pools

	return stat, nil
}
