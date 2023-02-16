package mongo

import (
	"context"
	"fmt"
	"gcex-contract-go/config"
	"gcex-contract-go/consts"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
	"time"
)

type LPFloat = consts.LPFloat
type Decimal128 = primitive.Decimal128
type Decimal = decimal.Decimal

var client mongo.Client
var database *mongo.Database

type Kline struct {
	Id     int64   `json:"id" bson:"time"`
	Amount Decimal `json:"amount"`
	Vol    Decimal `json:"vol"`
	Open   Decimal `json:"open"`
	Close  Decimal `json:"close"`
	High   Decimal `json:"high"`
	Low    Decimal `json:"low"`
}

type MongoKline struct {
	Id     int64      `json:"id" bson:"time"`
	Amount Decimal128 `json:"amount"`
	Vol    Decimal128 `json:"vol"`
	Open   Decimal128 `json:"open"`
	Close  Decimal128 `json:"close"`
	High   Decimal128 `json:"high"`
	Low    Decimal128 `json:"low"`
}

type Trade struct {
	Side   string     `json:"side"`
	Price  Decimal128 `json:"price"`
	Vol    Decimal128 `json:"vol"`
	Amount Decimal128 `json:"amount"`
}

type OriginTrade struct {
	Amount    Decimal128 `json:"amount"`
	Price     Decimal128 `json:"price"`
	TrendSide string     `json:"trendSide"`
	Turnover  Decimal128 `json:"turnover"`
}
type EventRepTradeResult struct {
	EventRep string  `json:"event_rep"`
	Channel  string  `json:"channel"`
	CbId     string  `json:"cb_id"`
	Ts       int64   `json:"ts"`
	Status   string  `json:"status"`
	Data     []Trade `json:"data"`
}

func init() {
	mongodbConf := config.Config.Mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//auth := options.Credential{
	//	AuthSource: mongodbConf.AuthenticationDatabase,
	//	Username:   mongodbConf.Username,
	//	Password:   mongodbConf.Password,
	//}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+mongodbConf.Address))
	if err != nil {
		log.Printf("init mongo error %v", err)
	}
	database = client.Database(mongodbConf.Database)
}

func GetKline(channel string, endIdx int64, pageSize int64) (result []Kline) {
	collection := database.Collection(channel)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{"_id", -1}})
	if pageSize > 0 {
		opts = opts.SetLimit(pageSize)
	}
	d := bson.D{}
	if endIdx > 0 {
		d = bson.D{{"time", bson.D{{"$lt", endIdx}}}}
	}
	cur, err := collection.Find(ctx, d, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var mongoK MongoKline
		s := cur.Current.String()
		fmt.Println(s)
		err := cur.Decode(&mongoK)
		if err != nil {
			log.Println(err)
		}
		if mongoK.Id > 100000000000 {
			mongoK.Id = mongoK.Id / 1000
		}
		k := Kline{
			Id:     mongoK.Id,
			Amount: Decimal128ToDecimal(mongoK.Amount),
			Vol:    Decimal128ToDecimal(mongoK.Vol),
			Open:   Decimal128ToDecimal(mongoK.Open),
			Close:  Decimal128ToDecimal(mongoK.Close),
			High:   Decimal128ToDecimal(mongoK.High),
			Low:    Decimal128ToDecimal(mongoK.Low),
		}
		result = append(result, k)
	}
	// 升序排序(稳定排序)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Id < result[j].Id {
			return true
		}
		return false
	})

	if err := cur.Err(); err != nil {
		log.Println(err)
	}
	return
}

func GetTrade(channel string, pageSize int64) (result []Trade) {
	collection := database.Collection(channel)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{{"_id", -1}})
	if pageSize > 0 {
		opts = opts.SetLimit(pageSize)
	}
	cur, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Printf("mongo select err%v", err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var trade OriginTrade
		err := cur.Decode(&trade)
		if err != nil {
			log.Println(err)
		}
		result = append(result, Trade{
			Side:   trade.TrendSide,
			Price:  trade.Price,
			Vol:    trade.Amount,
			Amount: trade.Turnover,
		})
	}

	if err := cur.Err(); err != nil {
		log.Printf("mongo select err%v", err)
	}
	return
}

func Decimal128ToDecimal(d128 Decimal128) Decimal {
	bigInt, i, _ := d128.BigInt()
	return decimal.NewFromBigInt(bigInt, int32(i))
}
