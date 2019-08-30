package apps

import (
	"net"

	"log"
	"time"

	"fmt"
	"github.com/miekg/dns"
)

const (
	notIPQuery = 0
	_IP4Query  = 1
	_IP6Query  = 28
	_IPCname   = 5
)

type Question struct {
	qname  string
	qtype  string
	qclass string
}

func (q *Question) String() string {
	return q.qname + " " + q.qclass + " " + q.qtype
}

func NewHandler() *GODNSHandler {
	var cache, memoryCache Cache
	memoryCache = &MemoryCache{
		Backend:  make(map[string]Mesg, Config.Cache.MaxCount),
		Expire:   time.Duration(Config.Cache.Expire) * time.Second,
		MaxCount: Config.Cache.MaxCount,
	}
	return &GODNSHandler{cache, memoryCache, NewResolver(), redisEngine}

}

type GODNSHandler struct {
	cache, memoryCache Cache
	resolver           *Resolver
	db                 *RedisEngine
}

func (h *GODNSHandler) do(Net string, w dns.ResponseWriter, req *dns.Msg) {

	q := req.Question[0]
	Q := Question{RemoveDomain(q.Name), dns.TypeToString[q.Qtype], dns.ClassToString[q.Qclass]}

	var remote net.IP

	if Net == "tcp" {
		remote = w.RemoteAddr().(*net.TCPAddr).IP
	} else {
		remote = w.RemoteAddr().(*net.UDPAddr).IP
	}
	logger.Infof("remote ip %s", remote)
	queryMemoryCache := h.isQueryCache(q)
	memoryMapKey := KeyGen(Q.String())

	//  the type  A or AAAA   query  memory cache
	if queryMemoryCache > 0 {
		mesg, err := h.memoryCache.Get(memoryMapKey)
		if err != nil {
			logger.Infof("%s din't hit memory cache", Q.String())
		} else {
			logger.Infof("the domain  %s hit memory cache", Q.String())

			msg := *mesg
			msg.Id = req.Id
			w.WriteMsg(&msg)
			return
		}

	}
	switch Config.DbType {
	case "REDIS":
		logger.Infof("find %v from redis engine, phase : start ", Q.String())

		//ips, err := h.db.HGet("www.baidu.com","wwww")
		//if err == nil && len(ips) > 0 {
		//	m := new(dns.Msg)
		//	m.SetReply(req)
		//	rr_header := dns.RR_Header{
		//		Name:   q.Name,
		//		Rrtype: dns.TypeA,
		//		Class:  dns.ClassINET,
		//		Ttl:    600,
		//	}
		//	var ip net.IP
		//	for _, v := range ips {
		//		ip = net.ParseIP(v).To4()
		//		m.Answer = append(m.Answer, &dns.A{rr_header, ip})
		//
		//	}
		//
		//	w.WriteMsg(m)
		//	logger.Infof("find %v from db engine, phase : end ", Q.String())
		//	return
		//
		//}

	default:
		logger.Infof("other ")

	}

	log.Printf("find %v from resolver , phase : start ", Q.String())

	mesg, err := h.resolver.Lookup(Net, req)
	if err != nil {
		fmt.Println("ssf")
		dns.HandleFailed(w, req)

		// cache the failure, too!
	} else {
		h.memoryCache.Set(memoryMapKey, mesg)

		w.WriteMsg(mesg)
	}
	// set query to memory

}

func (h *GODNSHandler) DoTCP(w dns.ResponseWriter, req *dns.Msg) {
	h.do("tcp", w, req)
}

func (h *GODNSHandler) DoUDP(w dns.ResponseWriter, req *dns.Msg) {
	h.do("udp", w, req)
}

// Only query cache when dns qclass == 'IN' and qtype == 'A'|'AAAA'
func (h *GODNSHandler) isQueryCache(q dns.Question) int {

	if q.Qclass != dns.ClassINET {
		return notIPQuery
	}
	switch q.Qtype {
	case dns.TypeA:
		return _IP4Query
	case dns.TypeAAAA:
		return _IP6Query
	case dns.TypeCNAME:
		return _IPCname

	default:
		return notIPQuery

	}
}

// domain name is fully qualified and  remove domain name .
func RemoveDomain(s string) string {
	if dns.IsFqdn(s) {
		return s[:len(s)-1]
	}
	return s
}
