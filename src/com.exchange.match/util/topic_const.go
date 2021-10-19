package util

var TradeData = "contract_trade_data"

// 交易下单MQ
var TradeContractOrder = "trade_contract_order"

func GetOrderTopic() string {
	orderTopic := "Order_topic"
	return orderTopic
}

var quit = make(chan string)

func GetChan() chan string {
	return quit
}
