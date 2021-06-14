package impl

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type DB struct {
	conn *redis.Client
}

func NewRedisConnect(opts *redis.Options) *DB {
	return &DB{
		conn: redis.NewClient(opts),
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
	cmd := db.conn.Set(ctx, key, val, 24*time.Hour)
	return cmd.Err()
}
