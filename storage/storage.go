// storage is backend storage
// you can define more backend storage mysql ....

package storage

import (
	"github.com/miekg/dns"
	"strings"
)

type Storage interface {
	List() (msgs []*dns.Msg, err error)
	ApiList() (records []Record, err error)

	// Get a msg for backend storage, if use handler please you msg.SetReply(reqMsg)
	Get(name string, qtype uint16) (msg *dns.Msg, err error)
	ApiGet(name string, qtype uint16) (records []Record, err error)

	Set(records []Record) error
	ApiSet(records []Record) error

	// test backend storage connect
	Ping() bool
}

type Record struct {
	// domain json struct
	Rtype  uint16 `json:"rtype"`  // 记录类型 example:  dns.TYPEA
	Host   string `json:"host"`   // 主机记录 host www
	Domain string `json:"domain"` // 域名 qianbao-inc.com prefix.Domain = dns.Name
	Line   int    `json:"line"`   // 线路 实现智能DNS 开发环境/测试环境/预发环境/生产环境/联通/电信
	Value  string `json:"value"`  // A -> 8.8.8.8 CNAME -> www.qianbao.com.
	Ttl    uint32 `json:"ttl"`    // ttl
	Port   uint16 `json:"port"`   // SRV
}

func GetStorage(stype, Addr, Password string, db int) Storage {
	storageName := strings.ToLower(stype)

	switch storageName {
	case "redis":
		return NewRedisBackendStorage(Addr, Password, db)
	default:
		return NewRedisBackendStorage(Addr, Password, db)
	}
}