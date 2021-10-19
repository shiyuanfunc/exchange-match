package consumer

import (
	"context"
	"exchange-match/src/com.exchange.match/domain"
	"exchange-match/src/com.exchange.match/message/producer"
	"exchange-match/src/com.exchange.match/service"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"log"
)

func handlerCancelOrder(ctx context.Context,
	msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	log.Println("received order cancel >>>>> ", len(msgList))
	for _, msg := range msgList {
		contractOrderReq := domain.ParseContractOrderReq(msg.Body)
		if contractOrderReq == nil {
			//
			log.Println("取消订单参数转换异常 >>>>> params", string(msg.Body))
			return consumer.ConsumeRetryLater, nil
		}
		userId := contractOrderReq.UserId
		for _, orderId := range contractOrderReq.IdList {
			cancelResult := service.CancelOrder(orderId)
			// 发送取消MQ
			go sendCancelMq(orderId, userId, cancelResult)
		}
	}
	return consumer.ConsumeSuccess, nil
}

// 发送取消结果
func sendCancelMq(orderId int64, userId int64, cancelResult bool) {
	orderCancelResult := domain.BuildContractOrderCancelResult(orderId, userId, cancelResult)
	producer.SendContractOrderCancelResult(orderCancelResult)
}
