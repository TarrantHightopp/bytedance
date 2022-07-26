package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisOpts struct {
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                   // redis的哪个数据库
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`             // 服务器地址:端口
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
	Context  context.Context
}

type Redis struct {
	Client  *redis.Client
	Context context.Context
}

func NewRedis(opt *RedisOpts) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("redis connect ping failed, err: ", err)
		panic(err)
	} else {
		fmt.Println("redis connect ping response: ", pong)
		return &Redis{Client: client, Context: opt.Context}
	}
}

func (r *Redis) Set(key string, val interface{}, timeout time.Duration) error {
	return r.Client.Set(r.Context, key, val, timeout).Err()
}

func (r *Redis) Get(key string) (string, error) {
	return r.Client.Get(r.Context, key).Result()
}

func (r *Redis) IsExist(key string) bool {
	return false
}
func (r *Redis) Delete(key string) error {
	return nil
}
