package antminer

import (
	"bytes"
	"context"
	"sync"

	"github.com/ka2n/masminer/machine"
	"github.com/ka2n/masminer/machine/asic/base"
	mnet "github.com/ka2n/masminer/net"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

// GetSystemInfo returns SystemInfo
func (c *Client) GetSystemInfo() (info SystemInfo, err error) {
	return c.GetSystemInfoContext(context.Background())
}

// GetSystemInfoContext returns SystemInfo
func (c *Client) GetSystemInfoContext(ctx context.Context) (info SystemInfo, err error) {
	// Read from cache
	c.MU.RLock()
	if c.systemInfo != nil {
		c.MU.RUnlock()
		return *c.systemInfo, nil
	}
	c.MU.RUnlock()

	info, err = c.getSystemInfo(ctx)
	if err != nil {
		return info, err
	}

	// Cache
	c.MU.Lock()
	defer c.MU.Unlock()
	c.systemInfo = &info
	return info, nil
}

func (c *Client) getSystemInfo(ctx context.Context) (info SystemInfo, err error) {
	var mu sync.Mutex
	client := c.SSH
	info.MinerType = c.minerType

	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		ret, err := getMacAddr(ctx, client, c.ipCMDPath)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.MACAddr = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getIPAddr(ctx, client, c.ipCMDPath)
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
		ret, err := getModel(ctx, client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.Model = ret
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
		ret, err := getFileSystemVersion(ctx, client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.FileSystemVersion = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getMinerVersion(ctx, client, c.versionCMD)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.MinerVersion = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getHardwareVersions(ctx, client, c.statsCMD)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.HardwareVersions = ret
		return nil
	})
	return info, wg.Wait()
}

func getMacAddr(ctx context.Context, client *ssh.Client, ipCMD string) (string, error) {
	cmd := ipCMD + ` link show eth0 | grep -o 'link/.*' | cut -d' ' -f2`
	ret, err := base.OutputRemoteShell(ctx, client, cmd)
	return string(bytes.TrimSpace(ret)), err
}

func getModel(ctx context.Context, client *ssh.Client) (machine.Model, error) {
	cmd := `sed -n 2p ` + metadataPath
	ret, err := base.OutputRemoteShell(ctx, client, cmd)
	if err != nil {
		return machine.ModelUnknown, err
	}
	return MinerTypeFromString(string(ret))
}

func getFileSystemVersion(ctx context.Context, client *ssh.Client) (string, error) {
	cmd := `sed -n 1p ` + metadataPath
	ret, _ := base.OutputRemoteShell(ctx, client, cmd)
	// let ignore error. some model have FileSystemVersion in it's binary, support it someday.
	return string(bytes.TrimSpace(ret)), nil
}

func getMinerVersion(ctx context.Context, client *ssh.Client, cmd string) (string, error) {
	ret, err := base.OutputRemoteShell(ctx, client, cmd)
	if err != nil {
		return "", err
	}
	return parseCGMinerVersion(ret)
}

func getHardwareVersions(ctx context.Context, client *ssh.Client, cmd string) ([]string, error) {
	ret, err := base.OutputRemoteShell(ctx, client, cmd)
	if err != nil {
		return nil, err
	}
	return parseHWVersionsFromCGMinerStats(bytes.TrimSpace(ret))
}

func getIPAddr(ctx context.Context, client *ssh.Client, ipCMD string) (string, error) {
	cmd := ipCMD + ` addr show eth0 | grep -o 'inet\s.*' | cut -d' ' -f2`
	ret, err := base.OutputRemoteShell(ctx, client, cmd)
	if err != nil {
		return string(ret), err
	}
	return mnet.ParseIPAddr(string(bytes.TrimSpace(ret)))
}
