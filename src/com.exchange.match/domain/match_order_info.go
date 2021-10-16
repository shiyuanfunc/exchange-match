package domain

import "exchange-match/src/com.exchange.match/util"

type MatchOrderInfo struct {
	MatchId     int64   `json:"matchId"`
	BuyOrderId  int64   `json:"buyOrderId"`
	BuyUserId   string  `json:"buyUserId"`
	SellOrderId int64   `json:"sellOrderId"`
	SellUserId  string  `json:"sellUserId"`
	DealAmount  float64 `json:"dealAmount"`
	DealPrice   float64 `json:"dealPrice"`
}

func BuildMatchOrder(sellOrder OrderInfo, buyOrder OrderInfo, dealAmount float64, dealPrice float64) MatchOrderInfo {
	return MatchOrderInfo{
		MatchId:     util.NextId(),
		BuyOrderId:  buyOrder.OrderId,
		BuyUserId:   buyOrder.UserId,
		SellOrderId: sellOrder.OrderId,
		SellUserId:  sellOrder.UserId,
		DealAmount:  dealAmount,
		DealPrice:   dealPrice,
	}
}
