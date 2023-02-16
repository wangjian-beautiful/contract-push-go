package server

import (
	"gcex-contract-go/config"
	"gcex-contract-go/consts"
	"gcex-contract-go/redis"
	"gcex-contract-go/websocket"
	"log"
	"sync"
)

type SubPub struct {
	Ws    *websocket.Server
	cbMap sync.Map
}

func initRedisSubChannel() (channels []string) {
	allKlineChannel := consts.GetAllKlineChannel(config.Config.Symbols...)
	allTicker := consts.GetAllTickerChannel(config.Config.Symbols...)
	return append(allKlineChannel, allTicker...)
}

//var ctx = context.Background()

func NewSubPub(ws *websocket.Server) *SubPub {
	subPub := &SubPub{
		Ws:    ws,
		cbMap: sync.Map{},
	}
	channels := initRedisSubChannel()
	pubSub := redis.Cluster.Subscribe(ctx, channels...)
	//defer func(pubsub *redis.PubSub) {
	//	err := pubsub.Close()
	//	if err != nil {
	//
	//	}
	//}(pubsub)
	//for {
	//	msg, err := pubSub.ReceiveMessage(ctx)
	//	if err != nil {
	//		log.Print(msg)
	//	}
	//
	//	fmt.Println(msg.Channel, msg.Payload)
	//}
	go func() {
		for {
			msg, err := pubSub.ReceiveMessage(ctx)
			if err != nil {
				log.Print(msg)
			}

			if msg != nil {
				go subPub.OnReceiveMessage(msg.Channel, msg.Payload)
			}
		}
	}()
	return subPub
}

func (c *SubPub) OnReceiveMessage(topic string, content string) {
	// redis发布认证推送暂时没有， 后续有用上可以放开
	//if consts.EventChannelAccount == topic {
	//	var pushMsf = consts.RedisPushMsg{}
	//	if err := json.Unmarshal([]byte(content), &pushMsf); err != nil {
	//		return
	//	}
	//	uid := pushMsf.Uid
	//	c.cbMap.Range(func(key, value any) bool {
	//		if value == uid {
	//			conn := c.Ws.GetConnection(key.(string))
	//			conn.WriteDefault([]byte(content))
	//		}
	//		return false
	//	})
	//	return
	//}
	conns := c.Ws.GetConnectionsByRoom(topic)
	if conns == nil || len(conns) == 0 {
		log.Printf("未有用户订阅该频道：%s", topic)
		return

	}
	for _, c := range conns {
		c.WriteDefault([]byte(content))
	}
}
func (c *SubPub) Subscribe(channel string, clientId string, uid string) {
	//if consts.EventChannelAccount == channel {
	//	if len(uid) == 0 {
	//		log.Println("订阅账号信息 uid不能为空")
	//		c.Ws.GetConnection(clientId).Write(websocket2.TextMessage, []byte("订阅账号信息 uid不能为空"))
	//		return
	//	}
	//	c.cbMap.Store(clientId, uid)
	//	return
	//}
	c.Ws.Join(channel, clientId)
}

func (c *SubPub) UnSubscribe(channel string, clientId string, uid string) {
	//if consts.EventChannelAccount == channel {
	//	c.cbMap.Delete(clientId)
	//	return
	//}
	c.Ws.Leave(channel, clientId)
}
