package apps

import (
	"fmt"
	"github.com/miekg/dns"
	"testing"
)

func TestNewResolver(t *testing.T) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("www.baidu.com"), dns.TypeA)
	r := NewResolver()
	dd, e := r.Lookup("udp", m)
	fmt.Println(e)
	fmt.Println(dd)

}
