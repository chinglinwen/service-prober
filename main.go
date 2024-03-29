package main

import (
	"flag"
	"log"
	"time"
	"wen/service-prober/pkg/cloudprober"
)

var (
	gprcServer = flag.String("server", "localhost:9314", "gRPC server address")
	interval   = flag.String("i", "1m", "interval for sync service")
	enableAll  = flag.Bool("all", false, "enable for all service")

	pathAnnotation   = flag.String("anPath", "prober.haodai.net/path", "service path annotation key")
	enableAnnotation = flag.String("anEnable", "prober.haodai.net/enable", "service prober enable annotation key")

	probeInterval = flag.Int("httpInterval", 3000, "http probe interval (milli-second) ")
	probeTimeout  = flag.Int("httpTimeout", 5000, "http probe timeout (milli-second)")
)

// see probe item: http://t.com:9313/status
func main() {
	log.Printf("starting... server: %v, sync interval: %v\n", *gprcServer, *interval)
	i, err := time.ParseDuration(*interval)
	if err != nil {
		log.Fatal("error interval time duration format", err)
	}

	check(*gprcServer)
	server := newserver(*gprcServer, i)
	server.sync()
}

func check(server string) {
	s := cloudprober.New(server)
	_, err := getOldTargets(s)
	if err != nil {
		log.Fatal("check err", err)
	}
}
