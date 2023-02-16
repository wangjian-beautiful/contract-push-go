package test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"gcex-contract-go/config"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/configor"
	"io"
	"log"
	"os"
	"os/signal"
	"testing"
	"time"
)

// 测试每秒qps为100并订阅全部频道 并10000台客户端同时在线
func TestWsServerConn(t *testing.T) {
	//install signal
	//err := redis.InitRedis("127.0.0.1:7001", "123456789", 0)
	//if err != nil {
	//	fmt.Printf("InitRedis error: %s\n", err)
	//	return
	//}

	//configor.Load(&config.Config, "../"+config.GetProfilesConf())
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for i := 0; i < 10; i++ {
		go ConnWs(interrupt)
		if i%100 == 0 {
			time.Sleep(time.Second * 1)
		}
	}
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
		}
	}

}

func TestGzip(t *testing.T) {
	configor.Load(&config.Config, "../"+config.GetProfilesConf())
	data := []byte("122288882asdfkjadskfkjasdjfjaksd")
	dozip := doGzip(data)
	data = UnGZipBytes(dozip)
	dataStr := string(UnGZipBytes(dozip))
	log.Printf("123%s", dataStr)
	fmt.Printf("data:%s", string(UnGZipBytes(dozip)))

}

func UnGZipBytes(data []byte) []byte {
	var out bytes.Buffer
	var in bytes.Buffer
	in.Write(data)
	r, _ := gzip.NewReader(&in)
	r.Close() //这句放在后面也没有问题，不写也没有任何报错
	//机翻注释：关闭关闭读者。它不会关闭底层的io.Reader。为了验证GZIP校验和，读取器必须完全使用，直到io.EOF。

	io.Copy(&out, r) //这里我看了下源码不是太明白，
	//我个人想法是这样的，Reader本身就是go中表示一个压缩文件的形式，r转化为[]byte就是一个符合压缩文件协议的压缩文件
	return out.Bytes()
}

/*
*
压缩bytes内容
1.根据指定目录创建文件
2.根据文件资源对象生成gzip Writer对象
3.往gzip Writer对象写入内容
*/
func doGzip(data []byte) []byte {
	var in bytes.Buffer
	gzipWriter := gzip.NewWriter(&in)
	defer gzipWriter.Close()
	_, err := gzipWriter.Write(data)
	if err != nil {
		log.Print(err)
	}
	return in.Bytes()
}

func ConnWs(interrupt chan os.Signal) {
	//ws := `ws://localhost:9090/ws?`
	ws := `wss://futuresws.gcex41.com/kline-api/ws`
	c, _, err := websocket.DefaultDialer.Dial(ws, nil)
	if err != nil {
		log.Println("dial:", err)
	}
	defer c.Close()
	fmt.Println("connect...")

	acceptMsg := func() {
		for {
			_, data, _ := c.ReadMessage()
			//c.SetCompressionLevel()
			//c.ReadMessage()
			//data := UGZipBytes(message)
			//msg := make(map[string]interface{})
			//if err := json.Unmarshal(data, &msg); err != nil {
			//	return
			//}
			//if err != nil {
			//	log.Println("read:", err)
			//	return
			//}
			//redis.Publish(msg["channel"].(string), string(data))
			log.Printf("recv : \n %v \n\n", string(data))
		}
	}

	heartbeat := func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		keepLive := `{"ping"}`
		for {
			select {
			case <-ticker.C:
				if _err := c.WriteMessage(websocket.PingMessage, []byte(fmt.Sprintf(keepLive))); _err != nil {
					log.Println("ws心跳失败")
					return
				}
			}
		}
	}

	subscribe := func(symbol string) {
		subscribeMember := `{"event":"sub","params":{"channel":"market_e_%s_ticker","cb_id":"e_%s"}}`
		subscribeMember = fmt.Sprintf(subscribeMember, symbol, symbol)
		if err_ := c.WriteMessage(websocket.TextMessage, []byte(subscribeMember)); err != nil {
			log.Fatalln("subscribe failed :", err_)
		}
	}
	go acceptMsg()
	go heartbeat()
	for _, symbol := range config.Config.Symbols {
		subscribe(symbol)
	}
	//// 退出信号
	done := make(chan struct{})
	defer close(done)
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			if err_ := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err_ != nil {
				log.Println("write close:", err)
				return
			}
		}
	}

}

func TestTicker(t *testing.T) {
	fmt.Printf("123ddd")
	ticker := time.NewTicker(500 * time.Millisecond)
	handlerDepthSub := func() {
		log.Println("ddddd")
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				go handlerDepthSub()
			}
		}
	}()
	done := make(chan struct{})
	for {
		select {
		case <-done:
			return
		}
	}

}
