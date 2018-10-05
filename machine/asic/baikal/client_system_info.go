package baikal

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ka2n/masminer/machine/asic/base"
	"github.com/ka2n/masminer/minerapi"

	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

func (c *Client) GetSystemInfo() (info SystemInfo, err error) {
	return c.GetSystemInfoContext(nil)
}

func (c *Client) GetSystemInfoContext(ctx context.Context) (info SystemInfo, err error) {
	// Read from cache
	c.MU.RLock()
	if c.systemInfo != nil {
		c.MU.RUnlock()
		return *c.systemInfo, nil
	}
	c.MU.RUnlock()

	info, err = getSystemInfo(ctx, c.SSH)
	if err != nil {
		return info, err
	}

	// Cache
	c.MU.Lock()
	defer c.MU.Unlock()
	c.systemInfo = &info
	return info, nil
}

func getSystemInfo(ctx context.Context, client *ssh.Client) (info SystemInfo, err error) {
	var wg errgroup.Group
	var mu sync.Mutex

	wg.Go(func() error {
		ret, err := base.GetMacAddr(ctx, client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.MACAddr = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := base.GetIPAddr(ctx, client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.IPAddr = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := base.GetHostname(ctx, client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.Hostname = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := base.GetKernelVersion(ctx, client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.KernelVersion = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := base.GetFileSystemVersion(ctx, client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.FileSystemVersion = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := base.GetUptimeSeconds(ctx, client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.UptimeSeconds = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := base.OutputMinerRPC(ctx, client, "stats+version", "")
		if err != nil {
			return err
		}

		var resp struct {
			minerapi.MultipleResponse
			Version []SGVersionResponse `json:"version"`
			Stats   []SGStatsResponse   `json:"stats"`
		}
		err = json.Unmarshal(ret, &resp)
		if err != nil {
			return err
		}

		if !(len(resp.Version) == 1 && len(resp.Version[0].Version) == 1) {
			return fmt.Errorf("error sgminer RPC response")
		}

		mu.Lock()
		version := resp.Version[0].Version[0]
		info.MinerDescription = version.Miner
		info.MinerVersion = version.SGMiner
		info.APIVersion = version.API
		mu.Unlock()

		if !(len(resp.Stats) != 0 && len(resp.Stats[0].Stats) != 0) {
			return fmt.Errorf("error sgminer RPC response")
		}

		stat := resp.Stats[0].Stats[0]

		mu.Lock()
		defer mu.Unlock()

		info.ProductType, err = modelFromAPIHWV(stat.HWV.String())
		if err != nil {
			return err
		}
		info.ProductVersion, err = minerVersionFromFWV(stat.FWV.String())
		if err != nil {
			return err
		}
		return nil
	})

	return info, wg.Wait()
}
