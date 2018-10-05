package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ka2n/masminer/machine"
	"github.com/ka2n/masminer/machine/asic"
)

type config struct {
	ip       string
	hostname string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.ip, "ip", "", "Target IP Address(required)")
	flag.StringVar(&cfg.hostname, "host", "", "Hostname to determine what kind of hardware")
	flag.Parse()

	log.SetPrefix("[inspect] ")

	ctx := context.Background()
	client, err := asic.DialTimeout(machine.RemoteRig{
		IPAddr:   cfg.ip,
		Hostname: cfg.hostname,
	}, time.Second*10)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	info, err := client.RigInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", info)

	stat, err := client.RigStat(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", stat)
}
