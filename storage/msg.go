package storage

import (
	"github.com/golang/glog"
	"github.com/miekg/dns"
	"net"
	"strings"
)

func SwitchMsg(records []Record) (msg *dns.Msg, err error) {
	for _, record := range records {
		switch record.Rtype {
		case dns.TypeA:
			msg, err = switchA(records)
		case dns.TypeAAAA:
			msg, err = switchAAAA(records)
		case dns.TypeCNAME:
			msg, err = switchCNAME(records)
		case dns.TypeMX:
			msg, err = switchMX(records)
		case dns.TypeTXT:
			msg, err = switchTXT(records)
		case dns.TypeNS:
			msg, err = switchNS(records)
		case dns.TypeSRV:
			msg, err = switchSRV(records)
		default:
			// dns.TypeA
			msg, err = switchA(records)
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func switchA(records []Record) (msg *dns.Msg, err error) {
	respMsg := new(dns.Msg)

	for _, record := range records {
		ip := net.ParseIP(record.Value)
		name := ParseName(record.Host, record.Domain)

		respMsg.Answer = append(respMsg.Answer, &dns.A{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    record.Ttl,
			},
			A: ip,
		})
	}
	return respMsg, nil
}

func switchAAAA(records []Record) (msg *dns.Msg, err error) {
	respMsg := new(dns.Msg)

	for _, record := range records {
		ip := net.ParseIP(record.Value)
		name := ParseName(record.Host, record.Domain)

		respMsg.Answer = append(respMsg.Answer, &dns.A{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    record.Ttl,
			},
			A: ip,
		})
	}
	return respMsg, nil
}

func switchCNAME(records []Record) (msg *dns.Msg, err error) {
	respMsg := new(dns.Msg)

	for _, record := range records {
		name := ParseName(record.Host, record.Domain)

		respMsg.Answer = append(respMsg.Answer, &dns.CNAME{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: dns.TypeCNAME,
				Class:  dns.ClassINET,
				Ttl:    record.Ttl,
			},
			Target: record.Value,
		})
	}
	return respMsg, nil
}

func switchMX(records []Record) (msg *dns.Msg, err error) {
	respMsg := new(dns.Msg)

	for _, record := range records {
		name := ParseName(record.Host, record.Domain)

		respMsg.Answer = append(respMsg.Answer, &dns.MX{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: dns.TypeMX,
				Class:  dns.ClassINET,
				Ttl:    record.Ttl,
			},
			Preference: 0,
			Mx:         record.Value,
		})
	}
	return respMsg, nil
}

func switchNS(records []Record) (msg *dns.Msg, err error) {
	respMsg := new(dns.Msg)

	for _, record := range records {
		name := ParseName(record.Host, record.Domain)

		respMsg.Answer = append(respMsg.Answer, &dns.NS{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: dns.TypeNS,
				Class:  dns.ClassINET,
				Ttl:    record.Ttl,
			},
			Ns: record.Value,
		})
	}
	return respMsg, nil
}

func switchTXT(records []Record) (msg *dns.Msg, err error) {
	respMsg := new(dns.Msg)

	for _, record := range records {
		name := ParseName(record.Host, record.Domain)

		respMsg.Answer = append(respMsg.Answer, &dns.TXT{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassINET,
				Ttl:    record.Ttl,
			},
			Txt: []string{record.Value},
		})
	}
	return respMsg, nil
}

func switchSRV(records []Record) (msg *dns.Msg, err error) {
	respMsg := new(dns.Msg)

	for _, record := range records {
		name := ParseName(record.Host, record.Domain)

		respMsg.Answer = append(respMsg.Answer, &dns.SRV{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: dns.TypeSRV,
				Class:  dns.ClassINET,
				Ttl:    record.Ttl,
			},
			Priority: 0,
			Weight:   0,
			Port:     record.Port,
			Target:   record.Value,
		})
	}
	return respMsg, nil
}

// TODO sup PTR type
//func switchPTR(records []Record) (msg *dns.Msg, err error) {
//	for _, record := range records {
//		name := ParseName(record.Host, record.Domain)
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

func ParseName(host, domain string) (name string) {
	name = host + "." + domain
	if !strings.HasSuffix(name, ".") {
		name += "."
	}
	glog.Infof("name: %s", name)

	return name
}
