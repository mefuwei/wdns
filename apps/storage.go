package apps

import (
	"log"

	"time"

	"encoding/json"

	"github.com/go-redis/redis"
)

/*
TODO: Redis storage engine
*/

type JsonSerializer struct {
	ParseResult string `json:"parse_result"`
	ParseType   string `json:"parse_type"`
	Ttl         int    `json:"ttl"`
}

func Loads(data []byte) ([]*JsonSerializer, error) {
	var ret []*JsonSerializer
	err := json.Unmarshal(data, &ret)
	return ret, err
}

func Dumps(c []*JsonSerializer) (encoded []byte, err error) {
	encoded, err = json.Marshal(c)
	return

}

var redisEngine *RedisEngine

type RedisEngine struct {
	Backend *redis.Client
	// Serializer JsonSerializer
}

// add parse area
func (c *RedisEngine) SetArea(v ...string) error {
	_, err := c.Backend.SAdd(Config.Area, v).Result()
	if err != nil {
		return err
	}
	return nil
}

// get area  default_cname_baidu_com_a_60
func (c *RedisEngine) GetArea() ([]string, error) {
	s, err := c.Backend.SMembers(Config.Area).Result()
	if err != nil {
		return nil, err
	}
	return s, nil

}

// redis
func (c *RedisEngine) Set(key string, v ...string) error {

	if _, err := c.Backend.SAdd(key, v).Result(); err != nil {
		return err
	} else {
		return nil
	}

}

// get area
func (c *RedisEngine) Get(key string) ([]string, error) {
	result, err := c.Backend.SMembers(key).Result()
	if err != nil {
		log.Printf("redis error : %v : pharse Get %s ", err.Error(), key)
		return nil, err
	}

	return result, nil

}

func (c *RedisEngine) HSet(area, domain string, value *JsonSerializer) error {

	var data, err = c.HGet(area, domain)
	if ok := VerifyDomainConflict(value, data); !ok {
		return recordConflict{key: domain}
	}
	data = append(data, value)
	encoded, err := Dumps(data)
	if err == nil {
		c.Backend.HSet(area, domain, encoded)

	} else {

		logger.Info(err)
		return err
	}
	return nil
}

func (c *RedisEngine) HGet(area, domain string) ([]*JsonSerializer, error) {

	value, err := c.Backend.HGet(area, domain).Result()
	if err == nil {
		return Loads([]byte(value))
	}
	return nil, err
}

func (c *RedisEngine) Ping() bool {
	result := c.Backend.Ping().Err()
	if result != nil {
		log.Printf("try connenct redis failed , error : %v", result)
		return false
	}
	return true
}

func NewRedisEngine() *RedisEngine {

	return &RedisEngine{
		Backend: redis.NewClient(&redis.Options{
			Addr:         Config.Redis.Addr(),
			Password:     Config.Redis.Password,
			DB:           Config.Redis.DB,
			PoolSize:     10,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		}),
	}

}

func init() {
	redisEngine = NewRedisEngine()
	// redis  health
	go func() {
		for {
			ok := redisEngine.Ping()

			if Config.Redis.Enable != ok {
				Config.Redis.Enable = redisEngine.Ping()
				log.Printf("redis enable change %v ", ok)
			}
			time.Sleep(3 * time.Second)
		}

	}()

}
