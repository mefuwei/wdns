package core

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/mefuwei/dns/apis"
	"github.com/miekg/dns"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
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

func NewWebServer(addr string) *WebServer {
	return &WebServer{
		addr: addr,
	}
}

type WebServer struct {
	addr    string
}

func (w *WebServer) Start() {
	service := new(restful.WebService)
	service.Path("/api/v1").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	service.Route(service.POST("dns").To(apis.DnsAdd))

	restful.Add(service)
	go glog.Fatal(http.ListenAndServe(w.addr, nil))
}