package core

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/miekg/dns"
	"net"
	"os"
	"strings"
)

func NewDnsServer(addr string) *DnsServer {
	return &DnsServer{
		addr:    addr,
		Handler: &DnsHandler{},
	}
}

type DnsServer struct {
	addr    string
	Handler *DnsHandler
}

func (s *DnsServer) Start() {
	schemes := []string{"udp", "tcp"}

	for _, p := range schemes {
		netAddr := net.JoinHostPort(s.addr, p)

		srv := &dns.Server{Addr: s.addr, Net: p, Handler: s.Handler}
		fmt.Printf("dns starting %s Server, listen %s...\n", strings.ToUpper(p), s.addr)
		glog.Infof("dns starting %s Server, listen %s...", strings.ToUpper(p), s.addr)
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				glog.Fatalf("dns starting %s Server, listen %s failed, %s", strings.ToUpper(p), netAddr, err.Error())
				os.Exit(1)
			}
		}()
	}
}
