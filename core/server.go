package core

import (
	"github.com/golang/glog"
	"github.com/miekg/dns"
	"net"
	"os"
	"strings"
	"time"
)

func NewServer(addr string) *Server {
	return &Server{
		addr:    addr,
		Handler: &DnsHandler{},
	}
}

type Server struct {
	addr string
	Handler *DnsHandler
}

func (s *Server) Start() {
	prots := []string{"udp", "tcp"}

	for _, p := range prots {
		netAddr := net.JoinHostPort(s.addr, p)
		srv := &dns.Server{Addr: s.addr, Net: p, Handler: s.Handler}
		glog.Infof("dns starting %s Server, listen %s...", strings.ToUpper(p), netAddr)
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				glog.Fatalf("dns starting %s Server, listen %s failed, %s", strings.ToUpper(p), netAddr, err.Error())
				os.Exit(1)
			}
		}()
	}

	for {
		time.Sleep(time.Second * 3)
	}
}