package sshutil

import (
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	defaultTimeout = time.Second * 120
)

// TimeoutDialer SSH dialer with timeout
// borrowed from https://github.com/cjcullen/kubernetes/blob/cde4f6d613b42b06946a6b68b0d55f11e8aedadb/pkg/ssh/ssh.go
type TimeoutDialer struct{}

func (d *TimeoutDialer) DialTimeout(network, addr string, config *ssh.ClientConfig, timeout time.Duration) (client *ssh.Client, err error) {
	if timeout == 0 {
		timeout = defaultTimeout
	}

	if config.Timeout == 0 {
		config.Timeout = timeout
	}

	conn, err := net.DialTimeout(network, addr, config.Timeout)
	if err != nil {
		return nil, err
	}

	if config.Timeout > 0 {
		conn.SetReadDeadline(time.Now().Add(config.Timeout))
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return nil, err
	}

	if config.Timeout > 0 {
		conn.SetReadDeadline(time.Time{})
	}
	return ssh.NewClient(c, chans, reqs), nil
}
