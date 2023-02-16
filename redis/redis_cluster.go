package redis

import (
	"context"
	"fmt"
	"gcex-contract-go/config"
	"github.com/go-redis/redis/v8"
	"time"
)

var Cluster *redis.ClusterClient

var ctx = context.Background()

func init() {
	Cluster = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        config.Config.Redis.Adders,
		DialTimeout:  3000 * time.Microsecond,
		ReadTimeout:  3000 * time.Microsecond,
		WriteTimeout: 3000 * time.Microsecond,
		Password:     config.Config.Redis.Auth,
	})
	// 发送一个ping命令,测试是否通
	s := Cluster.Do(ctx, "ping").String()
	fmt.Println(s)
}
