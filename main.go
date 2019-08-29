package main

import (
	"flag"
	"log"
	"time"
	"wen/service-prober/pkg/cloudprober"
)

var (
	server    = flag.String("server", "localhost:9314", "gRPC server address")
	interval  = flag.String("i", "10s", "interval for sync service")
	enableAll = flag.Bool("all", false, "enable for all service")
)

// see probe item: http://t.com:9313/status
func main() {

	log.Printf("starting... server: %v, sync interval: %v\n", *server, *interval)
	i, err := time.ParseDuration(*interval)
	if err != nil {
		log.Fatal("error interval time duration format", err)
	}

	check(*server)
	sync(*server, i)
}

func check(server string) {
	s := cloudprober.New(server)
	_, err := getOldTargets(s)
	if err != nil {
		log.Fatal("check err", err)
	}
}
