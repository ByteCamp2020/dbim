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
		pipe.RPush(key, 0)
		// set ttl
		pipe.Expire(key, per)
		_, _ = pipe.Exec()
	} else {
		l.rc.RPushX(key, 0)
	}

	return true
}

func (l *Limiter) UserLimit(userID string) bool {
	getRes := l.rc.Get(userID)
	pipe:= l.rc.TxPipeline()
	pipe.Set(userID, 0, l.Dur)
	_, _ = pipe.Exec()
	if getRes.Err() != redis.Nil {
		return false
	}
	return true
}
