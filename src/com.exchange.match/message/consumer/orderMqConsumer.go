package consumer

import (
	"context"
	"exchange-match/src/com.exchange.match/domain"
	"exchange-match/src/com.exchange.match/service"
	"exchange-match/src/com.exchange.match/util"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"log"
	"os"
)

var defaultClient rocketmq.PushConsumer

func init() {
	// RocketMQ 日志
	rlog.SetLogLevel("warn")
	start()
}

func start() {

	client, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(util.GetConsumerGroup()),
		consumer.WithNameServer(util.GetRocketNameSrv()),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeMessageBatchMaxSize(1),
		consumer.WithInstance(util.GetConsumerGroup()),
	)
	if err != nil {
		fmt.Println("consumer init error >>> ", err)
		return
	}

	// 合约下单
	initErr := client.Subscribe(util.TradeContractOrder, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			msg := msgs[i]
			orderInfo := domain.ConvertTradeOrder2OrderBookInfo(msg.Body)
			if orderInfo == nil {
				log.Println("消息解析对象异常 params: ", string(msg.Body), err)
				return consumer.ConsumeSuccess, nil
			}
			// handler orderInfo
			service.ReceivedOrderInfo(*orderInfo)
		}
		return consumer.ConsumeSuccess, nil
	})
	if initErr != nil {
		fmt.Println("consumer sub error >>> ", initErr)
		os.Exit(-1)
	}

	// 合约单取消
	_ = client.Subscribe(util.TradeContractOrderCancelReq, consumer.MessageSelector{}, handlerCancelOrder)

	err = client.Start()
	if err != nil {
		fmt.Println("consumer start error >>> ", err)
		os.Exit(-1)
	}
	defaultClient = client
	//go shutdown()
	//<-util.GetChan()
	////客户端自带的优雅退出
	//client.Shutdown()
}

// 优雅退出
func ShutDownConsumer() {
	if defaultClient != nil {
		_ = defaultClient.Shutdown()
	}
}
