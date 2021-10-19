package util

var TradeData = "contract_trade_data"

// 交易下单MQ
var TradeContractOrder = "trade_contract_order"

// 取消合约单
var TradeContractOrderCancelReq = "trade_contract_cancel_req"

// 合约单取消结果
var TradeContractOrderCanceled = "shiyuan-trade_contract_canceled"

func GetOrderTopic() string {
	orderTopic := "Order_topic"
	return orderTopic
}

var quit = make(chan string)

func GetChan() chan string {
	return quit
}
