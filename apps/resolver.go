package apps

import (
	"fmt"
	"github.com/miekg/dns"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	dnsc *dns.Client
	dnsf *dns.ClientConfig
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

	dnsc = &dns.Client{
		Net:          net,
		ReadTimeout:  r.timeout,
		WriteTimeout: r.timeout,
	}

	qname := req.Question[0].Name

	res := make(chan *RResponse, 1)
	var wg sync.WaitGroup
	L := func(nameserver string) {
		defer wg.Done()
		r, rtt, err := dnsc.Exchange(req, nameserver)
		if err != nil {
			logger.Debugf("%s socket error on %s error: ", nameserver, err.Error())
			return
		}
		//if r == nil || r.Rcode == dns.RcodeNameError || r.Rcode == dns.RcodeSuccess{
		//if r != nil && r.Rcode != dns.RcodeSuccess {
		//		//	log.Printf("%s failed to get an valid answer on %s", qname, nameserver)
		//		//	if r.Rcode == dns.RcodeServerFailure {
		//		//		return
		//		//	}
		//		//}

		if r == nil || r.Rcode == dns.RcodeNameError || r.Rcode == dns.RcodeSuccess {

			re := &RResponse{r, nameserver, rtt}
			select {
			case res <- re:
			default:
			}
		} else {
			return
		}
	}

	ticker := time.NewTicker(time.Duration(200) * time.Millisecond)
	defer ticker.Stop()

	for _, server := range dnsf.Servers {
		wg.Add(1)
		go L(server + ":" + dnsf.Port)
		// but exit early, if we have an answer
		select {
		case re := <-res:

			return re.msg, nil
		case <-ticker.C:
			continue
		}
	}
	// wait for all the namserver to finish, catch not data return error
	wg.Wait()
	select {
	case re := <-res:
		return re.msg, nil
	default:
		return nil, ResolvError{qname, net, dnsf.Servers}
	}

}

func init() {
	var err error
	dnsf, err = dns.ClientConfigFromFile("/etc/resolv.conf")
	fmt.Println(dnsf)
	if err != nil || dnsf == nil {
		fmt.Printf("Cannot initialize the local resolver: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println("load nameserver success")
	}

}
