package discovery

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type Discovery struct {
	conn redis.Conn
}

func NewDiscovery(redisAddr string) *Discovery {
	c, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return nil
	}
	d := &Discovery{conn: c}
	return d
}

func (d *Discovery) GetCometAddr() []string {
	num, err := d.conn.Do("llen", "cometlist")
	if err != nil {
		fmt.Println("Cometlist get len err", err)
	}
	var res []string
	val, _ := redis.Values(d.conn.Do("lrange", "cometlist", "0", num))
	for _, v := range val {
		res = append(res, string(v.([]byte)))
	}
	return res
}

func (d *Discovery) RegComet(addr string) {
	_, err := d.conn.Do("lpush", "cometlist", addr)
	if err != nil {
		fmt.Println("Register comet failed")
	}
}

func (d *Discovery) DelComet(addr string) {
	_, err := d.conn.Do("lrem", "cometlist", 1, addr)
	if err != nil {
		fmt.Println("Del comet failed", err)
	}
}
