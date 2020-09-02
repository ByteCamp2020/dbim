package discovery

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Discovery struct {
	conn *redis.Client
}

func NewDiscovery(redisAddr string) *Discovery {
	fmt.Println("Redis addr", redisAddr)
	opts, err := redis.ParseURL(redisAddr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(opts)
	rc := redis.NewClient(opts)
	if err := rc.Ping().Err(); err != nil {
		fmt.Println(err)
		return nil
	}
	d := &Discovery{conn: rc}
	return d
}

func (d *Discovery) GetCometAddr() []string {
	num := d.conn.LLen("cometlist").Val()

	var res []string
	val, _ := d.conn.LRange("cometlist", 0, num).Result()
	for _, v := range val {
		res = append(res, v)
	}
	return res
}

func (d *Discovery) RegComet(addr string) {
	pipe := d.conn.TxPipeline()
	pipe.RPush("cometlist", addr)
	_, err := pipe.Exec()
	if err != nil {
		fmt.Println("Register comet failed")
	}
}

func (d *Discovery) DelComet(addr string) {
	pipe := d.conn.TxPipeline()
	pipe.LRem("cometlist", 1, addr)
	_, err := pipe.Exec()
	if err != nil {
		fmt.Println("Del comet failed", err)
	}
}
