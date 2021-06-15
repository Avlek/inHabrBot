package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type DB struct {
	conn       *redis.Client
	expiration time.Duration
}

func NewRedisConnect(config RedisConfig) *DB {
	return &DB{
		conn: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		}),
		expiration: time.Duration(config.Expiration) * time.Hour,
	}
}

func (db *DB) Get(ctx context.Context, key string) (interface{}, error) {
	cmd := db.conn.Get(ctx, key)
	if cmd.Err() == redis.Nil {
		return nil, nil
	}
	return cmd.Val(), cmd.Err()
}

func (db *DB) Set(ctx context.Context, key string, val interface{}) error {
	cmd := db.conn.Set(ctx, key, val, db.expiration)
	return cmd.Err()
}
