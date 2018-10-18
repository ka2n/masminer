package asic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ka2n/masminer/machine"
	"github.com/ka2n/masminer/machine/asic/antminer"
	"github.com/ka2n/masminer/machine/asic/baikal"
	"github.com/ka2n/masminer/sshutil"
)

// Dial : get asic client from RemoteRig
func Dial(r machine.RemoteRig) (Client, error) {
	return DialTimeout(r, 0)
}

// DialTimeout : get asic client from RemoteRig with connection timeout
func DialTimeout(r machine.RemoteRig, timeout time.Duration) (Client, error) {
	c, err := dial(r, timeout)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c, c.Setup(ctx)
}

func dial(r machine.RemoteRig, timeout time.Duration) (Client, error) {
	rw, err := dialByHostname(r.IPAddr, r.Hostname, timeout)
	if err != nil {
		return nil, err
	} else if rw != nil {
		return rw, nil
	}

	rw, err = dialBySSH(r.IPAddr, timeout)
	if err != nil {
		return nil, err
	} else if rw != nil {
		return rw, nil
	}

	return nil, fmt.Errorf("unknown minerType")
}

func dialByHostname(ipAddr string, hostname string, timeout time.Duration) (Client, error) {
	hostLower := strings.ToLower(hostname)
	var (
		client Client
		dialer sshutil.TimeoutDialer
		err    error
	)

	switch {
	case strings.Contains(hostLower, "baikal"):
		client = &baikal.Client{}
	case strings.Contains(hostLower, "antminer"):
		client = &antminer.Client{}
	default:
		return nil, nil
	}

	addr, cfg := client.SSHConfig(ipAddr)
	cfg.Timeout = timeout
	conn, err := dialer.DialTimeout("tcp", addr, cfg, timeout)
	if err != nil {
		return nil, errors.Wrap(err, "antminer client dial failed")
	}
	client.SetSSH(conn)
	return client, nil
}

func dialBySSH(ipAddr string, timeout time.Duration) (Client, error) {
	var dialer sshutil.TimeoutDialer

	tries := []Client{
		&antminer.Client{},
		&baikal.Client{},
	}
	for _, t := range tries {
		addr, cfg := t.SSHConfig(ipAddr)
		cfg.Timeout = timeout
		conn, err := dialer.DialTimeout("tcp", addr, cfg, timeout)
		if err != nil {
			continue
		}
		t.SetSSH(conn)
		return t, nil
	}
	return nil, nil
}
