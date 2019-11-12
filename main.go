package main

import (
	"github.com/mefuwei/wdns/apps"
	"github.com/mefuwei/wdns/core"
	"os"
	"os/signal"
)

func main() {

	dnsServer := core.NewDnsServer(apps.Config.Server.Host)
	dnsServer.Start()
	webServer := core.NewWebServer("0.0.0.0:8989")
	webServer.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

forever:
	for {
		select {
		case <-sig:
			break forever
		}
	}

}
