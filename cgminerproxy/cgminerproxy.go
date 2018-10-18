package cgminerproxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

// CGMinerProxy is a simple proxy between ASIC
type CGMinerProxy struct {
	HostPort string
	Timeout  time.Duration
}

// Usable command
const (
	CommandSummary  = "summary"
	CommandRestart  = "restart"
	CommandQuit     = "quit"
	CommandDevs     = "devs"
	CommandChipstat = "chipstat"
)

// New returns a CGMinerProxy
func New(hostname string, port int) *CGMinerProxy {
	return &CGMinerProxy{
		HostPort: fmt.Sprintf("%s:%d", hostname, port),
		Timeout:  time.Second * 10,
	}
}

// RunCommand calls RPC command
func (proxy *CGMinerProxy) RunCommand(command, argument string) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", proxy.HostPort, proxy.Timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return proxy.RunCommandConn(conn, command, argument)
}

// RunCommandConn calls RPC command on existing net.Conn
func (proxy *CGMinerProxy) RunCommandConn(conn net.Conn, command, argument string) ([]byte, error) {
	type Req struct {
		Command   string `json:"command"`
		Parameter string `json:"parameter,omitempty"`
	}

	req := Req{Command: command, Parameter: argument}
	reqb, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(conn, "%s", reqb)

	resp, err := bufio.NewReader(conn).ReadBytes('\x00')
	if err != nil && err != io.EOF {
		return nil, err
	}
	return bytes.TrimRight(resp, "\x00"), nil
}
