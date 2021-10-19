package main

import (
	"exchange-match/src/com.exchange.match/domain"
	_ "exchange-match/src/com.exchange.match/message/consumer"
	_ "exchange-match/src/com.exchange.match/message/producer"
	"exchange-match/src/com.exchange.match/service"
	"exchange-match/src/com.exchange.match/shutdown"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	// 启动统计
	go stat()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)
	// 阻塞 主线程
	<-signals
	// 收到阻塞信号 全局退出
	shutdown.ShutDown()
}

func TestMatchOrder() {
	index := 0
	go stat()
	for {
		orderInfo := domain.NewOrderInfoBatch(1000.0+float64(index), float64(rand.Uint32()%50), "sellUser01"+strconv.Itoa(index), rand.Intn(2), 20211801+int64(index))
		service.ReceivedOrderInfo(orderInfo)
		index++
		//fmt.Println(service.ToJsonStr(orderInfo))
		if index > 100 {
			break
		}
	}
	time.Sleep(8 * time.Second)
}

// 测试
func testMatch() {
	for i := 0; i < 10; i++ {
		amount := 0
		for {
			amount = rand.Intn(30)
			if amount > 0 {
				break
			}
		}
		orderInfo := domain.NewOrderInfoBatch(1000.0+float64(i), float64(amount), "sellUser01"+strconv.Itoa(i), 1, 20211801+int64(i))
		service.ReceivedOrderInfo(orderInfo)
	}

	for i := 0; i < 10; i++ {
		amount := 0
		for {
			amount = rand.Intn(30)
			if amount > 0 {
				break
			}
		}
		orderInfo := domain.NewOrderInfoBatch(300.0+float64(i), float64(amount), "buyUserId02"+strconv.Itoa(i), 2, 20221801+int64(i))
		service.ReceivedOrderInfo(orderInfo)
	}

	amount := 0
	for {
		amount = rand.Intn(30)
		if amount > 0 {
			break
		}
	}
	dealBuyOrderInfo := domain.NewOrderInfoBatch(1500, float64(amount), "buyUserIdadmin", 2, 20231801)
	service.ReceivedOrderInfo(dealBuyOrderInfo)

	for {
		amount = rand.Intn(20)
		if amount > 0 {
			break
		}
	}
	dealSellOrderInfo := domain.NewOrderInfoBatch(100, float64(amount), "sellUserIdadmin", 1, 20241801)
	service.ReceivedOrderInfo(dealSellOrderInfo)
}

// 统计
func stat() {
	ticker := time.NewTicker(1 * time.Second)
	lastCount := 0
	for range ticker.C {
		nowCount := service.GetTotalMatched()
		log.Println("[Stat] >>>>>> ", time.Now().Format("2006-01-02 15:04:05"), " diff ", nowCount-lastCount)
		lastCount = nowCount
	}
	ticker.Stop()

}

func initLogger() log.Logger {
	file := "./" + "stat" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[Stat]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return log.Logger{}
}
