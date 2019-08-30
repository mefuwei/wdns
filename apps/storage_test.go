package apps

import (
	"testing"
)

func TestRedisEngine_HSet(t *testing.T) {
	var data JsonSerializer
	data.Domain = "www.dsb.com"
	data.Area = 1
	data.Ttl = 600
	data.Prefix = "www"
	data.Rtype = "A"
	data.Value = "19.19.19.191"
	redisEngine.HSet("wdns:dsb.com", "www", &data)
	//redisEngine.Ping()
}
