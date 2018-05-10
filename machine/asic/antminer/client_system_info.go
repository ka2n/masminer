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
		ret, err := getProductType(client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.ProductType = ret
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
		ret, err := getCGMinerVersion(client)
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		info.CGMinerVersion = ret
		return nil
	})

	wg.Go(func() error {
		ret, err := getHardwareVersions(client)
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

func getMacAddr(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `ip link show eth0 | grep -o 'link/.*' | cut -d' ' -f2`)
	return string(bytes.TrimSpace(ret)), err
}

func getHostname(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `hostname`)
	return string(ret), err
}

func getProductType(client *ssh.Client) (machine.MinerType, error) {
	ret, err := outputRemoteShell(client, `sed -n 2p `+metadataPath)
	if err != nil {
		return machine.MinerTypeUnknown, err
	}
	return MinerTypeFromString(string(ret))
}

func getKernelVersion(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `uname -srv`)
	return string(bytes.TrimSpace(ret)), err
}

func getFileSystemVersion(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `sed -n 1p `+metadataPath)
	return string(bytes.TrimSpace(ret)), err
}

func getCGMinerVersion(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, minerAPIVersionCMD)
	return parseCGMinerVersion(ret), err
}

func getHardwareVersions(client *ssh.Client) ([]string, error) {
	ret, err := outputRemoteShell(client, minerAPIStatsCMD)
	if err != nil {
		return nil, err
	}
	return parseHWVersionsFromCGMinerStats(bytes.TrimSpace(ret))
}

func getIPAddr(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `ip a show eth0 | grep -o 'inet.*' | cut -d' ' -f2`)
	if err != nil {
		return string(ret), err
	}
	return mnet.ParseIPAddr(string(bytes.TrimSpace(ret)))
}
