package cache

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var conf = map[string]interface{}{
	"MaxIdle":        5,
	"MaxActive":      20,
	"MaxIdleTimeout": 120,
	"Host":           "localhost:6379",
	"Password":       "",
	"Db":             1,
}

var redisClient *redis.Pool

func init() {
	maxIdle := 10
	if v, ok := conf["MaxIdle"]; ok {
		maxIdle = int(v.(int))
	}
	maxActive := 10
	if v, ok := conf["MaxActive"]; ok {
		maxActive = int(v.(int))
	}
	MaxIdleTimeout := 60
	if v, ok := conf["MaxIdleTimeout"]; ok {
		MaxIdleTimeout = int(v.(int))
	}
	timeout := time.Duration(5)

	// 建立连接池
	redisClient = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(MaxIdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", conf["Host"].(string),
				redis.DialPassword(conf["Password"].(string)),
				redis.DialDatabase(int(conf["Db"].(int))),
				redis.DialConnectTimeout(timeout*time.Second),
				redis.DialReadTimeout(timeout*time.Second),
				redis.DialWriteTimeout(timeout*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
}

func GetRedisConn() (redis.Conn, error) {
	return redisClient.Dial()
}
