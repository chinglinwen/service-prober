package main

import (
	"flag"
	"log"
	"time"
	"wen/service-prober/pkg/cloudprober"

	"github.com/pkg/profile"
)

var (
	server    = flag.String("server", "localhost:9314", "gRPC server address")
	interval  = flag.String("i", "1m", "interval for sync service")
	enableAll = flag.Bool("all", false, "enable for all service")

	pathAnnotation   = flag.String("anPath", "prober.haodai.net/path", "service path annotation key")
	enableAnnotation = flag.String("anEnable", "prober.haodai.net/enable", "service prober enable annotation key")

	probeInterval = flag.Int("httpInterval", 3000, "http probe interval (milli-second) ")
	probeTimeout  = flag.Int("httpTimeout", 5000, "http probe timeout (milli-second)")
)

// see probe item: http://t.com:9313/status
func main() {
	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
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
