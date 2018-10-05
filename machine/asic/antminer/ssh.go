package antminer

import (
	"time"

	"github.com/ka2n/masminer/sshutil"
	"golang.org/x/crypto/ssh"
)

var (
	sshDialer sshutil.TimeoutDialer
)

// NewSSHClient returns *ssh.Client with default setting
func NewSSHClient(host string) (*ssh.Client, error) {
	return NewSSHClientTimeout(host, 0)
}

// NewSSHClientTimeout returns *ssh.Client with default setting with connection timeout
func NewSSHClientTimeout(host string, timeout time.Duration) (*ssh.Client, error) {
	var c Client
	addr, cfg := c.SSHConfig(host, timeout)
	return sshDialer.DialTimeout("tcp", addr, cfg, timeout)
}
