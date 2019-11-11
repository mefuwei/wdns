package core

import (
	"github.com/golang/glog"
	"github.com/mefuwei/dns/storage"
	"github.com/miekg/dns"
	"net"
	"strconv"
)

const (
	resovePath = "/etc/resolv.conf"
)

var (
	defaultServers = []string{"114.114.114.114"}
	defaultPort    = 53

	// TODO used config
	storageType = "redis"
	redisAddr   = "localhost:6379"
	redisPasswd = ""
	redisDb     = 1
)

type DnsHandler struct{}

func (d *DnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	h := NewHandler(w, r)
	h.Do()
}

func NewHandler(w dns.ResponseWriter, r *dns.Msg) *Handler {
	question := r.Question[0]
	name, qtype := question.Name, question.Qtype
	remoteAddr := w.RemoteAddr().String()
	respMsg := new(dns.Msg)
	respMsg.SetReply(r)

	h := &Handler{
		Client:     new(dns.Client),
		W:          w,
		ReqMsg:     r,
		RespMsg:    respMsg,
		Name:       name,
		Qtype:      qtype,
		RemoteAddr: remoteAddr,
	}
	return h
}

type Handler struct {
	Client     *dns.Client        // exchange client
	W          dns.ResponseWriter // write at msg
	ReqMsg     *dns.Msg           // dns client request msg type
	RespMsg    *dns.Msg           // output to dns client msg
	Name       string             // dns question name
	Qtype      uint16             // dns question type
	RemoteAddr string
}

// TODO get backend sorage for common object and this config.
func (h *Handler) Do() {

	bs := storage.GetStorage(storageType, redisAddr, redisPasswd, redisDb)
	if msg, err := bs.Get(h.Name, h.Qtype); err != nil {
		// if not match local dns proxy to resolve
		h.Exchange()
		return
	} else {
		// backend storage return a null msg
		msg = msg.SetReply(h.RespMsg)
		h.Write(msg)
		return
	}
}

func (h *Handler) Exchange() {
	var Servers []string
	var Port string

	config, err := dns.ClientConfigFromFile(resovePath)
	if err != nil {
		glog.Errorf("Parse %s failed use default nameserver, %s", resovePath, err.Error())
		Servers = defaultServers
		Port = strconv.Itoa(defaultPort)
	} else {
		Servers = config.Servers
		Port = config.Port
	}

	// do exchange
	for _, srv := range Servers {
		server := net.JoinHostPort(srv, Port)
		if respMsg, _, err := h.Client.Exchange(h.ReqMsg, server); err == nil {
			h.Write(respMsg)
			return
		}
	}
}

// Write msg to client
func (h *Handler) Write(msg *dns.Msg) {
	if err := h.W.WriteMsg(msg); err != nil {
		glog.Errorf("[%s] fuck-dns write dns failed, %s", h.RemoteAddr, err.Error())
	} else {
		glog.Infof("[%s] query name: %s type: %d write success", h.RemoteAddr, h.Name, h.Qtype)
	}
}
