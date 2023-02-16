package consts

import "fmt"

const (
	EventSub   = "sub"
	EventReq   = "req"
	EventUnsub = "unsub"
	EventPub   = "pub"
)

const (
	EventChannelAccount = "sub_pub_account"
	//MqTopicMatchTrade   = "match_trade_go_test"
	MqTopicMatchTrade = "match_trade"
	RedisTickerKey    = "Ticker_key_"
	//ReqChannelReview  review channel
	ReqChannelReview = "review"
	PingMessageText  = "PING"
	PongMessageText  = "PONG"
)

const (
	//FundingRateChannelPrefix 资金费率跟价格通道前缀
	FundingRateChannelPrefix = "funding_rate_"
	//FundingRateKey 资金费率
	FundingRateKey = "funding_rate_key"
	//LatestPriceKey 最新价格
	LatestPriceKey = "latest_price_key"
)

const RedisTokenPrefix = "user_"

const (
	// PositionSetPrefix redis用户持仓前缀 position_set_prefix:uid
	PositionSetPrefix  = "position_set_prefix:"
	PositionDataPrefix = "position_data_prefix:"
)

const TimeFormatPatter = "2006-01-02 15:04:05"

const (
	SideBuy  = "BUY"
	SideSell = "SELL"
)

// SubMsg 订阅的消息格式定义
type SubMsg struct {
	ID     string            `json:"id"`    //请求ID
	Event  string            `json:"event"` //订阅时固定为sub,取消订阅时固定为unsub
	Params map[string]string `json:"params"`
}

type RedisPushMsg struct {
	Uid     string            `json:"uid"`
	Source  string            `json:"Source"`
	Channel string            `json:"channel"`
	Ts      int64             `json:"ts"`
	Payload map[string]string `json:"payload"`
}

// PushMsg 平台推送消息定义
type PushMsg struct {
	Channel string                 `json:"channel"`
	Ts      int64                  `json:"ts"`
	tick    map[string]interface{} `json:"tick"`
}

// PushDepthMsg 推送深度
type PushDepthMsg struct {
	Channel string                 `json:"channel"`
	Ts      int64                  `json:"ts"`
	tick    map[string]interface{} `json:"tick"`
}

type LPFloat struct {
	Value  float64
	Digits int32
}

func (l LPFloat) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.*f", l.Digits, l.Value)
	return []byte(s), nil
}

type User struct {
	AuthLevel                 int    `json:"authLevel"`
	AuthStatus                int    `json:"authStatus"`
	AuthType                  int    `json:"authType"`
	CapitalPword              string `json:"capitalPword"`
	CertificateTime           int64  `json:"certificateTime"`
	CountryCode               string `json:"countryCode"`
	Ctime                     int64  `json:"ctime"`
	Email                     string `json:"email"`
	ExcStatus                 int    `json:"excStatus"`
	FamilyName                string `json:"familyName"`
	GoogleAuthenticatorKey    string `json:"googleAuthenticatorKey"`
	GoogleAuthenticatorStatus int    `json:"googleAuthenticatorStatus"`
	Id                        int    `json:"id"`
	LastLoginTime             int64  `json:"lastLoginTime"`
	LoginPword                string `json:"loginPword"`
	LoginStatus               int    `json:"loginStatus"`
	MobileAuthenticatorStatus int    `json:"mobileAuthenticatorStatus"`
	MobileNumber              string `json:"mobileNumber"`
	Mtime                     int64  `json:"mtime"`
	Name                      string `json:"name"`
	Nickname                  string `json:"nickname"`
	RealName                  string `json:"realName"`
	RealnameTime              int64  `json:"realnameTime"`
	ShowMobileNumber          string `json:"showMobileNumber"`
	UserType                  int    `json:"userType"`
	WithdrawStatus            int    `json:"withdrawStatus"`
}
