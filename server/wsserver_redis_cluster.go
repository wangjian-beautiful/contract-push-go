package server

import (
	"context"
	"encoding/json"
	"fmt"
	"gcex-contract-go/config"
	"gcex-contract-go/consts"
	"gcex-contract-go/mongo"
	"gcex-contract-go/redis"
	"gcex-contract-go/websocket"
	websocket2 "github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type EventRepResult struct {
	EventRep string      `json:"event_rep"`
	Channel  string      `json:"channel"`
	CbId     string      `json:"cb_id"`
	Ts       int64       `json:"ts"`
	Data     interface{} `json:"data"`
	Status   string      `json:"status"`
	Code     string
}

type ReviewData struct {
	Amount any `json:"amount"`
	Close  any `json:"close"`
	High   any `json:"high"`
	Low    any `json:"low"`
	Open   any `json:"open"`
	Vol    any `json:"vol"`
}

var ctx = context.Background()

type WsServer struct {
	Ws             *websocket.Server
	authentication bool
	SubPub         *SubPub
}

func NewNominalServer() *WsServer {
	return &WsServer{
		Ws: websocket.New(websocket.Config{
			ReadBufferSize:   1024,
			WriteBufferSize:  10240,
			BinaryMessages:   true,
			EvtMessagePrefix: []byte("pc:"),
			GzipContent:      true,
			Authentication:   false,
		}),
		authentication: false,
	}
}

func NewAuthServer() *WsServer {
	return &WsServer{
		Ws: websocket.New(websocket.Config{
			ReadBufferSize:   1024,
			WriteBufferSize:  10240,
			BinaryMessages:   false,
			EvtMessagePrefix: []byte("bjs:"),
			CheckOrigin: func(r *http.Request) bool {
				token := r.URL.Query().Get("token")
				if len(token) == 0 {
					return false
				}
				return true
			},
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				w.WriteHeader(status)
				w.Write([]byte(fmt.Sprintf("authentication fail %v", reason)))
			},
			Authentication: true,
		}),
		authentication: true,
	}
}

func (w *WsServer) Start() (err error) {
	w.Ws.OnConnection(w.handleConnection)
	w.SubPub = NewSubPub(w.Ws)
	return nil
}

func (w *WsServer) handleConnection(c websocket.Connection) {
	log.Println("client connected,id=", c.ID())
	c.Write(1, []byte("welcome client"))
	c.OnDisconnect(func() {
		log.Printf("Total conn %d", c.Server().GetTotalConnections())
		log.Println("client Disconnect,id=", c.ID())
	})
	// 接收到发布消息事件
	c.On(consts.EventPub, func(subMsg *consts.SubMsg) {
		// c .Context（）是gin的http上下文。
		log.Printf("%s pub publish: %v\n", c.Context().ClientIP(), subMsg)
		if data, e := json.Marshal(subMsg.Params); e == nil {
			//发布消息到Redis
			redis.Cluster.Publish(redis.Cluster.Context(), subMsg.Params["channel"], string(data))
			//redis.Publish(subMsg.Params["channel"], string(data))
		}
	})

	// 接收到订阅的事件
	c.On(consts.EventSub, func(subMsg *consts.SubMsg) {
		// 将消息打印到控制台，
		log.Printf("%s sub msg: %v\n", c.Context().ClientIP(), subMsg)
		channel := subMsg.Params["channel"]
		if consts.IsAuthChannel(channel) && !w.authentication {
			c.Write(websocket2.TextMessage, []byte("未认证接口， 不允许订阅账号信息，请走/auth"))
			return
		}
		w.SubPub.Subscribe(channel, c.ID(), c.Uid())

	})
	// 接收到取消订阅的事件
	c.On(consts.EventUnsub, func(subMsg *consts.SubMsg) {
		log.Printf("%s unsub  msg: %v\n", c.Context().ClientIP(), subMsg)
		w.SubPub.UnSubscribe(subMsg.Params["channel"], c.ID(), c.Uid())
	})

	c.On(consts.EventReq, func(subMsg *consts.SubMsg) {
		result := HandleReqChannel(subMsg.Params)
		marshal, err := json.Marshal(result)
		if err != nil {
			log.Printf("EventReq json 转换出错 %v \t%v", subMsg, err)
		}
		c.WriteDefault(marshal)
	})
}

func HandleReqChannel(params map[string]string) any {
	channel := params["channel"]
	if consts.IsKlineChannel(channel) {
		endIdx, _ := strconv.ParseInt(params["endIdx"], 10, 64)
		pageSize, _ := strconv.ParseInt(params["pageSize"], 10, 64)
		if pageSize == 0 {
			pageSize = 100
		}
		channel := params["channel"]
		var repResult EventRepResult
		kLines := mongo.GetKline(channel, endIdx, pageSize)
		if kLines == nil {
			kLines = []mongo.Kline{}
		}

		repResult.Data = kLines
		repResult.Channel = channel
		repResult.EventRep = "rep"
		repResult.Ts = time.Now().UnixMilli()
		return repResult
	}
	if consts.IsTradeTicker(channel) {
		channel := params["channel"]
		var repResult EventRepResult
		repResult.Data = []interface{}{mongo.GetTrade(channel, 10)}
		repResult.Channel = channel
		repResult.EventRep = "rep"
		repResult.Ts = time.Now().UnixMilli()
		return repResult
	}
	if consts.ReqChannelReview == channel {
		var repResult EventRepResult
		allTicker := consts.GetAllTickerChannel(config.Config.Symbols...)
		var data = make(map[string]ReviewData)
		for _, channel := range allTicker {
			val, _ := redis.Cluster.HGetAll(ctx, consts.RedisTickerKey+channel).Result()
			log.Printf("ReqChannelReview kay:%s,value:%s", consts.RedisTickerKey+channel, val)
			split := strings.Split(channel, "_")
			if len(split) > 3 {
				if val == nil {
					val = make(map[string]string, 0)
				}
				key := split[1] + "_" + split[2]
				var RData = ReviewData{
					Open:   redisResultStrArrayToString(val["open"]),
					Close:  redisResultStrArrayToString(val["close"]),
					High:   redisResultStrArrayToString(val["high"]),
					Low:    redisResultStrArrayToString(val["low"]),
					Amount: redisResultStrArrayToString(val["amount"]),
					Vol:    redisResultStrArrayToString(val["vol"]),
				}
				data[key] = RData
			}
		}
		repResult.Data = data
		repResult.Channel = consts.ReqChannelReview
		repResult.Ts = time.Now().UnixMilli()
		repResult.Status = "ok"
		repResult.EventRep = consts.EventReq
		return repResult
	}
	if consts.IsTicker(channel) {
		var repResult EventRepResult
		val, err := redis.Cluster.Get(ctx, consts.RedisTickerKey+channel).Result()
		if err != nil {

		}
		repResult.Data = val
		repResult.Channel = channel
		repResult.Ts = time.Now().UnixMilli()
		repResult.Status = "ok"
		repResult.EventRep = consts.EventReq
		return repResult
	}
	if strings.HasPrefix(channel, consts.FundingRateChannelPrefix) {
		symbol := strings.TrimPrefix(channel, consts.FundingRateChannelPrefix)
		if len(symbol) > 0 {
			var repResult EventRepResult
			symbol = strings.ToUpper(symbol)
			fundingRateMap := make(map[string]string)
			fundingRateString, _ := redis.Cluster.HGet(ctx, consts.FundingRateKey, symbol).Result()
			latestPriceString, _ := redis.Cluster.HGet(ctx, consts.LatestPriceKey, symbol).Result()
			err := json.Unmarshal([]byte(fundingRateString), &fundingRateMap)
			log.Print("资金费率redis 数据", fundingRateString, latestPriceString)
			if err != nil {
				log.Printf("%s value:%s 转换出错 %v", consts.FundingRateKey+symbol, fundingRateString, err)
			}
			data := struct {
				CurrentFundRate string `json:"currentFundRate"`
				IndexPrice      string `json:"indexPrice"`
				TagPrice        string `json:"tagPrice"`
				NextFundRate    string `json:"nextFundRate"`
			}{IndexPrice: latestPriceString, TagPrice: latestPriceString}
			data.NextFundRate = fundingRateMap["nextFundRate"]
			data.CurrentFundRate = fundingRateMap["currentFundRate"]
			repResult.Data = data
			repResult.Channel = channel
			repResult.Ts = time.Now().UnixMilli()
			repResult.Status = "ok"
			repResult.EventRep = consts.EventReq
			repResult.Code = "0"
			return repResult
		}
	}
	return nil
}

func redisResultStrArrayToString(data string) any {
	var arr []any
	if err := json.Unmarshal([]byte(data), &arr); err != nil {
		log.Printf("redisResultStrArrayToString:%s\t%v", data, err)
	}
	return arr[1]
}
