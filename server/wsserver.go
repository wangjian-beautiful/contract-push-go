package server

// 第一版 redis 走的单链接，先不用
//
//import (
//	"encoding/json"
//	"fmt"
//	websocket2 "github.com/gorilla/websocket"
//	"log"
//	"net/http"
//	"sync"
//	"v2/consts"
//	"v2/redis"
//	"v2/websocket"
//)
//
//var sub redis.Subscriber
//
//type WsServer struct {
//	Ws             *websocket.Server
//	authentication bool
//	SubPub    *redis.SubPub
//}
//
//func NewNominalServer() *WsServer {
//	return &WsServer{
//		Ws: websocket.New(websocket.Config{
//			ReadBufferSize:   1024,
//			WriteBufferSize:  10240,
//			BinaryMessages:   false,
//			EvtMessagePrefix: []byte("bjs:"),
//		}),
//		authentication: false,
//	}
//}
//
//func NewAuthServer() *WsServer {
//	return &WsServer{
//		Ws: websocket.New(websocket.Config{
//			ReadBufferSize:   1024,
//			WriteBufferSize:  10240,
//			BinaryMessages:   false,
//			EvtMessagePrefix: []byte("bjs:"),
//			CheckOrigin: func(r *http.Request) bool {
//				token := r.URL.Query().Get("token")
//				if len(token) == 0 {
//					return false
//				}
//				return true
//			},
//			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
//				w.WriteHeader(status)
//				w.Write([]byte(fmt.Sprintf("authentication fail %v", reason)))
//			},
//		}),
//		authentication: true,
//	}
//}
//
//func (w *WsServer) Start() (err error) {
//	sub.Init(w.Ws)
//	w.Ws.OnConnection(w.handleConnection)
//	//redis client收到的消息分发到websocket
//	sub.OnReceive(func(topicMap sync.Map, topic, msg string) {
//		log.Printf("sub msg received,topic=%s,msg=%s\n", topic, msg)
//		if consts.EventChannelAccount == topic && w.authentication {
//			var pushMsf = consts.RedisPushMsg{}
//			if err = json.Unmarshal([]byte(msg), &pushMsf); err != nil {
//				return
//			}
//			uid := pushMsf.Uid
//			topicMap.Range(func(key, value any) bool {
//				if value == uid {
//					conn := w.Ws.GetConnection(key.(string))
//					conn.WriteDefault([]byte(msg))
//				}
//				return false
//			})
//			return
//		}
//
//		conns := sub.Ws.GetConnectionsByRoom(topic)
//		if conns == nil || len(conns) == 0 {
//			log.Printf("未有用户订阅该频道：%s", topic)
//		}
//		for _, c := range conns {
//			c.WriteDefault([]byte(msg))
//		}
//	})
//	return nil
//}
//
//func (w *WsServer) handleConnection(c websocket.Connection) {
//	log.Println("client connected,id=", c.ID())
//	c.Write(1, []byte("welcome client"))
//	c.OnDisconnect(func() {
//		log.Printf("Total conn %d", c.Server().GetTotalConnections())
//		log.Println("client Disconnect,id=", c.ID())
//	})
//	// 接收到发布消息事件
//	c.On(consts.EventPub, func(subMsg *consts.SubMsg) {
//		// c .Context（）是gin的http上下文。
//		log.Printf("%s pub publish: %v\n", c.Context().ClientIP(), subMsg)
//		if data, e := json.Marshal(subMsg.Params); e == nil {
//			//发布消息到Redis
//			redis.Publish(subMsg.Params["channel"], string(data))
//		}
//	})
//
//	// 接收到订阅的事件
//	c.On(consts.EventSub, func(subMsg *consts.SubMsg) {
//		// 将消息打印到控制台，
//		log.Printf("%s sub msg: %v\n", c.Context().ClientIP(), subMsg)
//		if consts.EventChannelAccount == subMsg.Params["channel"] && !w.authentication {
//			c.Write(websocket2.TextMessage, []byte("未认证接口， 不允许订阅账号信息，请走/auth"))
//			return
//		}
//		sub.Subscribe(subMsg.Params["channel"], c.ID(), subMsg.Params["uid"])
//
//	})
//	// 接收到取消订阅的事件
//	c.On(consts.EventUnsub, func(subMsg *consts.SubMsg) {
//		log.Printf("%s unsub  msg: %v\n", c.Context().ClientIP(), subMsg)
//		sub.UnSubscribe(subMsg.Params["channel"], c.ID(), subMsg.Params["uid"])
//	})
//
//	c.On(consts.EventReq, func(subMsg *consts.SubMsg) {
//		log.Printf("%s req  msg: %v\n", c.Context().ClientIP(), subMsg)
//		//c.Write()
//		//sub.UnSubscribe(subMsg.Params["channel"], c.ID())
//	})
//}
