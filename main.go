package main

import (
	"mefuwei/wdns/apps"
	"os"
	"os/signal"
	"time"
)

func main() {

	svr := &apps.Server{
		Host:         apps.Config.Server.Host,
		Port:         apps.Config.Server.Port,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	svr.Run()

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
