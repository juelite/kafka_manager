package common

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

type RedisClient struct {
}

// 获取redis实例
func (this *RedisClient) GetInstance() redis.Conn {
	var client redis.Conn
	var host, pass string
	host = beego.AppConfig.String("REDIS_HOST")
	pass = beego.AppConfig.String("REDIS_PASS")
	client, err := redis.Dial("tcp", host)
	if err != nil {
	}
	_, err = client.Do("AUTH", pass)
	if err != nil {
	}
	return client
}

// redis分布式锁 - 加锁
func (this *RedisClient) Lock(lockKey string, uUid string, ttl uint64) (bool, error) {
	client := this.GetInstance()
	defer client.Close()

	resp, err := redis.Int(client.Do("SET", lockKey, uUid, "NX", "EX", ttl))
	if err != nil {
		return false, err
	}
	if resp == 0 {
		return false, nil
	}
	return true, nil
}

// redis分布式锁 - 释放锁
func (this *RedisClient) UnLock(lockKey string, uUid string) (bool, error) {
	luaScript := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	rScript := redis.NewScript(1, luaScript)
	client := this.GetInstance()
	defer client.Close()

	resp, err := redis.Int(rScript.Do(client, lockKey, uUid))
	if err != nil {
		return false, err
	}
	if resp == 0 {
		return false, nil
	}
	return true, nil
}
