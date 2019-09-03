// storage is backend storage
// you can define more backend storage mysql ....

package storage

import (
	"github.com/mefuwei/wdns/apps"
	"github.com/miekg/dns"
	"strings"
)

type Storage interface {
	// list view all dns
	List(viewName string) (msgs []*dns.Msg, err error)
	// use name qtype query backend storage dns.msg
	Get(name string, qtype uint16) (*dns.Msg, error)
	// use dns.msg parse name and qtype to write backend storage.
	Set(msg *dns.Msg) error
}

// json string parse to dns.msg
func StringToMsg(m string) (msg *dns.Msg, err error) {

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

