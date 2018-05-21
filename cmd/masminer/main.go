package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ka2n/masminer/machine"
	"github.com/ka2n/masminer/machine/asic"
)

var (
	network = flag.String("network", "", "")
	port    = flag.Int("port", 4028, "")
)

func main() {
	flag.Parse()

	networks := make([]string, 0)
	if network != nil && *network != "" {
		nws := strings.Split(*network, ",")
		for _, nw := range nws {
			nw = strings.TrimSpace(nw)
			if nw != "" {
				networks = append(networks, nw)
			}
		}
	}

	ctx := context.Background()
	rigs, err := asic.DiscoverByIPScan(ctx, networks, *port, *port)
	if err != nil {
		panic(err)
	}

	if len(rigs) == 0 {
		fmt.Println("no rigs found(or cannot get MAC address)")
		return
	}

	for i, r := range rigs {
		fmt.Printf("#%d %s IP: %s, Host: %s\n", i, r.Name, r.IPAddr, r.Hostname)
	}
	watchRigs(rigs)
}

func watchRigs(rigs []machine.RemoteRig) error {
	ctx := context.Background()
	tick := time.NewTicker(time.Second * 10)
	defer tick.Stop()

	result := make(map[string]machine.RigInfo)
	mu := sync.Mutex{}

	for {
		<-tick.C
		var wg sync.WaitGroup
		wg.Add(len(rigs))
		fmt.Println("Start fetching...")
		now := time.Now()
		for _, rig := range rigs {
			rig := rig
			go func() {
				defer wg.Done()
				var client asic.Client
				var err error
				client, err = asic.Dial(rig)
				if err != nil {
					panic(err)
				}
				defer client.Close()
				stat, err := client.RigInfo(ctx)
				if err != nil {
					panic(err)
				}
				mu.Lock()
				defer mu.Unlock()
				result[stat.Rig.Name] = stat
			}()
		}
		wg.Wait()

		fmt.Println("========")
		for name, r := range result {
			fmt.Printf("%s - %s [%s], %s\n", name, r.Model, r.HardwareVersion, r.FirmwareVersion)
		}
		fmt.Println("========", time.Now().Sub(now))
	}
}
