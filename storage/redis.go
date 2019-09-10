package storage

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"github.com/miekg/dns"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	RedisEmpty = "redis: nil"
)

var (
	dnsMsgKey = "dns:%s:%d" // dns:{domainName}:{qtype} dns:www.qianbao-inc.com:1
	dnsPrefixKey = "dns:*"

	// falied message
	RedisGetFailed = "redis backend storage get key: %s name: %s type: %d failed, %s"
	RedisSetFailed = "redis backend storage set key: %s failed, %s"
	JsonParseFailed = "redis backend storage json parse msg failed name: %s type %d, %s"
	ParseDnsMsgFailed = "redis backend storage parse dns.msg failed name: %s type: %d, %s"
)

func NewRedisBackendStorage(Addr, Password string, db int) *RedisBackendStorage {
	rbs := &RedisBackendStorage{
		Client: redis.NewClient(&redis.Options{
			Addr:         Addr,
			Password:     Password,
			DB:           db,
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
	glog.Infof("keys: %#v", keys)

	return rbs.ParseMsg(keys)
}

// TODO redis.Keys is a bug, not patten dns:* need track this bug.
func (rbs *RedisBackendStorage) Get(name string, qtype uint16) (msg *dns.Msg, err error) {
	records, err := rbs.get(name, qtype)
	if err != nil {
		return msg, err
	}

	msg, err = SwitchMsg(records)
	glog.Infof("records: %#v", records)
	if err != nil {
		glog.Error(ParseDnsMsgFailed, name, qtype, err.Error())
		return
	}
	return msg, nil
}

func (rbs *RedisBackendStorage) get(name string, qtype uint16) (records []Record, err error) {
	if !strings.HasSuffix(name, ".") {
		name += "."
	}

	key := fmt.Sprintf(dnsMsgKey, name, qtype)
	res, err := rbs.Client.Get(key).Result()

	// TODO go-redis not ensure record is None or Error.
	if err != nil {
		// 这段逻辑需要注意，这里的err并不为空，只是为了不输出日志
		if err.Error() == RedisEmpty {
			return records, err // not fund record
		}
		// redis failed.
		glog.Errorf(RedisGetFailed, key, name, qtype, err.Error())
		return records, err
	}

	if err := json.Unmarshal([]byte(res), &records); err != nil {
		glog.Errorf(JsonParseFailed, name, qtype, err.Error())
		return records, err
	}
	return records, nil

}

func (rbs *RedisBackendStorage) Set(records []Record) error {
	var name string
	var qtype uint16
	var key string

	for _, record := range records {
		name = ParseName(record.Host, record.Domain)
		qtype = record.Rtype
		break
	}
	key = fmt.Sprintf(dnsMsgKey, name, qtype)
	return rbs.set(key, records)
}

func (rbs *RedisBackendStorage) set(key string, records []Record) error {

	data, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf(RedisSetFailed, key, err.Error())
	}

	if _, err := rbs.Client.Set(key, data, time.Microsecond * 0).Result(); err != nil {
		return err
	}
	return nil
}

func (rbs *RedisBackendStorage) Keys() (keys []string, err error) {
	key := fmt.Sprintf(dnsPrefixKey)

	// 这么写是为了好扩展
	// TODO bug, keys not fund anything.
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

// API
// Todo list rewrite
func (rbs *RedisBackendStorage) ApiList() (records []Record, err error)  {

}

func (rbs *RedisBackendStorage) ApiGet(name string, qtype uint16) (records []Record, err error) {
	records, err = rbs.get(name, qtype)
	if err != nil {
		return records, err
	}
	return records, nil
}

func (rbs *RedisBackendStorage) ApiSet(records []Record) error {
	return rbs.Set(records)
}

