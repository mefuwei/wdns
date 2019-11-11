package core

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/mefuwei/wdns/apis"
	"net/http"
	"os"
)

func NewWebServer(addr string) *WebServer {
	return &WebServer{
		addr: addr,
	}
}

type WebServer struct {
	addr string
}

func (w *WebServer) Start() {
	service := new(restful.WebService)
	service.Path("/api/v1").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	service.Route(service.POST("dns").To(apis.DnsAdd))

	restful.Filter(beforeFilers)
	restful.Add(service)

	go func() {
		glog.Infof("start web server of %s", w.addr)
		if err := http.ListenAndServe(w.addr, nil); err != nil {
			glog.Infof("!!! web server startup failed. %s", err.Error())
			os.Exit(2)
		}
	}()
}

func beforeFilers(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	glog.Infof("[%s] %s", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}
