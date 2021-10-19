package domain

import (
	"encoding/json"
	"log"
	"strconv"
)

type TradeContractOrder struct {
	ContractId     int64   `json:"contractId"`     //合约Id
	OrderDirection int     `json:"orderDirection"` // 下单方向
	OrderId        int64   `json:"orderId"`        // 订单ID
	OrderType      int8    `json:"orderType"`      // 下单方向 1空; 2多
	PositionMode   int8    `json:"positionMode"`   // 仓位模式 1空 2多
	Price          float64 `json:"price"`          // 委托价格
	SubjectId      int     `json:"subjectId"`      // 合约Id
	SubjectName    string  `json:"subjectName"`    // 合约名字
	TotalAmount    float64 `json:"totalAmount"`    // 委托数量
	TradersId      int64   `json:"tradersId"`      //交易员Id
	UserId         int64   `json:"userId"`         // 下单用户id
}

// 解析订单对象
func ParseTradeContractOrderOrder(bytes []byte) *TradeContractOrder {
	var tradeContractOrder TradeContractOrder
	err := json.Unmarshal(bytes, &tradeContractOrder)
	if err != nil {
		log.Println("下单消息解析异常, params:", string(bytes))
		return nil
	}
	return &tradeContractOrder
}

func ConvertTradeOrder2OrderBookInfo(bytes []byte) *OrderInfo {
	var tradeContractOrder TradeContractOrder
	err := json.Unmarshal(bytes, &tradeContractOrder)
	if err != nil {
		log.Println("下单消息解析异常, params:", string(bytes))
		return nil
	}
	orderInfo := OrderInfo{
		Price:          tradeContractOrder.Price,
		Amount:         tradeContractOrder.TotalAmount,
		UserId:         strconv.FormatInt(tradeContractOrder.UserId, 16),
		OrderDirection: tradeContractOrder.OrderDirection,
		OrderId:        tradeContractOrder.OrderId,
	}
	return &orderInfo
}
