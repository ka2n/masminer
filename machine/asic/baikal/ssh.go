package baikal

import (
	"os"

	"github.com/ka2n/masminer/cgminerproxy"
	"golang.org/x/crypto/ssh"
)

// NewSSHClient returns *ssh.Client with default setting
func NewSSHClient(host string) (*ssh.Client, error) {
	cfg := &ssh.ClientConfig{
		User: "baikal",
		Auth: []ssh.AuthMethod{
			ssh.Password("baikal"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
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
	sess.Stderr = os.Stderr
	return sess.Run(in)
}

func outputMinerRPC(client *ssh.Client, command, argument string) ([]byte, error) {
	conn, err := client.Dial("tcp", "127.0.0.1:4028")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	proxy := new(cgminerproxy.CGMinerProxy)
	return proxy.RunCommandConn(conn, command, argument)
}
