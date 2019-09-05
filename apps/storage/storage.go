// storage is backend storage
// you can define more backend storage mysql ....

package storage

import (
	"github.com/mefuwei/wdns/apps"
	"github.com/miekg/dns"
	"net"
	"strings"
)

type Storage interface {
	// list view all dns
	List(viewName string) (msgs []*dns.Msg, err error)
	// use name qtype query backend storage dns.msg
	Get(viewName, name string, qtype uint16) (msg *dns.Msg, err error)
	// use dns.msg parse name and qtype to write backend storage.
	Set(msg *dns.Msg) error
}

type Record struct {
	// domain json struct
	Rtype  uint16        `json:"rtype"` // 记录类型 example:  dns.TYPEA
	Host string        `json:"host"` // 主机记录 host www
	Domain string        `json:"domain"` // 域名 qianbao-inc.com prefix.Domain = dns.Name
	Line   int           `json:"line"` // 线路 实现智能DNS 开发环境/测试环境/预发环境/生产环境/联通/电信
	Value  string        `json:"value"` // A -> 8.8.8.8 CNAME -> www.qianbao.com.
	Ttl    uint32 `json:"ttl"` // ttl
	Port uint16 `json:"port"` // SRV
}

func GetStorage() *Storage {
	storageName := strings.ToLower(apps.Config.DbType)

	switch storageName {
	case "redis":
		return NewRedisBackendStorage()
	default:
		return NewRedisBackendStorage()
	}
}

func SwitchMsg(records []Record) (msg *dns.Msg, err error) {
	for _, record := range records {
		switch record.Rtype {
		case dns.TypeA:
			msg = switchA([]Record)
		case dns.TypeAAAA:
			msg = switchAAAA([]Record)
		case dns.TypeCNAME:
			msg = switchCNAME([]Record)
		case dns.TypeMX:
			msg = switchMX([]Record)
		case dns.TypeTXT:
			msg = switchTXT([]Record)
		case dns.TypeNS:
			msg = switchNS([]Record)
		case dns.TypeSRV:
			msg = switchSRV([]Record)
		default:
			// dns.TypeA
			msg = switchA([]Record)
		}
		break
	}
	return
}

func switchA(records []Record) (msg *dns.Msg, err error) {
	for _, record := range records {
		ip := net.ParseIP(record.Value)
		name := parseName(record.Host, record.Domain)

		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{
				Name: name,
				Rrtype: dns.TypeA,
				Class: dns.ClassINET,
				Ttl: record.Ttl,
			},
			A:   ip,
		})
	}
	return
}

func switchAAAA(records []Record) (msg *dns.Msg, err error) {
	for _, record := range records {
		ip := net.ParseIP(record.Value)
		name := parseName(record.Host, record.Domain)

		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{
				Name: name,
				Rrtype: dns.TypeAAAA,
				Class: dns.ClassINET,
				Ttl: record.Ttl,
			},
			A:   ip,
		})
	}
	return
}

func switchCNAME(records []Record) (msg *dns.Msg, err error) {
	for _, record := range records {
		name := parseName(record.Host, record.Domain)

		msg.Answer = append(msg.Answer, &dns.CNAME{
			Hdr:    dns.RR_Header{
				Name: name,
				Rrtype: dns.TypeCNAME,
				Class: dns.ClassINET,
				Ttl: record.Ttl,
			},
			Target: record.Value,
		})
	}
	return
}

func switchMX(records []Record) (msg *dns.Msg, err error) {
	for _, record := range records {
		name := parseName()

		msg.Answer = append(msg.Answer, &dns.MX{
			Hdr:        dns.RR_Header{
				Name: name,
				Rrtype: dns.TypeMX,
				Class: dns.ClassINET,
				Ttl: record.Ttl,
			},
			Preference: 0,
			Mx:         record.Value,
		})
	}
	return
}

func switchNS(records []Record) (msg *dns.Msg, err error) {
	for _, record := range records {
		name := parseName()

		msg.Answer = append(msg.Answer, &dns.NS{
			Hdr: dns.RR_Header{
				Name: name,
				Rrtype: dns.TypeNS,
				Class: dns.ClassINET,
				Ttl: record.Ttl,
			},
			Ns:  record.Value,
		})
	}
	return
}

func switchTXT(records []Record) (msg *dns.Msg, err error) {
	for _, record := range records {
		name := parseName()

		msg.Answer = append(msg.Answer, &dns.TXT{
			Hdr: dns.RR_Header{
				Name: name,
				Rrtype: dns.TypeTXT,
				Class: dns.ClassINET,
				Ttl: record.Ttl,
			},
			Txt: []string{record.Value},
		})
	}
	return
}

func switchSRV(records []Record) (msg *dns.Msg, err error) {
	for _, record := range records {
		name := parseName()

		msg.Answer = append(msg.Answer, &dns.SRV{
			Hdr:      dns.RR_Header{
				Name: name,
				Rrtype: dns.TypeSRV,
				Class: dns.ClassINET,
				Ttl: record.Ttl,
			},
			Priority: 0,
			Weight:   0,
			Port:     record.Port,
			Target:   record.Value,
		})
	}
	return
}

// TODO sup PTR type
//func switchPTR(records []Record) (msg *dns.Msg, err error) {
//	for _, record := range records {
//		name := parseName()
//
//		msg.Answer = append(msg.Answer, &dns.PTR{
//			Hdr: dns.RR_Header{
//				Name: name,
//				Rrtype: dns.TypeSRV,
//				Class: dns.ClassINET,
//				Ttl: record.Ttl,
//			},
//			Ptr: "",
//		})
//	}
//	return
//}

func parseName(host, domain string) (name string) {
	name = host + "." + domain
	if strings.HasSuffix(name, ".") {
		name += "."
	}
	return name
}