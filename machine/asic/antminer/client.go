package antminer

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/ka2n/masminer/inspect"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	mu         sync.RWMutex
	ssh        *ssh.Client
	systemInfo *SystemInfo
}

func (c *Client) MineStop(ctx context.Context) error {
	return runRemoteShell(c.ssh, fmt.Sprintf(minerInitdCMD, "stop"))
}

func (c *Client) MineStart(ctx context.Context) error {
	return runRemoteShell(c.ssh, fmt.Sprintf(minerInitdCMD, "start"))
}

func (c *Client) Restart(ctx context.Context) error {
	return runRemoteShell(c.ssh, fmt.Sprintf(minerInitdCMD, "restart"))
}

func (c *Client) Reboot(ctx context.Context) error {
	return runRemoteShell(c.ssh, "shutdown -r +5")
}

func (c *Client) SetSSH(client *ssh.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ssh = client
}

func (c *Client) Close() error {
	return c.ssh.Close()
}

func (c *Client) RigInfo(ctx context.Context) (inspect.RigInfo, error) {
	var info inspect.RigInfo
	si, err := c.GetSystemInfo()
	if err != nil {
		return info, err
	}
	info.Rig.IPAddr = c.ssh.RemoteAddr().String()
	info.Rig.Hostname = si.Hostname
	info.Rig.MACAddr = si.MACAddr
	info.Rig.Name = inspect.ShortName(si.MACAddr)
	info.MinerType = si.ProductType
	info.Manufacture = manufactureName
	info.HardwareVersion = strings.Join(si.HardwareVersions, ",")
	info.FirmwareVersion = si.FileSystemVersion
	info.MinerVersion = si.CGMinerVersion
	info.Algos = Algos(si.ProductType)
	return info, nil
}

func (c *Client) RigStat(ctx context.Context) (inspect.RigStat, error) {
	var stat inspect.RigStat

	ms, err := c.GetStats()
	if err != nil {
		return stat, err
	}

	stat.GHS5s = ms.Summary.GHS5s
	stat.GHSAvarage = ms.Summary.GHSAvarage
	stat.Accepted = ms.Summary.Accepted
	stat.Rejected = ms.Summary.Rejected
	stat.HardwareErrors = ms.Summary.HardwareErrors
	stat.Utility = ms.Summary.Utility

	stat.Devices = make([]inspect.DeviceStat, len(ms.Devs.Chains))
	for i, c := range ms.Devs.Chains {
		var st inspect.DeviceStat
		st.TempChip = c.TempChip
		st.TempPCB = c.TempPCB
		st.Frequency = c.Freq
		st.Chips = strings.Count(c.Status, "o")
		st.HardwareErrors = c.Hw
		st.Hashrate = c.Rate
		stat.Devices[i] = st
	}

	stat.Pools = make([]inspect.PoolStat, len(ms.Pools))
	for i, pl := range ms.Pools {
		var p inspect.PoolStat
		p.URL = pl.URL
		p.User = pl.User
		p.Status = pl.Status
		p.StratumActive = pl.StratumActive == "true"
		p.Priority, _ = strconv.Atoi(pl.Priority)
		p.Getworks = pl.Getworks
		p.Accepted = pl.Accepted
		p.Rejected = pl.Rejected
		p.Discarded = pl.Discarded
		p.Stale = pl.Stale
		p.DifficultyAccepted = pl.DifficultyAccepted
		p.DifficultyRejected = pl.DifficultyRejected
		p.DifficultyStale = pl.DifficultyStale
		p.LastShareDifficulty = pl.LastShareDifficulty
		p.LastShareTime = pl.LastShareTime
		stat.Pools[i] = p
	}

	return stat, nil
}

func (c *Client) MinerSetting(ctx context.Context) (inspect.MinerSetting, error) {
	ms, err := c.GetMinerSetting()
	if err != nil {
		return inspect.MinerSetting{}, err
	}
	return ms.CommonMinerSetting(), nil
}

func (c *Client) SetMinerSetting(ctx context.Context, setting inspect.MinerSetting) error {
	var ms MinerSetting
	if err := ms.LoadCommonMinerSetting(setting); err != nil {
		return err
	}
	return c.WriteCGMinerSetting(ms)
}
