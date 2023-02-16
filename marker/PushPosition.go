package marker

import (
	"encoding/json"
	"gcex-contract-go/consts"
	"gcex-contract-go/redis"
	"log"
	"strings"
	"time"
)

func getPositionCacheDataKey(symbol string, uid string, side string) string {
	return strings.Join([]string{consts.PositionDataPrefix, uid, symbol, side}, "")

}

func StarterPositionPush() {
	pushPositionSub := func() {
		conns := GetAllConnectionsByRoom(consts.PositionChannel)
		for _, conn := range conns {
			var repResult EventSubResult
			positionSetVal, _ := redis.Cluster.SMembers(ctx, consts.PositionSetPrefix+conn.Uid()).Result()
			var resultData = make([]any, 0)
			if len(positionSetVal) > 0 {
				for _, pkv := range positionSetVal {
					pkv = consts.DeleteDoubleQuotationMark(pkv)
					if strings.HasSuffix(pkv, consts.SideBuy) {
						symbol := strings.TrimSuffix(pkv, consts.SideBuy)
						positionCacheDataKey := getPositionCacheDataKey(symbol, conn.Uid(), consts.SideBuy)
						redisPosiData, err := redis.Cluster.HGetAll(redis.Cluster.Context(), positionCacheDataKey).Result()
						if err == nil {
							resultData = append(resultData, redisPosiData)
						}
					}
					if strings.HasSuffix(pkv, consts.SideSell) {
						symbol := strings.TrimSuffix(pkv, consts.SideSell)
						positionCacheDataKey := getPositionCacheDataKey(symbol, conn.Uid(), consts.SideSell)
						redisPosiData, err := redis.Cluster.HGetAll(redis.Cluster.Context(), positionCacheDataKey).Result()
						if err == nil {
							resultData = append(resultData, redisPosiData)
						}
					}
				}

			}
			log.Printf("用户持仓数据 %s\t%v", conn.Uid(), resultData)

			repResult.Data = resultData
			repResult.Channel = consts.PositionChannel
			repResult.Ts = time.Now().UnixMilli()
			repResult.Status = "ok"
			repResult.Code = "0"
			msg, _ := json.Marshal(repResult)
			conn.WriteDefault(msg)
		}

	}
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				go pushPositionSub()
			}
		}
	}()
}
