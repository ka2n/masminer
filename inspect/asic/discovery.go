package asic

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ka2n/masminer/inspect"

	"github.com/ka2n/masminer/netscan"
	"github.com/mostlygeek/arp"
)

// DiscoverByIPScan : find ASICs by scan network
func DiscoverByIPScan(ctx context.Context, ipPrefixes []string, portBegin int, portEnd int) ([]RemoteRig, error) {
	// Collect networks
	networks := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			switch a.(type) {
			case *net.IPNet:
				if v4 := a.(*net.IPNet).IP.To4(); v4 != nil && !v4.IsLoopback() {
					_, ipnet, _ := net.ParseCIDR(a.String())
					ipnetStr := ipnet.String()
					if len(ipPrefixes) > 0 {
						for _, prefix := range ipPrefixes {
							if strings.HasPrefix(ipnetStr, prefix) {
								networks = append(networks, ipnetStr)
							}
						}
					} else {
						networks = append(networks, ipnetStr)
					}
				}
			default:
			}
		}
	}

	// Scan IPs
	result := make([]RemoteRig, 0)
	for _, netmask := range networks {
		ret, err := netscan.Scan(netmask, portBegin, portEnd)
		if err != nil {
			continue
		}
		for _, r := range ret {
			// Collect mac addr from ARP table
			mac := arp.Search(r.IP)
			if mac == "" {
				continue
			}
			result = append(result, RemoteRig{
				Name:    inspect.ShortName(mac),
				MACAddr: mac,
				IPAddr:  r.IP,
				APIPort: r.Port,
			})
		}
	}
	return result, nil
}

// DiscoverByMCast find ASICs by broadcast message
func DiscoverByMCast(ctx context.Context, mcastCode string, mcastAddr string, mcastListenPort int, timeout time.Duration, handler func(RemoteRig)) error {
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
			var rig RemoteRig
			buffer = bytes.Trim(buffer, "\x00")
			body := string(buffer)
			ans := strings.SplitN(body, "-", 4)
			if len(ans) >= 3 && ans[0] == "cgm" && ans[1] == mcastCode {
				port, err := strconv.ParseInt(ans[2], 10, 64)
				if err != nil {
					fmt.Println(err)
					continue
				}
				mac := arp.Search(src.IP.String())
				if len(ans) > 3 && len(ans[3]) > 0 {
					mac = ans[3]
				}
				if mac == "" {
					continue
				}
				rig.Name = mac
				rig.APIPort = int(port)
				rig.IPAddr = src.IP.String()
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
