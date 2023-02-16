// // 发布 订阅
// // 默认使用redis，如需要更换 修改Subscriber.client
package redis

//
//import (
//	"gcex-contract-go/consts"
//	"github.com/gomodule/redigo/redis"
//	websocket2 "github.com/gorilla/websocket"
//
//	"gcex-contract-go/websocket"
//	"log"
//	"sync"
//)
//
//func Publish(topic string, msg string) (interface{}, error) {
//	resp, err := Redo("Publish", topic, msg)
//	if err != nil {
//		log.Println(err)
//	}
//	return resp, err
//}
//
//type SubscribeCallback func(topicMap sync.Map, topic, msg string)
//
//type Subscriber struct {
//	client   redis.PubSubConn
//	Ws       *websocket.Server //websocket
//	cbMap    sync.Map
//	CallBack interface {
//		OnReceive(SubscribeCallback)
//	}
//}
//
//var fnSubReceived SubscribeCallback
//
//func (c *Subscriber) OnReceive(cb SubscribeCallback) {
//	fnSubReceived = cb
//}
//
//func (c *Subscriber) Init(ws *websocket.Server) {
//
//	conn := RedisClient.Get()
//
//	c.client = redis.PubSubConn{conn}
//	c.Ws = ws
//	go func() {
//		for {
//			log.Println("redis wait...")
//			switch res := c.client.Receive().(type) {
//			case redis.Message:
//				topic := res.Channel
//				message := string(res.Data)
//				fnSubReceived(c.cbMap, topic, message)
//			case redis.Subscription:
//				log.Printf("%s: %s %d\n", res.Channel, res.Kind, res.Count)
//			case error:
//				log.Println("error handle", res)
//				if IsConnError(res) {
//					conn, err := RedisClient.Dial()
//					if err != nil {
//						log.Printf("err=%s\n", err)
//					}
//					c.client = redis.PubSubConn{conn}
//				}
//				continue
//			}
//		}
//	}()
//
//}
//
//func (c *Subscriber) IniCluster(ws *websocket.Server) {
//
//	conn := RedisClient.Get()
//
//	c.client = redis.PubSubConn{conn}
//	c.Ws = ws
//	go func() {
//		for {
//			log.Println("redis wait...")
//			switch res := c.client.Receive().(type) {
//			case redis.Message:
//				//fmt.Printf("receive:%#v\n", res)
//				topic := res.Channel
//				message := string(res.Data)
//				fnSubReceived(c.cbMap, topic, message)
//			case redis.Subscription:
//				log.Printf("%s: %s %d\n", res.Channel, res.Kind, res.Count)
//			case error:
//				log.Println("error handle", res)
//				if IsConnError(res) {
//					conn, err := RedisClient.Dial()
//					if err != nil {
//						log.Printf("err=%s\n", err)
//					}
//					c.client = redis.PubSubConn{conn}
//				}
//				continue
//			}
//		}
//	}()
//
//}
//func (c *Subscriber) Close() {
//	err := c.client.Close()
//	if err != nil {
//		log.Println("redis close error.")
//	}
//}
//
//func (c *Subscriber) Subscribe(channel string, clientid string, uid string) {
//	err := c.client.Subscribe(channel)
//	if err != nil {
//		log.Println("redis Subscribe error.", err)
//	}
//	if consts.EventChannelAccount == channel {
//		if len(uid) == 0 {
//			log.Println("订阅账号信息 uid不能为空")
//			c.Ws.GetConnection(clientid).Write(websocket2.TextMessage, []byte("订阅账号信息 uid不能为空"))
//			return
//		}
//		c.cbMap.Store(clientid, uid)
//		return
//	}
//	c.Ws.Join(channel, clientid)
//}
//
//func (c *Subscriber) PSubscribe(channel interface{}, clientid string) {
//	err := c.client.PSubscribe(channel)
//	if err != nil {
//		log.Println("redis PSubscribe error.", err)
//	}
//	//c.cbMap.Store(clientid, channel.(string))
//}
//
//func (c *Subscriber) UnSubscribe(channel string, clientid string, uid string) {
//	if consts.EventChannelAccount == channel {
//		c.cbMap.Delete(clientid)
//		return
//	}
//	c.Ws.Leave(channel, clientid)
//}
