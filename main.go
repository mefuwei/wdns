package main

import (
	"flag"
	"github.com/mefuwei/dns/core"
)

var (
	addr string
)

func init() {
	flag.Parse()
	_ = flag.Set("stderrthreshold", "info")
	flag.StringVar(&addr, "dnsAddr", "0.0.0.0:53", "bind host, example 192.168.1.1:53")
	flag.StringVar(&addr, "webAddr", "0.0.0.0:8080", "bind host, example 192.168.1.1:8080")
}

func main() {
	dnsServer := core.NewDnsServer(addr)
	dnsServer.Start()

	webServer := core.NewWebServer(addr)
	webServer.Start()
}
