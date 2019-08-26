package apps

import (
	"net"
	"strconv"
	"time"

	"github.com/miekg/dns"
)

type Server struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// join addr and port
func (s *Server) Addr() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))

}

// server fun function
func (s *Server) Run() {
	handler := NewHandler()
	tcpHandler := dns.NewServeMux()

	udpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", handler.DoTCP)
	udpHandler.HandleFunc(".", handler.DoUDP)

	tcpServer := &dns.Server{
		Addr:         s.Addr(),
		Net:          "tcp",
		Handler:      tcpHandler,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
	}

	udpServer := &dns.Server{
		Addr:         s.Addr(),
		Net:          "udp",
		Handler:      udpHandler,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
	}
	go s.start(tcpServer)
	go s.start(udpServer)

}

func (s *Server) start(ds *dns.Server) {
	logger.Infof("start %s listen on %s ", ds.Net, ds.Addr)
	err := ds.ListenAndServe()
	if err != nil {
		logger.Panicf("start server error :%v", err)
	}

}
