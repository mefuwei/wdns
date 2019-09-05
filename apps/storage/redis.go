package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/go-redis/redis"
	"github.com/mefuwei/wdns/apps"
	"github.com/miekg/dns"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	dnsMsgKey = "dns:%s:%s:%d" // dns:{domainName}:{qtype} dns:www.qianbao-inc.com:1
	dnsPrefixKey = "dns:*"
	redisBackendStorage *RedisBackendStorage

	// error
	RedisGetFailed = "redis backend storage get name: %s type: %d failed, %s"
	JsonParseFailed = "redis backend storage json parse msg failed name: %s type %d, %s"
	ParseDnsMsgFailed = "redis backend storage parse dns.msg failed name: %s type: %d, %s"
)

func init() {
	redisBackendStorage = NewRedisBackendStorage()
	// redis  health
	go func() {
		switch apps.Config.DbType {
		case "REDIS":
			for {
				ok := redisBackendStorage.Ping()
				if !ok {
					glog.Infof("the redis server is down ")
					redisBackendStorage = NewRedisBackendStorage()
				}
				time.Sleep(3 * time.Second)
			}
		default:
			glog.Infof("break check ")
		}
	}()
}

func NewRedisBackendStorage() *RedisBackendStorage {
	rbs := &RedisBackendStorage{
		Client: redis.NewClient(&redis.Options{
			Addr:         apps.Config.Redis.Addr(),
			Password:     apps.Config.Redis.Password,
			DB:           apps.Config.Redis.DB,
			PoolSize:     10,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		}),
	}
	return rbs
}

type RedisBackendStorage struct {
	Client *redis.Client
}

func (rbs *RedisBackendStorage) List() (msgs []*dns.Msg, err error) {
	keys, err := rbs.Keys()
	if err != nil {
		return msgs, err
	}

	return rbs.ParseMsg(keys)
}

func (rbs *RedisBackendStorage) Get(name string, qtype uint16) (msg *dns.Msg, err error) {
	key := fmt.Sprintf(dnsMsgKey, name, qtype)

	res, err := rbs.Client.Get(key).Result()
	if err != nil {
		glog.Errorf(RedisGetFailed, name, qtype, err.Error())
		return
	}

	var records []Record
	if err := json.Unmarshal([]byte(res), &records); err != nil {
		glog.Errorf(JsonParseFailed, name, qtype, err.Error())
		return
	}

	msg, err = SwitchMsg(records)
	if err != nil {
		glog.Error(ParseDnsMsgFailed, name, qtype, err.Error())
		return
	}
	return msg, nil
}

func (rbs *RedisBackendStorage) Set(msg *dns.Msg) error {

}

func (rbs *RedisBackendStorage) Keys() (keys []string, err error) {
	key := fmt.Sprintf(dnsPrefixKey)

	// 这么写是为了好扩展
	if keys, err := rbs.Client.Keys(key).Result(); err != nil {
		return keys, err
	}
	return keys, err
}

func (rbs *RedisBackendStorage) ParseMsg(keys []string) (msgs []*dns.Msg, err error) {
	for _, k := range keys {
		flag := strings.Split(k, ":")
		name := flag[2]
		qt, _ := strconv.Atoi(flag[3])
		qtype := uint16(qt)
		if msg, err := rbs.Get(name, qtype); err != nil {
			glog.Errorf("reids backend dont get name: %s type: %d, %s", name, qtype, err)
			continue
		} else {
			msgs = append(msgs, msg)
		}
	}
	return msgs, nil
}

func (rbs *RedisBackendStorage) Ping() bool {
	result := rbs.Client.Ping().Err()
	if result != nil {
		log.Printf("try connenct redis failed , error : %v", result)
		return false
	}
	return true
}

