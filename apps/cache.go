package apps

import (
	"sync"
	"time"

	"crypto/md5"
	"encoding/hex"
	"github.com/miekg/dns"
)

/*
TODO: cache interface
*/

type Cache interface {
	Get(key string) (Msg *dns.Msg, err error)
	Set(key string, Msg *dns.Msg) bool
	Remove(key string) error
}

type Mesg struct {
	Msg    *dns.Msg
	Expire time.Time
}

/*
TODO: memory cache Backend
*/

type MemoryCache struct {
	Backend  map[string]Mesg
	Expire   time.Duration
	MaxCount int
	mu       sync.RWMutex
}

// memory get function
func (c *MemoryCache) Get(key string) (*dns.Msg, error) {
	c.mu.RLock()
	mesg, ok := c.Backend[key]
	c.mu.RUnlock()
	if !ok {
		return nil, KeyNotFound{key}
	}
	if mesg.Expire.Before(time.Now()) {
		c.Remove(key)
		return nil, KeyExpired{key}
	}
	return mesg.Msg, nil
}

// memory set function
func (c *MemoryCache) Set(key string, msg *dns.Msg) bool {

	expire := time.Now().Add(c.Expire)
	mesg := Mesg{msg, expire}
	c.mu.Lock()
	c.Backend[key] = mesg
	c.mu.Unlock()
	return true
}

// memory remove
func (c *MemoryCache) Remove(key string) error {
	c.mu.Lock()
	delete(c.Backend, key)
	c.mu.Unlock()
	return nil
}

// if cache size,0 no limit
func (c *MemoryCache) IsFull() bool {
	if c.MaxCount == 0 {
		return false
	}
	return false
}

// create md5 key
func KeyGen(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))

}
