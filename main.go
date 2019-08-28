package main

import (
	"flag"
)

var (
	server   = flag.String("server", "localhost:9314", "gRPC server address")
	addProbe = flag.String("add_probe", "", "Path to probe config to add")
	rmProbe  = flag.String("rm_probe", "", "Probe name to remove")
)

// see probe item: http://t.com:9313/status
func main() {
	flag.Parse()

	sync(*server)
}
