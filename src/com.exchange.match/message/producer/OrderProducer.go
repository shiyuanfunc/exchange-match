package producer

import (
	"context"
	"exchange-match/src/com.exchange.match/util"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	MqRetryTimes = 2
)

var defaultProducer rocketmq.Producer

// 获取生产者
func init() {
	log.SetPrefix("[message producer] ")
	mqProducer, err := rocketmq.NewProducer(
		producer.WithNameServer(util.GetRocketNameSrv()),
		producer.WithRetry(MqRetryTimes),
		producer.WithGroupName(util.GetConfig().ProducerGroup),
		producer.WithQueueSelector(producer.NewRandomQueueSelector()),
		producer.WithVIPChannel(false),
	)
	if err != nil {
		log.Println("init producer error", err.Error())
		os.Exit(-1)
	}
	err = mqProducer.Start()
	if err != nil {
		log.Println("start producer error", err.Error())
		os.Exit(-1)
	}
	log.Println("producer started success >>>> ")
	defaultProducer = mqProducer
	return
}

func newProducer() rocketmq.Producer {
	return defaultProducer
}

func SendMessage(obj interface{}) bool {
	message := primitive.NewMessage(util.TEST_TOPIC, []byte(util.ToJsonString(obj)))
	return sendMessage(message)
}

// 发送消息
func sendMessage(message *primitive.Message) bool {
	message.WithTag("*")
	sync, err := defaultProducer.SendSync(context.Background(), message)
	if err != nil {
		log.Println("Send message error", err)
		return false
	}
	log.Println(util.ToJsonString(sync))
	return true
}

func main() {
	pro := newProducer()
	count := 1
	for {
		str := "this is a msg data " + strconv.Itoa(count)
		message := primitive.NewMessage("item_test", []byte(str))
		sendResult, err := pro.SendSync(context.Background(), message)
		if err != nil {
			fmt.Println("send message error ", err.Error())
			return
		}
		fmt.Println("send message success", sendResult)
		time.Sleep(1 * 1e9)
	}
}

// 优雅退出
func ShutDownProducer() {
	if defaultProducer != nil {
		_ = defaultProducer.Shutdown()
	}
}
