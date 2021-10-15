package service

import (
	"encoding/json"
	. "exchange-match/src/com.exchange.match/domain"
	"exchange-match/src/com.exchange.match/message/producer"
	match "exchange-match/src/com.exchange.match/orderbook"
	"log"
)

func init() {
	log.SetPrefix("[MatchService] ")
}

var totalCount = 0

// 处理订单
func ReceivedOrderInfo(orderInfo OrderInfo) {
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
				return
			}
			if buyOrder.Amount < orderInfo.Amount {
				go sendMatchInfo(*buyOrder, orderInfo, buyOrder.Amount, buyOrder.Price, "处理卖单")
				orderInfo.Amount -= buyOrder.Amount
				match.PollBuyOrder()
			}
		}
	}
	// 2、写入卖方委托单列表
	if orderInfo.Amount > 0 {
		match.PushSellOrder(&orderInfo)
	}
}

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
			// 满足匹配条件
			go logMatchOrder(sellOrder.Price, buyOrder.Price)

			if buyOrder.Amount > sellOrder.Amount {
				// 买方量 > 该笔卖单量 买方未完全成交 该笔卖单完全成交
				go sendMatchInfo(buyOrder, *sellOrder, sellOrder.Amount, sellOrder.Price, "处理买单")
				buyOrder.Amount -= sellOrder.Amount
				match.PollSellOrder()
				continue
			}
			if buyOrder.Amount == sellOrder.Amount {
				// 买单完全成交 该笔卖单完全成交 弹出
				// todo 异步更新 价格对应的委托量
				// 发送成交消息
				go sendMatchInfo(buyOrder, *sellOrder, sellOrder.Amount, sellOrder.Price, "处理买单")
				match.PollSellOrder()
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
