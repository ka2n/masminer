// The MIT License (MIT)
//
// Copyright (c) 2018 Jessica Frazelle
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
// https://github.com/jessfraz/netscan

package netscan

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	timeout = time.Second
)

type CheckResult struct {
	IP    string
	Proto string
	Port  int
}

func checkReachable(proto, addr string) bool {
	c, err := net.DialTimeout(proto, addr, timeout)
	if err == nil {
		c.Close()
		return true
	}
	return false
}

func scanIP(ip string, beginPort int, endPort int) []CheckResult {
	result := make([]CheckResult, 0)
	protos := []string{"tcp"}
	for _, proto := range protos {
		for port := beginPort; port <= endPort; port++ {
			addr := fmt.Sprintf("%s:%d", ip, port)
			if checkReachable(proto, addr) {
				result = append(result, CheckResult{
					IP:    ip,
					Proto: proto,
					Port:  port,
				})
			}
		}
	}
	return result
}

// Scan : scan IP addrs in port range
func Scan(s string, beginPort int, endPort int) ([]CheckResult, error) {
	if beginPort > endPort {
		return nil, fmt.Errorf("End port can not be greater than the beginning port: %d > %d", endPort, beginPort)
	}

	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		ip = net.ParseIP(s)
		return scanIP(ip.String(), beginPort, endPort), nil
	}

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 200)
	result := make([]CheckResult, 0)
	mu := sync.Mutex{}

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		sem <- struct{}{}
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			defer func() { <-sem }()

			ret := scanIP(ip, beginPort, endPort)
			mu.Lock()
			defer mu.Unlock()
			result = append(result, ret...)
		}(ip.String())
	}
	wg.Wait()
	return result, nil
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
