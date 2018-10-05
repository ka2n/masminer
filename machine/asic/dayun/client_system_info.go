package dayun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ka2n/masminer/machine"
	"github.com/ka2n/masminer/machine/asic/base"

	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

func (c *Client) GetSystemInfoContext(ctx context.Context) (info SystemInfo, err error) {
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
		v, err := getDashboardVersion(ctx, client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.DashboardVersion = v
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
		ret, err := getProductType(ctx, client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.ProductType = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := base.OutputMinerRPC(ctx, client, "version", "")
		if err != nil {
			return err
		}

		var resp VersionResponse
		err = json.Unmarshal(ret, &resp)
		if err != nil {
			return err
		}

		if !(len(resp.Version) == 1 &&
			len(resp.ResponseCommon.Status) > 0) {
			return fmt.Errorf("error cgminer RPC response")
		}

		mu.Lock()
		version := resp.Version[0]
		info.MinerDescription = resp.ResponseCommon.Status[0].Description
		info.MinerVersion = version.CGMiner
		info.APIVersion = version.API
		mu.Unlock()
		return nil
	})

	return info, wg.Wait()
}

func getProductType(ctx context.Context, client *ssh.Client) (machine.Model, error) {
	cmd := `grep -q 'Zig Z1' /var/www/html/src/Template/Layout/default.twig`
	if err := base.RunRemoteShell(ctx, client, cmd); err != nil {
		return machine.ModelUnknown, nil
	}
	return ModelZ1, nil
}

func getDashboardVersion(ctx context.Context, client *ssh.Client) (string, error) {
	cmd := `cat /var/www/html/VERSION`
	out, err := base.OutputRemoteShell(ctx, client, cmd)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(out)), nil
}
