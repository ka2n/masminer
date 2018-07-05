package antminer

import (
	"time"

	"github.com/ka2n/masminer/cgminerproxy"
	"golang.org/x/crypto/ssh"
)

// NewSSHClient returns *ssh.Client with default setting
func NewSSHClient(host string) (*ssh.Client, error) {
	return NewSSHClientTimeout(host, 0)
}

// NewSSHClientTimeout returns *ssh.Client with default setting with connection timeout
func NewSSHClientTimeout(host string, timeout time.Duration) (*ssh.Client, error) {
	cfg := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("admin"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}
	return ssh.Dial("tcp", host+":22", cfg)
}

func outputRemoteShell(client *ssh.Client, in string) ([]byte, error) {
	sess, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()
	return sess.Output(in)
}

func runRemoteShell(client *ssh.Client, in string) error {
	sess, err := client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	return sess.Run(in)
}

func outputMinerRPC(client *ssh.Client, command, argument string) (string, error) {
	conn, err := client.Dial("tcp", "127.0.0.1:4028")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	proxy := new(cgminerproxy.CGMinerProxy)
	ret, err := proxy.RunCommandConn(conn, command, argument)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}
