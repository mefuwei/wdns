package storage

import (
	"github.com/miekg/dns"
	"testing"
)

type TR struct {
	Name string
	Qtype uint16
}

func getRedis() Storage {
	redis := GetStorage("redis", "localhost:6379", "", 1)
	return redis
}

func TestRedisBackendStorage_Ping(t *testing.T) {
	redis := getRedis()
	if ok := redis.Ping(); !ok {
		t.Errorf("redis connect failed, please you cheack you config file and redis application")
	}
}

func TestRedisBackendStorage_Set(t *testing.T) {
	record := []Record{
		{
			Rtype: dns.TypeA,
			Host: "www",
			Domain: "qianbao-inc.com.",
			Line: 1,
			Value: "8.8.8.8",
			Ttl: 30,
		},
		{
			Rtype: dns.TypeA,
			Host: "www",
			Domain: "qianbao-inc.com.",
			Line: 1,
			Value: "9.9.9.9",
			Ttl: 30,
		},
	}

	redis := getRedis()
	err := redis.Set(record)
	if err != nil {
		t.Error(err.Error())
	}

}

func TestRedisBackendStorage_Get(t *testing.T) {
	redis := getRedis()

	cases := []TR{
		{"www.qianbao-inc.com.", 1},
	}

	for _, r := range cases {
		_ , err := redis.Get(r.Name, r.Qtype)
		if err != nil {
			t.Errorf("redis get name: %s type: %d failed", r.Name, r.Qtype)
		}
	}
}

//func TestRedisBackendStorage_List(t *testing.T) {
//	redis := getRedis()
//
//	if msgs, err := redis.List(); err != nil {
//		t.Error(err.Error())
//	} else {
//		t.Logf("%#v", msgs)
//	}
//}