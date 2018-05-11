package baikal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	mnet "github.com/ka2n/masminer/net"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

func (c *Client) GetSystemInfo() (info SystemInfo, err error) {
	// Read from cache
	c.mu.RLock()
	if c.systemInfo != nil {
		c.mu.RUnlock()
		return *c.systemInfo, nil
	}
	c.mu.RUnlock()

	info, err = getSystemInfo(c.ssh)
	if err != nil {
		return info, err
	}

	// Cache
	c.mu.Lock()
	defer c.mu.Unlock()
	c.systemInfo = &info
	return info, nil
}

func getSystemInfo(client *ssh.Client) (info SystemInfo, err error) {
	var wg errgroup.Group
	var mu sync.Mutex

	wg.Go(func() error {
		ret, err := getMacAddr(client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.MACAddr = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getIPAddr(client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.IPAddr = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getHostname(client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.Hostname = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getKernelVersion(client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.KernelVersion = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getFileSystemVersion(client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.FileSystemVersion = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := outputMinerRPC(client, "stats+version", "")
		if err != nil {
			return err
		}

		var resp struct {
			SGMultipleCMDResponse
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

		info.ProductType, err = minerTypeFromAPIHWV(stat.HWV.String())
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

func getMacAddr(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `ip link show eth0 | grep -o 'link/.*' | cut -d' ' -f2`)
	return string(bytes.TrimSpace(ret)), err
}

func getHostname(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `hostname`)
	return string(ret), err
}

func getKernelVersion(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `uname -srv`)
	return string(bytes.TrimSpace(ret)), err
}

func getFileSystemVersion(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `uname -v`)
	return string(bytes.TrimSpace(ret)), err
}

func getIPAddr(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `ip a show eth0 | grep -o 'inet.*' | cut -d' ' -f2`)
	if err != nil {
		return string(ret), err
	}
	return mnet.ParseIPAddr(string(bytes.TrimSpace(ret)))
}
