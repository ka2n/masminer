// MIT License
//
// Copyright © 2017 Jaime Pillora <dev@jpillora.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// https://github.com/jpillora/icmpscan/blob/master/names.go

package netscan

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

func LookupHostname(dnsClient *dns.Client, ipStr string) string {
	ip := net.ParseIP(ipStr)
	wg := sync.WaitGroup{}
	wg.Add(3)
	result := make(chan string, 3)
	//prepare query
	m := dns.Msg{}
	ipRev := strings.Split(ip.String(), ".")
	for i, j := 0, len(ipRev)-1; i < j; i, j = i+1, j-1 {
		ipRev[i], ipRev[j] = ipRev[j], ipRev[i]
	}
	target := strings.Join(ipRev, ".") + ".in-addr.arpa."
	m.SetQuestion(target, dns.TypePTR)
	//send to host (mdns)
	go func() {
		r, _, err := dnsClient.Exchange(&m, ip.String()+":5353")
		if err == nil {
			if len(r.Answer) > 0 {
				p := r.Answer[0].(*dns.PTR)
				result <- strings.TrimSuffix(strings.TrimSuffix(p.Ptr, "."), ".local")
			}
		}
		wg.Done()
	}()
	//send to router (dns)
	go func() {
		//give the host a slight headstart
		time.Sleep(5 * time.Millisecond)
		//use provided or guessed dns server address
		var server string
		if server == "" {
			b := make([]byte, 4)
			copy(b, []byte(ip.To4()))
			b[3] = 1
			server = net.IP(b).String()
		}
		r, _, err := dnsClient.Exchange(&m, server+":53")
		if err == nil {
			if len(r.Answer) > 0 {
				p := r.Answer[0].(*dns.PTR)
				result <- strings.TrimSuffix(p.Ptr, ".")
			}
		}
		wg.Done()
	}()
	//send to host (netbios name service)
	go func() {
		//give the host a slight headstart
		time.Sleep(5 * time.Millisecond)
		//use provided or guessed dns server address
		hostname, err := lookupNetBIOSName(ip)
		if err == nil {
			result <- hostname
		}
		wg.Done()
	}()
	//close after all 3 have returned
	go func() {
		wg.Wait()
		close(result)
	}()
	//use first result
	hostname := <-result
	return hostname
}

func lookupNetBIOSName(ip net.IP) (string, error) {
	m := dns.Msg{}
	m.SetQuestion("CKAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA.", 33)
	b, err := m.Pack()
	if err != nil {
		return "", err
	}
	conn, err := net.Dial("udp", ip.String()+":137")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	if _, err := conn.Write(b); err != nil {
		return "", err
	}
	buff := make([]byte, 512)
	n, err := conn.Read(buff)
	if err != nil {
		return "", err
	}
	if n < 12 {
		return "", fmt.Errorf("no header")
	}
	b = buff[:n]
	if m.Id != binary.BigEndian.Uint16(b[0:2]) {
		return "", fmt.Errorf("id mismatch")
	}
	//==== headers
	// flags := binary.BigEndian.Uint16(b[2:4])
	// questions := binary.BigEndian.Uint16(b[4:6])
	answers := binary.BigEndian.Uint16(b[6:8])
	// authority := binary.BigEndian.Uint16(b[8:10])
	// additional := binary.BigEndian.Uint16(b[10:12])
	if answers == 0 {
		return "", fmt.Errorf("no answers")
	}
	//==== answers
	b = b[12:]
	offset := 0
	for b[offset] != 0 {
		offset++
		if offset == len(b) {
			return "", fmt.Errorf("too short")
		}
	}
	// hostname := string(b[:offset])
	b = b[offset+1:]
	if len(b) < 12 {
		return "", fmt.Errorf("no answer")
	}
	// rtype := binary.BigEndian.Uint16(b[:2])
	// rclass := binary.BigEndian.Uint16(b[2:4])
	// ttl := binary.BigEndian.Uint32(b[4:8])
	// len := binary.BigEndian.Uint16(b[8:10])
	names := b[10]
	if names == 0 {
		return "", fmt.Errorf("no names")
	}
	b = b[11:]
	offset = 0
	for b[offset] != 0 {
		offset++
		if offset == len(b) {
			return "", fmt.Errorf("too short")
		}
	}
	netbiosName := strings.TrimSpace(string(b[:offset]))
	// macAddress := offset + 2 and onwards
	return netbiosName, nil
}
