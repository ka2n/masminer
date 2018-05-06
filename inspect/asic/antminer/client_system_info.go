package antminer

import (
	"bytes"

	"github.com/ka2n/masminer/inspect"
	"golang.org/x/crypto/ssh"
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
	info.MACAddr, err = getMacAddr(client)
	if err != nil {
		return info, err
	}
	info.Hostname, err = getHostname(client)
	if err != nil {
		return info, err
	}
	info.ProductType, err = getMinerType(client)
	if err != nil {
		return info, err
	}
	info.KernelVersion, err = getKernelVersion(client)
	if err != nil {
		return info, err
	}
	info.FileSystemVersion, err = getFileSystemVersion(client)
	if err != nil {
		return info, err
	}
	info.CGMinerVersion, err = getCGMinerVersion(client)
	if err != nil {
		return info, err
	}
	info.HardwareVersions, err = getHardwareVersions(client)
	if err != nil {
		return info, err
	}
	return info, nil
}

func getMacAddr(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `ip link show eth0 | grep -o 'link/.*' | cut -d' ' -f2`)
	return string(bytes.TrimSpace(ret)), err
}

func getHostname(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `hostname`)
	return string(ret), err
}

func getMinerType(client *ssh.Client) (inspect.MinerType, error) {
	ret, err := outputRemoteShell(client, `sed -n 2p `+metadataPath)
	if err != nil {
		return inspect.MinerTypeUnknown, err
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
