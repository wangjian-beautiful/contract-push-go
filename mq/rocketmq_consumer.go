package mq

import (
	"fmt"
	"gcex-contract-go/config"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"os"
)

var C rocketmq.PushConsumer

func init() {
	r := config.Config.Rocketmq
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(r.Producer.Group),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{r.NameServer})),
		//consumer.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
	)

	if err != nil {
		fmt.Println("init consumer error: " + err.Error())
		os.Exit(0)
	}
	C = c
}
