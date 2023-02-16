package redis

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io"
	"log"
	"strings"
	"time"
)

var (
	// RedisClient RD redis全局client
	RedisClient *redis.Pool
	maxretry    = 3 //redis 发布错误重试次数
)

// InitRedis 初始设置
func InitRedis(host string, auth string, db int) error {
	// 连接Redis
	RedisClient = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   4000,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host, redis.DialPassword(auth), redis.DialDatabase(db))
			if nil != err {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	rd := RedisClient.Get()
	defer rd.Close()

	c, err := redis.Dial("tcp", host, redis.DialPassword(auth), redis.DialDatabase(db))
	defer c.Close()
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return err
	}
	log.Println("Connect to redis ok")
	return nil

}

// Redo 在pool加入TestOnBorrow方法来去除扫描坏连接
func Redo(command string, opt ...interface{}) (interface{}, error) {
	if RedisClient == nil {
		return "", errors.New("error,redis client is null")
	}
	rd := RedisClient.Get()
	defer rd.Close()

	var conn redis.Conn
	var err error

	var needNewConn bool

	resp, err := rd.Do(command, opt...)
	needNewConn = IsConnError(err)
	if needNewConn == false {
		return resp, err
	} else {
		conn, err = RedisClient.Dial()
	}

	for index := 0; index < maxretry; index++ {
		if conn == nil && index+1 > maxretry {
			return resp, err
		}
		if conn == nil {
			conn, err = RedisClient.Dial()
		}
		if err != nil {
			continue
		}

		resp, err := conn.Do(command, opt...)
		needNewConn = IsConnError(err)
		if needNewConn == false {
			return resp, err
		} else {
			conn, err = RedisClient.Dial()
		}
	}

	conn.Close()
	return "", errors.New("redis error")
}

func IsConnError(err error) bool {
	var needNewConn bool

	if err == nil {
		return false
	}

	if err == io.EOF {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "use of closed network connection") {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "connect: connection refused") {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "connection closed") {
		needNewConn = true
	}
	return needNewConn
}
