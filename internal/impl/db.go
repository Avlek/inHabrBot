package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type DB struct {
	conn       *redis.Client
	channelID  int64
	adminID    int64
	expiration time.Duration
}

func NewRedisConnect(config *Config) *DB {
	return &DB{
		conn: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		}),
		channelID:  config.Telegram.ChannelID,
		adminID:    config.Telegram.AdminID,
		expiration: time.Duration(config.Redis.Expiration) * time.Hour,
	}
}

func (db *DB) Get(ctx context.Context, key string) (interface{}, error) {
	fullKey := fmt.Sprintf("post%d_%s", db.channelID, key)
	cmd := db.conn.Get(ctx, fullKey)
	if cmd.Err() == redis.Nil {
		return nil, nil
	}
	return cmd.Val(), cmd.Err()
}

func (db *DB) Set(ctx context.Context, key string, val interface{}) error {
	fullKey := fmt.Sprintf("post%d_%s", db.channelID, key)
	cmd := db.conn.Set(ctx, fullKey, val, db.expiration)
	return cmd.Err()
}
