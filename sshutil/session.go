package sshutil

import (
	"context"
	"io"

	"golang.org/x/crypto/ssh"
)

// Session is x/crypto/ssh wrapper with timeouts
type Session struct {
	*ssh.Session
}

func (s Session) zombieHack() (io.WriteCloser, error) {
	// Request ptty to prevent zombie process after disconnection
	err := s.RequestPty("xterm", 80, 40, ssh.TerminalModes{})
	if err != nil {
		return nil, err
	}
	return s.StdinPipe()
}

func (s Session) Output(cmd string) ([]byte, error) {
	stdin, err := s.zombieHack()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()

	return s.Session.Output(cmd)
}

func (s Session) Run(cmd string) error {
	stdin, err := s.zombieHack()
	if err != nil {
		return err
	}
	defer stdin.Close()

	return s.Session.Run(cmd)
}

// OutputContext is call Session.Output with context.Context
func (s Session) OutputContext(ctx context.Context, cmd string) (output []byte, err error) {
	if ctx == nil {
		return s.Output(cmd)
	}
	rchan := make(chan struct{})
	go func() {
		defer close(rchan)
		output, err = s.Output(cmd)
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
		err = s.Run(cmd)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-rchan:
		return err
	}
}
