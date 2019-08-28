package main

import (
	"flag"
	"log"
	"time"
)

var (
	server   = flag.String("server", "localhost:9314", "gRPC server address")
	interval = flag.String("i", "10s", "interval for sync service")
)

// see probe item: http://t.com:9313/status
func main() {
	flag.Parse()
	i, err := time.ParseDuration(*interval)
	if err != nil {
		log.Fatal("error interval time duration format", err)
	}
	sync(*server, i)
}
