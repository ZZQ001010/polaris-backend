package cache

import (
	"reflect"
	"sync"

	"github.com/gomodule/redigo/redis"
)

var host = "localhost:6379"

var mu sync.Mutex
var single RedisProxy

type RedisProxy struct {
	conn redis.Conn
}

func (rp *RedisProxy) ZAdd(key string, score float64, value string) {
	conn := rp.Connect()
	conn.Do("zadd", key, score, value)
}

func (rp *RedisProxy) SetEx(key string, value string, ex int64) {
	conn := rp.Connect()
	conn.Do("setex", key, ex, value)
}

func (rp *RedisProxy) Get(key string) string {
	conn := rp.Connect()
	rs, _ := conn.Do("get", key)
	if rs == nil {
		return ""
	}
	return string(rs.([]byte))
}

func (rp *RedisProxy) Del(key string) int64 {
	conn := rp.Connect()
	rs, _ := conn.Do("Del", key)
	if rs == nil {
		return 0
	}
	return rs.(int64)
}

func (rp *RedisProxy) Incrby(key string, v int64) int64 {
	conn := rp.Connect()
	rs, _ := conn.Do("INCRBY", key, v)
	return rs.(int64)
}

func (rp *RedisProxy) Exist(key string) bool {
	conn := rp.Connect()
	rs, _ := conn.Do("EXISTS", key)
	if rs == nil {
		return false
	}
	return rs.(int64) == 1
}

func (rp *RedisProxy) Expire(key string, expire int64) bool {
	conn := rp.Connect()
	rs, _ := conn.Do("EXPIRE", key, expire)
	return rs.(int64) == 1
}

func (rp *RedisProxy) Connect() redis.Conn {
	if rp.conn == nil || rp.conn.Err() != nil {
		mu.Lock()
		defer mu.Unlock()
		if rp.conn == nil || rp.conn.Err() != nil {
			conn, _ := GetRedisConn()
			rp.conn = conn
		}
	}
	return rp.conn
}

func (rp *RedisProxy) Close() {
	if rp.conn != nil && rp.conn.Err() == nil {
		rp.conn.Close()
	}
}

func (rp RedisProxy) IsEmpty() bool {
	return reflect.DeepEqual(rp, RedisProxy{})
}

//GetProxy get redis oper proxy
func GetProxy() *RedisProxy {
	return &RedisProxy{}
}
