package cgminerproxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
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
func (proxy *CGMinerProxy) RunCommand(command, argument string) (string, error) {
	conn, err := net.DialTimeout("tcp", proxy.HostPort, proxy.Timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	type Req struct {
		Command   string `json:"command"`
		Parameter string `json:"parameter,omitempty"`
	}

	req := Req{Command: command, Parameter: argument}
	reqb, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(conn, "%s", reqb)
	resp, err := bufio.NewReader(conn).ReadString('\x00')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(resp, "\x00"), nil
}
