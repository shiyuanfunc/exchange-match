package util

var TradeData string = "contract_trade_data"

func GetOrderTopic() string {
	orderTopic := "Order_topic"
	return orderTopic
}

var quit = make(chan string)

func GetChan() chan string {
	return quit
}
