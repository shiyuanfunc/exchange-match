package shutdown

import (
	"exchange-match/src/com.exchange.match/message/consumer"
	"exchange-match/src/com.exchange.match/message/producer"
	"log"
	"time"
)

const (
	YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"
)

// 全局优雅退出
func ShutDown() {
	log.Println("ShutDown ALL >>>> ", time.Now().Format(YYYY_MM_DD_HH_MM_SS))
	consumer.ShutDownConsumer()
	producer.ShutDownProducer()
}
