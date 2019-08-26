package apps

import (
	"fmt"
	"testing"
)

func TestRedisEngine_connect(t *testing.T) {
	if s, err := redisEngine.GetArea(); err != nil {
		fmt.Println(err)
	} else {
		for _, v := range s {
			fmt.Println(v)
		}
	}
}

func TestRedisEngine_HSet(t *testing.T) {
	var data JsonSerializer
	data.ParseResult = "10.10.10.10"
	data.ParseType = "A"
	data.Ttl = 600
	redisEngine.HSet("wdns:bj-sh", "www.test.com", &data)
	redisEngine.Ping()
}
