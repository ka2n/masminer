package asic

import (
	"fmt"
	"strings"

	"github.com/ka2n/masminer/machine"
	"github.com/ka2n/masminer/machine/asic/antminer"
	"github.com/ka2n/masminer/machine/asic/baikal"
	"golang.org/x/crypto/ssh"
)

// Dial : get asic client from RemoteRig
func Dial(r machine.RemoteRig) (Client, error) {
	rw, err := dialByHostname(r.IPAddr, r.Hostname)
	if err != nil {
		return nil, err
	} else if rw != nil {
		return rw, nil
	}

	rw, err = dialBySSH(r.IPAddr)
	if err != nil {
		return nil, err
	} else if rw != nil {
		return rw, nil
	}

	return nil, fmt.Errorf("unknown minerType")
}

func dialByHostname(ipAddr string, hostname string) (Client, error) {
	hostLower := strings.ToLower(hostname)
	switch {
	case strings.Contains(hostLower, "baikal"):
		var client baikal.Client
		sc, err := baikal.NewSSHClient(ipAddr)
		if err != nil {
			return nil, err
		}
		client.SetSSH(sc)
		return &client, nil
	case strings.Contains(hostLower, "antminer"):
		var client antminer.Client
		sc, err := antminer.NewSSHClient(ipAddr)
		if err != nil {
			return nil, err
		}
		client.SetSSH(sc)
		return &client, nil
	default:
		return nil, nil
	}
}

func dialBySSH(ipAddr string) (Client, error) {
	var sc *ssh.Client
	var err error

	sc, err = antminer.NewSSHClient(ipAddr)
	if err == nil {
		var client antminer.Client
		client.SetSSH(sc)
		return &client, nil
	}

	sc, err = baikal.NewSSHClient(ipAddr)
	if err == nil {
		var client baikal.Client
		client.SetSSH(sc)
		return &client, nil
	}

	return nil, nil
}
