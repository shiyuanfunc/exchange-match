package domain

import (
	"encoding/json"
	"errors"
)

// 订单薄结构
type OrderInfo struct {
	Price          float64 `json:"price"`
	Amount         float64 `json:"amount"`
	UserId         string  `json:"userId"`
	OrderDirection int     `json:"orderDirection"`
	OrderId        int64   `json:"orderId"`
}

// 初始化订单对象
func NewOrderInfo(price float64, amount float64) OrderInfo {
	return OrderInfo{Price: price, Amount: amount}
}

func NewOrderInfoBatch(price float64, amount float64, userId string, orderDirection int, orderId int64) OrderInfo {
	return OrderInfo{
		Price:          price,
		Amount:         amount,
		UserId:         userId,
		OrderDirection: orderDirection,
		OrderId:        orderId,
	}
}

// 消息反解成对象
func ParseObject(bytes []byte) (*OrderInfo, error) {
	var orderInfo OrderInfo
	err := json.Unmarshal(bytes, &orderInfo)
	if err != nil {
		return nil, errors.New("测试")
	}
	return &orderInfo, nil
}
