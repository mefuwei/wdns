package main

import (
	"fmt"
	"testing"

	"github.com/miekg/dns"
)

const (
	nameserver = "114.114.114.119:53"
	domain     = "www.qq.com"
)

func BenchmarkDig(b *testing.B) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	// var r *dns.Msg
	c := new(dns.Client)

	a := "d"
	fmt.Print(a)
	c.Exchange(m, nameserver)
	//fmt.Printf("-- %v -- \n",m)

}
