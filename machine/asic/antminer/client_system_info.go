package antminer

import (
	"bytes"
	"sync"

	"github.com/ka2n/masminer/machine"
	mnet "github.com/ka2n/masminer/net"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

// GetSystemInfo returns SystemInfo
func (c *Client) GetSystemInfo() (info SystemInfo, err error) {
	// Read from cache
	c.mu.RLock()
	if c.systemInfo != nil {
		c.mu.RUnlock()
		return *c.systemInfo, nil
	}
	c.mu.RUnlock()

	info, err = c.getSystemInfo()
	if err != nil {
		return info, err
	}

	// Cache
	c.mu.Lock()
	defer c.mu.Unlock()
	c.systemInfo = &info
	return info, nil
}

func (c *Client) getSystemInfo() (info SystemInfo, err error) {
	var wg errgroup.Group
	var mu sync.Mutex
	client := c.ssh
	info.MinerType = c.minerType

	wg.Go(func() error {
		ret, err := getMacAddr(client, c.ipCMDPath)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.MACAddr = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getIPAddr(client, c.ipCMDPath)
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
		ret, err := getModel(client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.Model = ret
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
		ret, err := getMinerVersion(client, c.versionCMD)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.MinerVersion = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getHardwareVersions(client, c.statsCMD)
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

func getMacAddr(client *ssh.Client, ipCMD string) (string, error) {
	cmd := ipCMD + ` link show eth0 | grep -o 'link/.*' | cut -d' ' -f2`
	ret, err := outputRemoteShell(client, cmd)
	return string(bytes.TrimSpace(ret)), err
}

func getHostname(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `hostname`)
	return string(ret), err
}

func getModel(client *ssh.Client) (machine.Model, error) {
	cmd := `sed -n 2p ` + metadataPath
	ret, err := outputRemoteShell(client, cmd)
	if err != nil {
		return machine.ModelUnknown, err
	}
	return MinerTypeFromString(string(ret))
}

func getKernelVersion(client *ssh.Client) (string, error) {
	cmd := `uname -srv`
	ret, err := outputRemoteShell(client, cmd)
	return string(bytes.TrimSpace(ret)), err
}

func getFileSystemVersion(client *ssh.Client) (string, error) {
	cmd := `sed -n 1p ` + metadataPath
	ret, err := outputRemoteShell(client, cmd)
	return string(bytes.TrimSpace(ret)), err
}

func getMinerVersion(client *ssh.Client, cmd string) (string, error) {
	ret, err := outputRemoteShell(client, cmd)
	if err != nil {
		return "", err
	}
	return parseCGMinerVersion(ret)
}

func getHardwareVersions(client *ssh.Client, cmd string) ([]string, error) {
	ret, err := outputRemoteShell(client, cmd)
	if err != nil {
		return nil, err
	}
	return parseHWVersionsFromCGMinerStats(bytes.TrimSpace(ret))
}

func getIPAddr(client *ssh.Client, ipCMD string) (string, error) {
	cmd := ipCMD + ` addr show eth0 | grep -o 'inet\s.*' | cut -d' ' -f2`
	ret, err := outputRemoteShell(client, cmd)
	if err != nil {
		return string(ret), err
	}
	return mnet.ParseIPAddr(string(bytes.TrimSpace(ret)))
}
