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
	sess, err := newSessionContext(ctx, client)
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	return sess.OutputContext(ctx, in)
}

func RunRemoteShell(ctx context.Context, client *ssh.Client, in string) error {
	sess, err := newSessionContext(ctx, client)
	if err != nil {
		return err
	}
	defer sess.Close()

	return sess.RunContext(ctx, in)
}

func OutputMinerRPC(ctx context.Context, d Dialer, command, argument string) ([]byte, error) {
	conn, err := d.Dial("tcp", "127.0.0.1:4028")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	proxy := new(cgminerproxy.CGMinerProxy)
	ret, err := proxy.RunCommandConn(conn, command, argument)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func newSessionContext(ctx context.Context, client *ssh.Client) (*sshutil.Session, error) {
	var sess *ssh.Session
	var err error

	done := make(chan struct{})
	go func() {
		defer close(done)
		sess, err = client.NewSession()
	}()

	select {
	case <-ctx.Done():
		if sess != nil {
			sess.Close()
		}
		return nil, ctx.Err()
	case <-done:
		return &sshutil.Session{Session: sess}, err
	}
}
