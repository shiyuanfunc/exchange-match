package domain

import (
	"encoding/json"
	"log"
)

type ContractOrderCancelReq struct {
	CancelType string  `json:"cancelType"`
	IdList     []int64 `json:"idList"`
	UserId     int64   `json:"userId"`
}

type ContractOrderCancelResult struct {
	OrderId      int64 `json:"orderId"`
	UserId       int64 `json:"userId"`
	CancelResult bool  `json:"cancelResult"`
}

func ParseContractOrderReq(bytes []byte) *ContractOrderCancelReq {
	var contractOrderReq ContractOrderCancelReq
	err := json.Unmarshal(bytes, &contractOrderReq)
	if err != nil {
		log.Println("解析取消订单消息异常, params:", string(bytes))
		return nil
	}
	return &contractOrderReq
}

// 构建合约单取消结果
func BuildContractOrderCancelResult(orderId int64, userId int64, cancelResult bool) *ContractOrderCancelResult {
	return &ContractOrderCancelResult{
		OrderId:      orderId,
		UserId:       userId,
		CancelResult: cancelResult,
	}
}
