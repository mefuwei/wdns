package apps

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"net"

	"github.com/BurntSushi/toml"
)

type DNSServerSetting struct {
	Host string
	Port int
}

type RedisSettings struct {
	//Enable   bool
	Host     string
	Port     int
	DB       int
	Password string
}

func (s RedisSettings) Addr() string {

	return s.Host + ":" + strconv.Itoa(s.Port)

}

type CacheSettings struct {
	Backend  string
	Expire   int
	MaxCount int
}

type ResolvSetting struct {
	Nameserver []string
}

func (c *ResolvSetting) DNSServers() []string {

	var result []string
	for _, nameserver := range c.Nameserver {
		result = append(result, net.JoinHostPort(nameserver, "53"))

	}
	return result

}

type LogSettings struct {
	Stdout bool
	File   string
	Level  string
}

type Settings struct {
	DbType  string `toml:"DbType"`
	Area    string
	Version string
	Debug   bool
	Author  string
	Server  DNSServerSetting `toml:"server"`
	Redis   RedisSettings    `toml:"redis"`
	Log     LogSettings      `toml:"log"`
	Cache   CacheSettings    `toml:"cache"`
	Resolv  ResolvSetting    `toml:"resolv"`
}

var (
	Config Settings
)

func init() {
	var configFile string
	flag.StringVar(&configFile, "c", "/Users/fuwei/go/src/mefuwei/wdns/etc/dns.conf", "./wdns -c etc/dns.conf")
	flag.Parse()
	if _, err := toml.DecodeFile(configFile, &Config); err != nil {
		fmt.Printf("%s is valid toml config \n", configFile)
		fmt.Println(err)
		os.Exit(1)
	}
}
