package marker

import (
	"encoding/json"
	"gcex-contract-go/consts"
	"gcex-contract-go/redis"
	"log"
	"strings"
	"time"
)

type EventSubResult struct {
	Channel string      `json:"channel"`
	CbId    string      `json:"cb_id"`
	Ts      int64       `json:"ts"`
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
	Code    string
}

func GetFundingRateRooms() (channels []string) {
	for _, ws := range wss {
		rooms := ws.Ws.GetRoomNames()
		for _, room := range rooms {
			if consts.IsFundingRateChannel(room) {
				channels = append(channels, room)
			}
		}
	}
	return
}

func StarterFundingRatePush() {
	handFundingRateSub := func() {
		rooms := GetFundingRateRooms()
		for _, channel := range rooms {
			symbol := strings.TrimPrefix(channel, consts.FundingRateChannelPrefix)
			if len(symbol) > 0 {
				var repResult EventSubResult
				symbol = strings.ToUpper(symbol)
				fundingRateMap := make(map[string]string)
				fundingRateString, _ := redis.Cluster.HGet(ctx, consts.FundingRateKey, symbol).Result()
				latestPriceString, _ := redis.Cluster.HGet(ctx, consts.LatestPriceKey, symbol).Result()
				log.Print("资金费率redis 数据", fundingRateString, latestPriceString)
				if len(fundingRateString) > 0 {
					err := json.Unmarshal([]byte(fundingRateString), &fundingRateMap)
					if err != nil {
						log.Printf("%s value:%s 转换出错 %v", consts.FundingRateKey+symbol, fundingRateString, err)
					}
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
				repResult.Code = "0"
				conns := GetAllConnectionsByRoom(channel)
				msg, _ := json.Marshal(repResult)
				for _, conn := range conns {
					conn.WriteDefault(msg)
				}
			}

		}
	}
	ticker := time.NewTicker(3000 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				go handFundingRateSub()
			}
		}
	}()
}
