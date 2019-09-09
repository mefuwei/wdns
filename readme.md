#### 一个简单快速的dns缓存服务器，由go编写。
> 类似 dnsmasq,支持功能 智能区域解析 缓存 转发 ,本地解析支持类型 A记录 AAAA 记录 CNAME

main -> core.Server -> Handlers -> storage...

## 任务列表
- [x] 项目分离测试
- [] 实现restful功能
- [] 实现权威DNS的功能
- [] 高性能支持低TTL
- [] View 支持根据不同IP地址加载不同配置，默认为default
  
## 说明 
1. 安装
```bash
go get 

cd wdns

./build.sh
```

2 启动

```bash

sudo ./wdns -c etc/dns.conf

```

3. 配置文件

> 配置文件 dns.conf 是TOML 格式
详情参考  https://github.com/mojombo/toml

Example
```
#Toml config file


Version = "0.0.1"
Author = "F.W"

Debug = false

Area = "area"
[server]
host = "127.0.0.1"
port = 53


[redis]
enable = true
host = "127.0.0.1"
port = 6379
db = 0
password =""


[log]
stdout = true
file = "./wdns.log"
level = "INFO"
#DEBUG | INFO |NOTICE | WARN | ERROR

[cache]
backend = "memory"
expire = 600
# 10 minutes
maxCount = 0
# If set zero. The Sum of cache itmes will be unlimit.


[resolv]
nameserver = [
"192.168.20.100",
"114.114.115.115",
"114.114.114.114",
"208.67.220.220",
"119.29.29.29",
"180.76.76.76",
"223.6.6.6",
"223.5.5.5",
"8.8.8.8"
]
```

#### 缓存

> 默认是 memory 作为缓存,

```toml
[cache]
backend = "memory"
expire = 600
# 10 minutes
maxCount = 0
# If set zero. The Sum of cache itmes will be unlimit.


```

####  存储  
> 支持  Redis sqlete3 默认是 redis

#### 架构

> 请求查询顺序 缓存 [memory] --> 存储[redis|mysql] --> nameserver  


> DNS记录类型与编号中文版   https://zh.wikipedia.org/wiki/DNS%E8%AE%B0%E5%BD%95%E7%B1%BB%E5%9E%8B%E5%88%97%E8%A1%A8

> DNS 记录类型与编英文版 https://en.wikipedia.org/wiki/List_of_DNS_record_types