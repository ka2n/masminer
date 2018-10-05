package base

import (
	"context"
	"net"

	"github.com/ka2n/masminer/cgminerproxy"
	"github.com/ka2n/masminer/sshutil"
	"golang.org/x/crypto/ssh"
)

type Dialer interface {
	Dial(n, addr string) (net.Conn, error)
}

func OutputRemoteShell(ctx context.Context, client *ssh.Client, in string) ([]byte, error) {
	sess, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	sw := sshutil.Session{Session: sess}
	return sw.OutputContext(ctx, in)
}

func RunRemoteShell(ctx context.Context, client *ssh.Client, in string) error {
	sess, err := client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	sw := sshutil.Session{Session: sess}
	return sw.RunContext(ctx, in)
}

func OutputMinerRPC(ctx context.Context, d Dialer, command, argument string) ([]byte, error) {
	var out []byte
	var err error

	done := make(chan struct{})
	go func() {
		defer close(done)
		out, err = outputMinerRPCInner(d, command, argument)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return out, err
	}
}

func outputMinerRPCInner(d Dialer, command, argument string) ([]byte, error) {
	conn, err := d.Dial("tcp", "127.0.0.1:4028")
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
