package baikal

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ka2n/masminer/minerapi"
	"golang.org/x/sync/errgroup"

	"github.com/ka2n/masminer/machine/asic/base"
)

func (c *Client) GetStats() (stat MinerStats, err error) {
	return c.GetStatsContext(context.Background())
}

func (c *Client) GetStatsContext(ctx context.Context) (stat MinerStats, err error) {
	var mu sync.Mutex
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		ret, err := base.GetCPUTemp(ctx, c.SSH)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		stat.System.TempCPU = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := base.OutputMinerRPC(ctx, c.SSH, "devs+pools+stats+summary", "")
		if err != nil {
			return err
		}

		var resp struct {
			minerapi.MultipleResponse
			Devs    []SGDevsResponse    `json:"devs"`
			Pools   []SGPoolsResponse   `json:"pools"`
			Stats   []SGStatsResponse   `json:"stats"`
			Summary []SGSummaryResponse `json:"summary"`
		}
		err = json.Unmarshal(ret, &resp)
		if err != nil {
			return err
		}

		if len(resp.Devs) == 0 {
			return fmt.Errorf("does not return DEVS")
		}

		if len(resp.Summary) == 0 || len(resp.Summary[0].Summary) == 0 {
			return fmt.Errorf("does not return SUMMARY")
		}

		if len(resp.Stats) == 0 {
			return fmt.Errorf("does not return STATS")
		}

		if len(resp.Pools) == 0 {
			return fmt.Errorf("does not return POOLS")
		}

		mu.Lock()
		defer mu.Unlock()
		stat.Summary = resp.Summary[0].Summary[0]
		stat.Devs = resp.Devs[0].Devs
		stat.Stats = resp.Stats[0].Stats
		stat.Pools = resp.Pools[0].Pools
		return nil
	})

	return stat, wg.Wait()
}
