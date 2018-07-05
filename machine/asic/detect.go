package asic

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ka2n/masminer/machine"
	"github.com/ka2n/masminer/machine/asic/antminer"
	"github.com/ka2n/masminer/machine/asic/baikal"
	"golang.org/x/crypto/ssh"
)

// Dial : get asic client from RemoteRig
func Dial(r machine.RemoteRig) (Client, error) {
	return DialTimeout(r, 0)
}

// DialTimeout : get asic client from RemoteRig with connection timeout
func DialTimeout(r machine.RemoteRig, timeout time.Duration) (Client, error) {
	c, err := dial(r, 0)
	if err != nil {
		return nil, err
	}
	return c, c.Setup()
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
	switch {
	case strings.Contains(hostLower, "baikal"):
		var client baikal.Client
		sc, err := baikal.NewSSHClientTimeout(ipAddr, timeout)
		if err != nil {
			return nil, errors.Wrap(err, "baikal client dial failed")
		}
		client.SetSSH(sc)
		return &client, nil
	case strings.Contains(hostLower, "antminer"):
		var client antminer.Client
		sc, err := antminer.NewSSHClientTimeout(ipAddr, timeout)
		if err != nil {
			return nil, errors.Wrap(err, "antminer client dial failed")
		}
		client.SetSSH(sc)
		return &client, nil
	default:
		return nil, nil
	}
}

func dialBySSH(ipAddr string, timeout time.Duration) (Client, error) {
	var sc *ssh.Client
	var err error

	sc, err = antminer.NewSSHClientTimeout(ipAddr, timeout)
	if err == nil {
		var client antminer.Client
		client.SetSSH(sc)
		return &client, nil
	}

	sc, err = baikal.NewSSHClientTimeout(ipAddr, timeout)
	if err == nil {
		var client baikal.Client
		client.SetSSH(sc)
		return &client, nil
	}

	return nil, nil
}
