package apps

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
	"log"
	"strings"
	"sync"
)

type ResolvError struct {
	qname, net  string
	nameservers []string
}

func (e ResolvError) Error() string {
	errmsg := fmt.Sprintf("%s resolv failed on %s (%s)", e.qname, strings.Join(e.nameservers, "; "), e.net)
	return errmsg
}

type RResponse struct {
	msg        *dns.Msg
	nameserver string
	rtt        time.Duration
}

type Resolver struct {
	timeout time.Duration
}

func NewResolver() *Resolver {

	r := &Resolver{
		timeout: 6 * time.Second,
	}
	return r
}

func (r Resolver) Lookup(net string, req *dns.Msg) (message *dns.Msg, err error) {

	c := &dns.Client{
		Net:          net,
		ReadTimeout:  r.timeout,
		WriteTimeout: r.timeout,
	}

	qname := req.Question[0].Name

	res := make(chan *RResponse, 1)
	var wg sync.WaitGroup
	L := func(nameserver string) {
		defer wg.Done()
		r, rtt, err := c.Exchange(req, nameserver)
		if err != nil {
			log.Printf("%s socket error on %s", qname, nameserver)
			log.Printf("error:%s", err.Error())
			return
		}
		// If SERVFAIL happen, should return immediately and try another upstream resolver.
		// However, other Error code like NXDOMAIN is an clear response stating
		// that it has been verified no such domain existas and ask other resolvers
		// would make no sense. See more about #20
		if r != nil && r.Rcode != dns.RcodeSuccess {
			log.Printf("%s failed to get an valid answer on %s", qname, nameserver)
			if r.Rcode == dns.RcodeServerFailure {
				return
			}
		}
		re := &RResponse{r, nameserver, rtt}
		select {
		case res <- re:
		default:
		}
	}

	ticker := time.NewTicker(time.Duration(200) * time.Millisecond)
	defer ticker.Stop()
	// Start lookup on each nameserver top-down, in every second
	nameservers := Config.Resolv.DNSServers()
	for _, nameserver := range nameservers {
		wg.Add(1)
		go L(nameserver)
		// but exit early, if we have an answer
		select {
		case re := <-res:
			log.Printf("%s resolv on %s rtt: %v", RemoveDomain(qname), re.nameserver, re.rtt)
			return re.msg, nil
		case <-ticker.C:
			continue
		}
	}
	// wait for all the namservers to finish
	wg.Wait()
	select {
	case re := <-res:
		log.Printf("%s resolv on %s rtt: %v", RemoveDomain(qname), re.nameserver, re.rtt)
		return re.msg, nil
	default:
		return nil, ResolvError{qname, net, nameservers}
	}

}
