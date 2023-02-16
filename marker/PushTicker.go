package marker

import (
	"context"
	"encoding/json"
	"fmt"
	"gcex-contract-go/consts"
	"gcex-contract-go/mq"
	"gcex-contract-go/server"
	. "gcex-contract-go/websocket"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

type Decimal = decimal.Decimal
type MatchTradeDetailsDTO struct {
	//价格
	Price float64
	//数量
	Amount float64
	//成交额
	Turnover float64
	//主动单方向
	TrendSide string
	//买订单号
	BuyOrderId int64
	//卖订单号
	SellOrderId int64
	//是否完成
	BuyFinish bool
	Time      int64
}

type MatchTradeDTO struct {
	//撤单id
	OrderId int64
	//币对
	Symbol string
	//消息类型
	//撮合数据
	DetailsDTO MatchTradeDetailsDTO
}

type PushTickerMsg struct {
	Channel string `json:"channel"`
	Ts      int64  `json:"ts"`
	Tick    struct {
		Id   int64 `json:"id"`
		Ts   int64 `json:"ts"`
		Data []struct {
			Side   string  `json:"side"`
			Price  float64 `json:"price"`
			Vol    float64 `json:"vol"`
			Amount float64 `json:"amount"`
			Ds     string  `json:"ds"`
		} `json:"data"`
	} `json:"tick"`
}

var wss []*server.WsServer

func AddWSService(server *server.WsServer) {
	wss = append(wss, server)
}

func GetAllConnectionsByRoom(channel string) []Connection {
	var conns []Connection
	for _, ws := range wss {
		connRoom := ws.Ws.GetConnectionsByRoom(channel)
		conns = append(conns, connRoom...)
	}
	return conns
}
func PushTradeTicker(msg []byte) {
	var trade = MatchTradeDTO{}
	err := json.Unmarshal(msg, &trade)
	if err != nil {
		log.Printf("ticker消息转换错误 msg:%s  error%v", msg, err)
		return
	}
	channel := consts.GetTradeTickerChannel(trade.Symbol)
	conns := GetAllConnectionsByRoom(channel)
	if len(conns) > 0 {
		for _, c := range conns {
			pushTickerMsg := &PushTickerMsg{}
			pushTickerMsg.Channel = channel
			pushTickerMsg.Ts = time.Now().Unix()
			tick := &pushTickerMsg.Tick
			tick.Id = trade.DetailsDTO.BuyOrderId
			tick.Ts = time.Now().UnixMilli()
			data := struct {
				Side   string  `json:"side"`
				Price  float64 `json:"price"`
				Vol    float64 `json:"vol"`
				Amount float64 `json:"amount"`
				Ds     string  `json:"ds"`
			}{trade.DetailsDTO.TrendSide,
				trade.DetailsDTO.Price,
				trade.DetailsDTO.Amount,
				trade.DetailsDTO.Amount,
				""}
			data.Ds = time.UnixMilli(trade.DetailsDTO.Time).Format(consts.TimeFormatPatter)
			pushTickerMsg.Tick.Data = append(pushTickerMsg.Tick.Data, data)
			pushMsg, err := json.Marshal(pushTickerMsg)
			if err != nil {
				return
			}
			c.WriteDefault(pushMsg)
		}
	}
}

//func main() {
//	push := MatchTradeDTO{
//		OrderId: 0,
//		Symbol:  "btcusdt",
//		DetailsDTO: MatchTradeDetailsDTO{
//			Price:     decimal.RequireFromString("123"),
//			Amount:    decimal.RequireFromString("123"),
//			Turnover:  decimal.RequireFromString("123"),
//			TrendSide: "BUY",
//		},
//	}
//	msg, _ := json.Marshal(push)
//	PushTradeTicker(msg)
//}

func StarterMqPush() {
	err := mq.C.Subscribe(consts.MqTopicMatchTrade, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("subscribe callback: %v \n", msgs)
		for _, msg := range msgs {
			PushTradeTicker(msg.Body)
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = mq.C.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
}
