package antminer

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/ka2n/masminer/machine"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	mu         sync.RWMutex
	ssh        *ssh.Client
	systemInfo *SystemInfo

	minerType  string
	configPath string

	ipCMDPath string

	summaryCMD string
	poolsCMD   string
	statsCMD   string
	versionCMD string
	initdCMD   string
}

func (c *Client) Setup() error {
	return c.SetupContext(nil)
}

func (c *Client) SetupContext(ctx context.Context) error {
	// Check miner is cgminer or bmminer
	switch {
	case runRemoteShell(ctx, c.ssh, "which cgminer-api") == nil:
		c.configPath = minerConfigPath
		c.summaryCMD = minerAPISummaryCMD
		c.poolsCMD = minerAPIPoolsCMD
		c.statsCMD = minerAPIStatsCMD
		c.versionCMD = minerAPIVersionCMD
		c.initdCMD = minerInitdCMD
		c.minerType = "CGMiner"
	case runRemoteShell(ctx, c.ssh, "which bmminer-api") == nil:
		c.configPath = minerBMMinerConfigPath
		c.summaryCMD = minerBMMinerAPISummaryCMD
		c.poolsCMD = minerBMMinerAPIPoolsCMD
		c.statsCMD = minerBMMinerAPIStatsCMD
		c.versionCMD = minerBMMinerAPIVersionCMD
		c.initdCMD = minerBMMinerInitdCMD
		c.minerType = "BMMiner"
	default:
		return fmt.Errorf("cannot detect miner program type")
	}

	switch {
	case runRemoteShell(ctx, c.ssh, "type /bin/ip") == nil:
		c.ipCMDPath = "/bin/ip"
	case runRemoteShell(ctx, c.ssh, "type /sbin/ip") == nil:
		c.ipCMDPath = "/sbin/ip"
	default:
		return fmt.Errorf("cannot detect ip command path")
	}

	return nil
}

func (c *Client) MineStop(ctx context.Context) error {
	return runRemoteShell(ctx, c.ssh, fmt.Sprintf(c.initdCMD, "stop"))
}

func (c *Client) MineStart(ctx context.Context) error {
	return runRemoteShell(ctx, c.ssh, fmt.Sprintf(c.initdCMD, "start"))
}

func (c *Client) Restart(ctx context.Context) error {
	return runRemoteShell(ctx, c.ssh, fmt.Sprintf(c.initdCMD, "restart"))
}

func (c *Client) Reboot(ctx context.Context) error {
	return runRemoteShell(ctx, c.ssh, "shutdown -r +5")
}

func (c *Client) SetSSH(client *ssh.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ssh = client
}

func (c *Client) Close() error {
	return c.ssh.Close()
}

func (c *Client) RigInfo(ctx context.Context) (machine.RigInfo, error) {
	var info machine.RigInfo
	si, err := c.GetSystemInfoContext(ctx)
	if err != nil {
		return info, err
	}
	info.Rig.IPAddr = si.IPAddr
	info.Rig.Hostname = si.Hostname
	info.Rig.MACAddr = si.MACAddr
	info.Rig.Name = machine.ShortName(si.MACAddr)
	info.Manufacture = manufactureName
	info.Model = si.Model
	info.HardwareVersion = strings.Join(si.HardwareVersions, ",")
	info.FirmwareVersion = si.FileSystemVersion
	info.MinerVersion = si.MinerVersion
	info.MinerType = si.MinerType
	info.Algos = Algos(si.Model)
	info.UptimeSeconds = si.UptimeSeconds
	return info, nil
}

func (c *Client) RigStat(ctx context.Context) (machine.RigStat, error) {
	var stat machine.RigStat

	ms, err := c.GetStatsContext(ctx)
	if err != nil {
		return stat, err
	}

	stat.GHS5s = ms.Summary.GHS5s
	stat.GHSAvarage = ms.Summary.GHSAvarage
	stat.Accepted = ms.Summary.Accepted
	stat.Rejected = ms.Summary.Rejected
	stat.HardwareErrors = ms.Summary.HardwareErrors
	stat.Utility = ms.Summary.Utility

	stat.Devices = make([]machine.DeviceStat, len(ms.Devs.Chains))
	for i, c := range ms.Devs.Chains {
		var st machine.DeviceStat
		st.TempChip = c.TempChip
		st.TempPCB = c.TempPCB
		st.Frequency = c.Freq
		st.Chips = strings.Count(c.Status, "o")
		st.HardwareErrors = c.Hw
		st.Hashrate = c.Rate
		stat.Devices[i] = st
	}

	stat.Pools = make([]machine.PoolStat, len(ms.Pools))
	for i, pl := range ms.Pools {
		var p machine.PoolStat
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

func (c *Client) MinerSetting(ctx context.Context) (machine.MinerSetting, error) {
	ms, err := c.GetMinerSettingContext(ctx)
	if err != nil {
		return machine.MinerSetting{}, err
	}
	return ms.CommonMinerSetting(), nil
}

func (c *Client) SetMinerSetting(ctx context.Context, setting machine.MinerSetting) error {
	var ms MinerSetting
	if err := ms.LoadCommonMinerSetting(setting); err != nil {
		return err
	}
	return c.WriteCGMinerSettingContext(ctx, ms)
}
