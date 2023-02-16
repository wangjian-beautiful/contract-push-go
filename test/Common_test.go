package test

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	channel := "market_$symbol_depth_step0"
	channel2 := channel[:len(channel)-1]
	fmt.Printf("%s", channel2)
}

func TestTrade(t *testing.T) {
	//var items []marker.TradePlate
	//for i := 0; i < 10; i++ {
	//	items = append(items, marker.TradePlate{
	//		Price:  decimal.RequireFromString("12.91000000000" + strconv.Itoa(i)),
	//		Amount: decimal.RequireFromString("1"),
	//	})
	//}
	//marshal, _ := json.Marshal(items)
	//result := marker.GroupTrade(marshal, 5)
	//bytes, _ := json.Marshal(result)
	//fmt.Printf(string(bytes))

}
