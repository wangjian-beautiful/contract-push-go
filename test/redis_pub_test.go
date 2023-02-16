package test

import (
	"fmt"
	"gcex-contract-go/redis"
	"testing"
)

// 向Redis发布10000条数据
func TestRedisPub(t *testing.T) {
	// 初始化连接redis
	err := redis.InitRedis("127.0.0.1:7001", "123456789", 0)
	if err != nil {
		fmt.Printf("InitRedis error: %s\n", err)
		return
	}
	topic := "market_e_compusdt_ticker"
	msg := `{
    "event":"sub",
    "params":{
        "channel":"market_$symbol_depth_step0", // $symbol E.g. 币币:btcusdt 合约:e_btcusdt
        "cb_id":"1" // 业务id 非必填
		}
	}`
	for i := 0; i < 10000; i++ {
		redis.Publish(topic, msg)
	}

}
