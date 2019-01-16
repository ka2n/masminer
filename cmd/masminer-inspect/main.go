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
	timeout  time.Duration
}

func main() {
	var cfg config
	flag.StringVar(&cfg.ip, "ip", "", "Target IP Address(required)")
	flag.StringVar(&cfg.hostname, "host", "", "Hostname to determine what kind of hardware")
	flag.DurationVar(&cfg.timeout, "timeout", time.Second*10, "timeout")
	flag.Parse()

	log.SetPrefix("[inspect] ")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	client, err := asic.DialTimeout(machine.RemoteRig{
		IPAddr:   cfg.ip,
		Hostname: cfg.hostname,
	}, cfg.timeout)
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
