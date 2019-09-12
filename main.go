package main

import (
	"flag"
	"github.com/mefuwei/dns/core"
	"time"
)

var (
	dnsAddr string
	webAddr string
)

func init() {
	flag.Parse()
	_ = flag.Set("stderrthreshold", "info")
	flag.StringVar(&dnsAddr, "dnsAddr", "0.0.0.0:53", "bind host, example 192.168.1.1:53")
	flag.StringVar(&webAddr, "webAddr", "0.0.0.0:8080", "bind host, example 192.168.1.1:8080")
}

func main() {
	dnsServer := core.NewDnsServer(dnsAddr)
	dnsServer.Start()

	webServer := core.NewWebServer(webAddr)
	webServer.Start()

	for {
		time.Sleep(time.Second * 3)
	}
}
