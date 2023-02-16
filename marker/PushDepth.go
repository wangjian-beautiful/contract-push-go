package marker

import (
	"context"
	"encoding/json"
	"gcex-contract-go/consts"
	"gcex-contract-go/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

type LPFloat = consts.LPFloat

type PushDepthMsg struct {
	Channel string `json:"channel"`
	Ts      int64  `json:"ts"`
	Tick    struct {
		Asks [][]LPFloat `json:"asks"`
		Buys [][]LPFloat `json:"buys"`
	} `json:"tick"`
}

type TradePlate struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

func GetDepthMap() map[string][]int32 {
	depthMap := make(map[string][]int32)
	for _, ws := range wss {
		rooms := ws.Ws.GetRoomNames()
		for _, room := range rooms {
			if consts.IsDepthChannel(room) {
				symbol := strings.Split(room, "_")[2]
				depth, _ := strconv.Atoi(room[len(room)-1:])
				if _, ok := depthMap[symbol]; ok {
					depthMap[symbol] = append(depthMap[symbol], int32(depth))
				} else {
					depthMap[symbol] = []int32{int32(depth)}
				}
			}
		}
	}
	return depthMap
}

var ctx = context.Background()

func StarterDepth() {
	handlerDepthSub := func() {
		depthMap := GetDepthMap()
		for symbol, depths := range depthMap {
			for _, depth := range depths {
				buyTradRedisKey := consts.GetTradePlateRedisKey(symbol, "BUY")
				sellTradRedisKey := consts.GetTradePlateRedisKey(symbol, "SELL")
				buyTradeData, _ := redis.Cluster.Get(ctx, buyTradRedisKey).Bytes()
				sellTradData, _ := redis.Cluster.Get(ctx, sellTradRedisKey).Bytes()
				log.Printf("深度redis数据 %s:%s\t%s:%s", buyTradRedisKey, string(buyTradeData), sellTradRedisKey, string(sellTradData))
				buys := GroupTrade(buyTradeData, depth)
				asks := GroupTrade(sellTradData, depth)
				channel := consts.GetDepthChannel(symbol, depth)
				pushDepthMsg := &PushDepthMsg{}
				pushDepthMsg.Channel = channel
				pushDepthMsg.Ts = time.Now().UnixMilli()
				if buys == nil {
					buys = [][]LPFloat{}
				}
				if asks == nil {
					asks = [][]LPFloat{}
				}
				pushDepthMsg.Tick.Buys = buys
				pushDepthMsg.Tick.Asks = asks
				msg, err := json.Marshal(pushDepthMsg)
				if err != nil {
					log.Printf("深度结果转换出错%v\t%v", pushDepthMsg, err)
					return
				}
				conns := GetAllConnectionsByRoom(channel)
				for _, conn := range conns {
					conn.WriteDefault(msg)
				}
			}
		}
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				go handlerDepthSub()
			}
		}
	}()
}

func GroupTrade(content []byte, depth int32) (result [][]LPFloat) {
	items := make([]TradePlate, 0)
	content = consts.JsonStringConvert(content)
	if json.Unmarshal(content, &items) != nil {
		log.Printf("GroupTrade json格式错误%V")
		return
	}
	for _, item := range items {
		price := item.Price
		amount := item.Amount
		isGroup := false
		for _, d := range result {
			if d[0].Value == (price) {
				d[1].Value = d[1].Value + amount
				isGroup = true
				break
			}
		}
		if isGroup {
			continue
		}
		result = append(result, []LPFloat{{price, 2}, {amount, 2}})
	}
	return

}
