package main

import (
	"flag"
	"github.com/mefuwei/dns/core"
)

var (
	addr string
)

func init()  {
	flag.Parse()
	_ = flag.Set("stderrthreshold", "info")
	flag.StringVar(&addr, "addr", "0.0.0.0:53", "bind host, example 192.168.1.1:53")
}

func main() {
	server := core.NewServer(addr)
	server.Start()
}
