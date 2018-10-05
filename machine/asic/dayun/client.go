package dayun

import (
	"context"
	"time"

	"github.com/ka2n/masminer/machine"
	"github.com/ka2n/masminer/machine/asic/base"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	base.Client

	systemInfo *SystemInfo
}

func (c *Client) SSHConfig(host string, timeout time.Duration) (string, *ssh.ClientConfig) {
	return host + ":22", &ssh.ClientConfig{
		User: "fa",
		Auth: []ssh.AuthMethod{
			ssh.Password("fa"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}
}

func (c *Client) MineStop(context.Context) error {
	panic("not implemented")
}

func (c *Client) MineStart(context.Context) error {
	panic("not implemented")
}

func (c *Client) Restart(context.Context) error {
	panic("not implemented")
}

func (c *Client) Reboot(context.Context) error {
	panic("not implemented")
}

func (c *Client) RigInfo(ctx context.Context) (machine.RigInfo, error) {
	var info machine.RigInfo
	si, err := c.GetSystemInfoContext(ctx)
	if err != nil {
		return info, err
	}

	info.Rig.IPAddr = si.IPAddr
	info.Rig.Hostname = si.Hostname
	info.Rig.Name = machine.ShortName(si.MACAddr)
	info.Rig.MACAddr = si.MACAddr
	info.Model = si.ProductType
	info.Manufacture = manufactureName
	info.FirmwareVersion = si.DashboardVersion
	info.MinerType = si.MinerDescription
	info.MinerVersion = si.MinerVersion
	info.UptimeSeconds = si.UptimeSeconds
	return info, nil
}

func (c *Client) RigStat(context.Context) (machine.RigStat, error) {
	return machine.RigStat{}, nil
}

func (c *Client) MinerSetting(context.Context) (machine.MinerSetting, error) {
	panic("not implemented")
}

func (c *Client) SetMinerSetting(context.Context, machine.MinerSetting) error {
	panic("not implemented")
}
