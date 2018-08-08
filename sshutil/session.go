package sshutil

import (
	"context"

	"golang.org/x/crypto/ssh"
)

// Session is x/crypto/ssh wrapper with timeouts
type Session struct {
	*ssh.Session
}

// OutputContext is call Session.Output with context.Context
func (s Session) OutputContext(ctx context.Context, cmd string) (output []byte, err error) {
	if ctx == nil {
		return s.Output(cmd)
	}

	rchan := make(chan struct{})
	go func() {
		defer close(rchan)
		output, err = s.Session.Output(cmd)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-rchan:
		return output, err
	}
}

// RunContext is call Session.Run with context.Context
func (s Session) RunContext(ctx context.Context, cmd string) (err error) {
	if ctx == nil {
		return s.Run(cmd)
	}

	rchan := make(chan struct{})
	go func() {
		defer close(rchan)
		err = s.Session.Run(cmd)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-rchan:
		return err
	}
}
