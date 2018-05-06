package asic

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/ka2n/masminer/inspect"
	"github.com/ka2n/masminer/netscan"
	"github.com/mostlygeek/arp"
)

// DiscoverByIPScan : find ASICs by scan network
func DiscoverByIPScan(ctx context.Context, networks []string, portBegin int, portEnd int) ([]inspect.RemoteRig, error) {
	// Scan IPs
	dnsClient := &dns.Client{}
	result := make([]inspect.RemoteRig, 0)
	for _, netmask := range networks {
		ret, err := netscan.Scan(netmask, portBegin, portEnd)
		if err != nil {
			continue
		}
		for _, r := range ret {
			// Collect mac addr from ARP table
			mac := arp.Search(r.IP)

			var rig inspect.RemoteRig
			if mac != "" {
				rig.Name = inspect.ShortName(mac)
				rig.MACAddr = mac
			}
			rig.IPAddr = r.IP
			rig.Hostname = netscan.LookupHostname(dnsClient, r.IP)
			result = append(result, rig)
		}
	}
	return result, nil
}

// DiscoverByMCast find ASICs by broadcast message
func DiscoverByMCast(ctx context.Context, mcastCode string, mcastAddr string, mcastListenPort int, timeout time.Duration, handler func(inspect.RemoteRig)) error {
	bufSize := 8192

	// Create socket to broadcast message
	addrBroad, err := net.ResolveUDPAddr("udp", mcastAddr)
	if err != nil {
		return err
	}
	connBroad, err := net.DialUDP("udp", nil, addrBroad)
	if err != nil {
		return err
	}
	defer connBroad.Close()

	// Create listener for receive message
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", mcastListenPort))
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	deadline := time.Now().Add(timeout)
	conn.SetReadBuffer(bufSize)
	conn.SetDeadline(deadline)

	dnsClient := &dns.Client{}

	readErr := make(chan error)
	go func() {
		for {
			buffer := make([]byte, bufSize)
			_, src, err := conn.ReadFromUDP(buffer)
			if err != nil {
				readErr <- err
				return
			}

			// Parse received message
			var rig inspect.RemoteRig
			buffer = bytes.Trim(buffer, "\x00")
			body := string(buffer)
			ans := strings.SplitN(body, "-", 4)
			if len(ans) >= 3 && ans[0] == "cgm" && ans[1] == mcastCode {
				mac := arp.Search(src.IP.String())
				if len(ans) > 3 && len(ans[3]) > 0 {
					mac = ans[3]
				}
				if mac == "" {
					continue
				}
				rig.Name = mac
				rig.IPAddr = src.IP.String()
				rig.Hostname = netscan.LookupHostname(dnsClient, rig.IPAddr)
				handler(rig)
			}
		}
	}()

	// broadcast message
	buf := fmt.Sprintf("cgminer-%s-%d", mcastCode, mcastListenPort)
	connBroad.Write([]byte(buf))
	connBroad.Close()

	<-readErr
	return nil
}
