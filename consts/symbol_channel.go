package consts

import (
	"fmt"
	"strings"
)

const CoinTradePlate = "CONTRACT_TRADE_PLATE:"
const marketESuffix = "market_e"

const PositionChannel = marketESuffix + "_position"

var allKlinePeriod = []string{"1min", "5min", "15min", "30min", "60min", "4hours", "1day", "1week", "1month"}

// GetTradePlateRedisKey redis TradePlate key
func GetTradePlateRedisKey(symbol string, side string) string {
	return fmt.Sprintf("%s%s_%s", CoinTradePlate, strings.ToUpper(symbol), side)
}

func IsDepthChannel(channel string) bool {
	return strings.Contains(channel, "depth") && strings.HasSuffix(channel[:len(channel)-1], "step")
}

func IsFundingRateChannel(channel string) bool {
	return strings.HasPrefix(channel, FundingRateChannelPrefix)
}

func IsPositionChannel(channel string) bool {
	return strings.EqualFold(PositionChannel, channel)
}

func GetDepthChannel(symbol string, depth int32) string {
	return fmt.Sprintf("market_e_%s_depth_step%d", symbol, depth)
}

func GetTickerChannel(symbol string) string {
	return fmt.Sprintf("market_e_%s_ticker", strings.ToLower(symbol))
}

func GetTradeTickerChannel(symbol string) string {
	return fmt.Sprintf("market_e_%s_trade_ticker", strings.ToLower(symbol))
}

func GetKlineChannel(symbol string, period string) string {
	return fmt.Sprintf("market_e_%s_kline_%s", strings.ToLower(symbol), period)
}

func GetAllKlineChannel(symbols ...string) (channels []string) {
	for _, symbol := range symbols {
		for i := range allKlinePeriod {
			join := strings.Join([]string{marketESuffix, symbol, "kline", allKlinePeriod[i]}, "_")
			channels = append(channels, join)
		}
	}
	return
}
func GetAllTickerChannel(symbols ...string) (channels []string) {
	for _, symbol := range symbols {
		//添加24小时ticker通道
		tickerChannel := strings.Join([]string{marketESuffix, symbol, "ticker"}, "_")
		channels = append(channels, tickerChannel)
	}
	return
}

func IsKlineChannel(channel string) bool {
	split := strings.Split(channel, "_")
	return len(split) > 4 && strings.Compare(split[3], "kline") == 0
}

func IsTradeTicker(channel string) bool {
	return len(channel) > 0 && strings.HasSuffix(channel, "_trade_ticker")
}

func IsTicker(channel string) bool {
	return len(channel) > 0 && strings.HasSuffix(channel, "_ticker")
}
