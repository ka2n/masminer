package antminer

import (
	"context"
	"time"

	"github.com/ka2n/masminer/cgminerproxy"
	"github.com/ka2n/masminer/sshutil"
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

func outputRemoteShell(ctx context.Context, client *ssh.Client, in string) ([]byte, error) {
	sess, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	sw := sshutil.Session{Session: sess}
	return sw.OutputContext(ctx, in)
}

func runRemoteShell(ctx context.Context, client *ssh.Client, in string) error {
	sess, err := client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	sw := sshutil.Session{Session: sess}
	return sw.RunContext(ctx, in)
}

func outputMinerRPC(ctx context.Context, client *ssh.Client, command, argument string) ([]byte, error) {
	var out []byte
	var err error

	done := make(chan struct{})
	go func() {
		defer close(done)
		out, err = outputMinerRPCInner(client, command, argument)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return out, err
	}
}

func outputMinerRPCInner(client *ssh.Client, command, argument string) ([]byte, error) {
	conn, err := client.Dial("tcp", "127.0.0.1:4028")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	proxy := new(cgminerproxy.CGMinerProxy)
	ret, err := proxy.RunCommandConn(conn, command, argument)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
