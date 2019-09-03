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
	// domain json struct
	Rtype  string        `json:"rtype"`
	Prefix string        `json:"prefix"`
	Domain string        `json:"domain"`
	Area   int           `json:"area"`
	Value  string        `json:"value"`
	Ttl    time.Duration `json:"ttl"`
}

type RedisSerializer []*JsonSerializer

func (c *RedisSerializer) Loads(str string) error {
	err := json.Unmarshal([]byte(str), &c)
	return err
}

func (c *RedisSerializer) Dumps() (str string, err error) {
	encoded, err := json.Marshal(c)

	return string(encoded), err

}

var redisEngine *RedisEngine

type RedisEngine struct {
	Backend *redis.Client
	// Serializer JsonSerializer
}

func (c *RedisEngine) HSet(domain, prefix string, data *JsonSerializer) error {
	var err error
	dataList := RedisSerializer{}

	dataList, err = c.HGet(domain, prefix)

	err = VerifyRecordRules(data, dataList)
	if err != nil {
		logger.Infoln(err)
		return err
	}
	dataList = append(dataList, data)
	jsonStr, err := dataList.Dumps()
	if err == nil {
		c.Backend.HSet(domain, prefix, jsonStr)

	} else {

		logger.Info(err)
		return err
	}
	return nil
}

func (c *RedisEngine) HGet(domain, prefix string) (RedisSerializer, error) {
	dataList := RedisSerializer{}

	str, err := c.Backend.HGet(domain, prefix).Result()
	if err == nil {

		return dataList, dataList.Loads(str)
	}
	return nil, err
}

func (c *RedisEngine) Exist(prefix string) (bool, error) {
	if res, err := c.Backend.Keys(prefix).Result(); err != nil {
		return false, err
	} else {
		if res == nil {
			return false, nil // not exist
		}
		return true, nil
	}
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
		switch Config.DbType {
		case "REDIS":
			for {

				ok := redisEngine.Ping()
				if !ok {
					logger.Infof("the redis server is down ")
				}

				time.Sleep(3 * time.Second)
			}
		default:

			logger.Infof("break check ")
		}

	}()
}
