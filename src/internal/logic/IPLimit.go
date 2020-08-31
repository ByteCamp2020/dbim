package logic

import (
	"bdim/src/internal/logic/conf"
	"time"

	"github.com/go-redis/redis"
)

type Limiter struct {
	// Redis client connection.
	rc *redis.Client
	Count int64
	Dur   time.Duration
}

func NewLimiter(c conf.HTTPServer) (*Limiter, error) {
	opts, err := redis.ParseURL(c.RedisAddr)
	if err != nil {
		return nil, err
	}
	rc := redis.NewClient(opts)
	if err := rc.Ping().Err(); err != nil {
		return nil, err
	}

	return &Limiter{
		rc:    rc,
		Count: c.Count,
		Dur:   c.Dur,
	}, nil
}

func (l *Limiter) Allow(key string, events int64, per time.Duration) bool {
	curr := l.rc.LLen(key).Val()
	if curr >= events {
		return false
	}

	if v := l.rc.Exists(key).Val(); v == 0 {
		pipe := l.rc.TxPipeline()
		pipe.RPush(key, key)
		// set ttl
		pipe.Expire(key, per)
		_, _ = pipe.Exec()
	} else {
		l.rc.RPushX(key, key)
	}

	return true
}
