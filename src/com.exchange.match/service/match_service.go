package service

import (
	"encoding/json"
	. "exchange-match/src/com.exchange.match/domain"
	"exchange-match/src/com.exchange.match/message/producer"
	match "exchange-match/src/com.exchange.match/orderbook"
	"fmt"
	"log"
	"sync"
)

func init() {
	log.SetPrefix("[MatchService] ")
}

var totalCount = 0

// 订单map <orderId, orderId>
// orderId存在 代表该orderId已经被撮合或者被取消
var orderMap sync.Map

// 处理订单
func ReceivedOrderInfo(orderInfo OrderInfo) {
	if checkOrderHasHandler(orderInfo.OrderId) {
		log.Println(fmt.Sprintf("订单orderId:%d已被处理", orderInfo.OrderId))
		return
	}
	if orderInfo.OrderDirection == 1 {
		// 卖单
		handlerBidOrder(orderInfo)
	} else {
		// 买单
		handlerBuyOrder(orderInfo)
	}
}

// 处理卖单
func handlerBidOrder(orderInfo OrderInfo) {
	// 1、获取买单元素
	if match.BuyOrderBookSize() > 0 {
		for {
			buyOrder := match.PeekBuyOrder()
			if buyOrder == nil {
				break
			}
			if buyOrder.Price < orderInfo.Price {
				// 买价 < 卖价
				break
			}

			if checkOrderHasHandler(buyOrder.OrderId) {
				// 该笔买单已经被处理过(被取消)
				match.PollBuyOrder()
				clearOrderHandlerMark(buyOrder.OrderId)
				continue
			}
			// 满足匹配条件
			go logMatchOrder(orderInfo.Price, buyOrder.Price)
			// 买价 >= 卖价
			if buyOrder.Amount > orderInfo.Amount {
				// 买方未完全成交 卖单完全成交
				go sendMatchInfo(*buyOrder, orderInfo, orderInfo.Amount, orderInfo.Price, "处理卖单")
				buyOrder.Amount -= orderInfo.Amount
				return
			}
			if buyOrder.Amount == orderInfo.Amount {
				go sendMatchInfo(*buyOrder, orderInfo, orderInfo.Amount, orderInfo.Price, "处理卖单")
				match.PollBuyOrder()
				// 买单已被处理完全
				// 清除标记
				clearOrderHandlerMark(buyOrder.OrderId)
				return
			}
			if buyOrder.Amount < orderInfo.Amount {
				go sendMatchInfo(*buyOrder, orderInfo, buyOrder.Amount, buyOrder.Price, "处理卖单")
				orderInfo.Amount -= buyOrder.Amount
				match.PollBuyOrder()
				// 买单已被处理完全
				// 清除标记
				clearOrderHandlerMark(buyOrder.OrderId)
			}
		}
	}
	// 2、写入卖方委托单列表
	if orderInfo.Amount > 0 {
		// 代表该笔还有待处理部分
		clearOrderHandlerMark(orderInfo.OrderId)
		match.PushSellOrder(&orderInfo)
	}
}

// 处理买单
func handlerBuyOrder(buyOrder OrderInfo) {

	// 1、获取卖单委托列表
	// 2、数量撮合
	// 3、写入买方委托列表
	if match.SellOrderBookSize() > 0 {
		for {
			sellOrder := match.PeekSellOrder()
			if sellOrder == nil {
				break
			}
			if sellOrder.Price > buyOrder.Price {
				// 卖方价格 > 买方价
				break
			}

			if checkOrderHasHandler(sellOrder.OrderId) {
				// 该卖单已被标记处理 (被取消)
				// 弹出该元素
				match.PollSellOrder()
				clearOrderHandlerMark(sellOrder.OrderId)
				continue
			}

			// 满足匹配条件
			go logMatchOrder(sellOrder.Price, buyOrder.Price)

			if buyOrder.Amount > sellOrder.Amount {
				// 买方量 > 该笔卖单量 买方未完全成交 该笔卖单完全成交
				go sendMatchInfo(buyOrder, *sellOrder, sellOrder.Amount, sellOrder.Price, "处理买单")
				buyOrder.Amount -= sellOrder.Amount
				match.PollSellOrder()
				clearOrderHandlerMark(sellOrder.OrderId)
				continue
			}
			if buyOrder.Amount == sellOrder.Amount {
				// 买单完全成交 该笔卖单完全成交 弹出
				// todo 异步更新 价格对应的委托量
				// 发送成交消息
				go sendMatchInfo(buyOrder, *sellOrder, sellOrder.Amount, sellOrder.Price, "处理买单")
				match.PollSellOrder()
				clearOrderHandlerMark(sellOrder.OrderId)
				return
			}

			if buyOrder.Amount < sellOrder.Amount {
				go sendMatchInfo(buyOrder, *sellOrder, buyOrder.Amount, buyOrder.Price, "处理买单")
				// 买方完全成交  卖方未完全成交
				sellOrder.Amount -= buyOrder.Amount
				return
			}
		}
	}

	if buyOrder.Amount > 0 {
		match.PushBuyOrder(&buyOrder)
	}
}

func sendMatchInfo(buyOrder OrderInfo, sellOrder OrderInfo, amount float64, price float64, direction string) {
	//log.Println("[Matched]"+direction, " BuyOrderInfo: ", ToJsonStr(buyOrder), " SellOrderInfo:", ToJsonStr(sellOrder), " dealAmount:", amount, " dealPrice:", price)
	matchOrder := BuildMatchOrder(sellOrder, buyOrder, amount, price)
	producer.SendMessage(matchOrder)
}

// log 条件达成日志
func logMatchOrder(sellPrice float64, buyPrice float64) {
	//log.Println(fmt.Sprintf("[logMatchOrder] sellPrice:%f buyPrice:%f", sellPrice, buyPrice))
	totalCount++
}

func ToJsonStr(obj interface{}) string {
	objByte, _ := json.Marshal(obj)
	return string(objByte)
}

func GetTotalMatched() int {
	return totalCount
}

// 取消委托单
// false 代表取消失败; true代表取消成功
func CancelOrder(orderId int64) bool {
	_, loaded := orderMap.LoadOrStore(orderId, orderId)
	if loaded {
		// 该orderId在map中存在,已经被取消或者被撮合
		return false
	}
	// map中不存在 放入成功
	return true
}

// 校验订单是否被处理过 false 未被处理； true 已被处理
func checkOrderHasHandler(orderId int64) bool {
	_, loaded := orderMap.LoadOrStore(orderId, orderId)
	return loaded
}

func clearOrderHandlerMark(orderId int64) {
	orderMap.Delete(orderId)
}
