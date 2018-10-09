package base

import (
	"bytes"
	"context"
	"sync"

	mnet "github.com/ka2n/masminer/net"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	SSH *ssh.Client
	MU  sync.RWMutex
}

func (c *Client) SetSSH(client *ssh.Client) {
	c.MU.Lock()
	defer c.MU.Unlock()
	c.SSH = client
}

func (c *Client) Setup() error {
	return nil
}

func (c *Client) Close() error {
	return c.SSH.Close()
}

func GetMacAddr(ctx context.Context, client *ssh.Client) (string, error) {
	ret, err := OutputRemoteShell(ctx, client, `ip link show eth0 | grep -o 'link/.*' | cut -d' ' -f2`)
	return string(bytes.TrimSpace(ret)), err
}

func GetHostname(ctx context.Context, client *ssh.Client) (string, error) {
	ret, err := OutputRemoteShell(ctx, client, `hostname`)
	return string(bytes.TrimSpace(ret)), err
}

func GetKernelVersion(ctx context.Context, client *ssh.Client) (string, error) {
	ret, err := OutputRemoteShell(ctx, client, `uname -srv`)
	return string(bytes.TrimSpace(ret)), err
}

func GetFileSystemVersion(ctx context.Context, client *ssh.Client) (string, error) {
	ret, err := OutputRemoteShell(ctx, client, `uname -v`)
	return string(bytes.TrimSpace(ret)), err
}

func GetUptimeSeconds(ctx context.Context, client *ssh.Client) (string, error) {
	cmd := "cut -d \".\" -f 1 /proc/uptime"
	ret, err := OutputRemoteShell(ctx, client, cmd)
	return string(bytes.TrimSpace(ret)), err
}

func GetCPUTemp(ctx context.Context, client *ssh.Client) (string, error) {
	cmd := `cat /sys/class/thermal/thermal_zone*/temp | awk '{sum+=$1} END {print sum/NR}'`
	ret, err := OutputRemoteShell(ctx, client, cmd)
	return string(bytes.TrimSpace(ret)), err
}

func GetIPAddr(ctx context.Context, client *ssh.Client) (string, error) {
	ret, err := OutputRemoteShell(ctx, client, `ip a show eth0 | grep -o 'inet\s.*' | cut -d' ' -f2`)
	if err != nil {
		return string(ret), err
	}
	return mnet.ParseIPAddr(string(bytes.TrimSpace(ret)))
}
