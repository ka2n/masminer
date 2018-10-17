package sshutil

import (
	"fmt"
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

	done := make(chan struct{}, 1)
	go func() {
		defer close(done)
		client, err = ssh.Dial(network, addr, config)
	}()

	select {
	case <-done:
		return client, err
	case <-time.After(timeout):
		if client != nil {
			client.Close()
		}
		return nil, fmt.Errorf("timed out dialing %s:%s", network, addr)
	}
}
