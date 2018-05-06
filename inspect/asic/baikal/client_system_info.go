package baikal

import (
	"bytes"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/ssh"
)

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
	info.MACAddr, err = getMacAddr(c.ssh)
	if err != nil {
		return info, err
	}
	info.Hostname, err = getHostname(c.ssh)
	if err != nil {
		return info, err
	}
	info.KernelVersion, err = getKernelVersion(c.ssh)
	if err != nil {
		return info, err
	}
	info.FileSystemVersion, err = getFileSystemVersion(c.ssh)
	if err != nil {
		return info, err
	}

	ret, err := outputMinerRPC(c.ssh, "stats+version", "")
	if err != nil {
		return info, err
	}

	var resp struct {
		SGMultipleCMDResponse
		Version []SGVersionResponse `json:"version"`
		Stats   []SGStatsResponse   `json:"stats"`
	}
	err = json.Unmarshal(ret, &resp)
	if err != nil {
		return info, err
	}

	if !(len(resp.Version) == 1 && len(resp.Version[0].Version) == 1) {
		return info, fmt.Errorf("error sgminer RPC response")
	}
	version := resp.Version[0].Version[0]
	info.MinerDescription = version.Miner
	info.MinerVersion = version.SGMiner
	info.APIVersion = version.API

	if !(len(resp.Stats) != 0 && len(resp.Stats[0].Stats) != 0) {
		return info, fmt.Errorf("error sgminer RPC response")
	}
	stat := resp.Stats[0].Stats[0]

	info.ProductType, err = minerTypeFromAPIHWV(stat.HWV.String())
	if err != nil {
		return info, err
	}
	info.ProductVersion, err = minerVersionFromFWV(stat.FWV.String())
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

func getKernelVersion(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `uname -srv`)
	return string(bytes.TrimSpace(ret)), err
}

func getFileSystemVersion(client *ssh.Client) (string, error) {
	ret, err := outputRemoteShell(client, `uname -v`)
	return string(bytes.TrimSpace(ret)), err
}
