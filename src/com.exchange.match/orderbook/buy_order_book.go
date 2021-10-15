package match

import (
	. "exchange-match/src/com.exchange.match/domain"
)

// 成交前提是 买方价 >= 卖方价
// 出高价买入的订单  优先成交, 因此按价格 倒序 (大顶堆)
// 出低价卖的订单 优先成交  因此按价格 升序    (小顶堆)

// 买方订单薄
// 按价格从大到小排序
var buyOrderList = make([]*OrderInfo, 0)

// 订单薄大小
var count = 0

// 订单列表大小
func BuyOrderBookSize() int {
	return count
}

// 获取订单列表
func BuyOrderBook() []*OrderInfo {
	return buyOrderList
}

// 添加订单
func PushBuyOrder(order *OrderInfo) {
	buyOrderList = append(buyOrderList, order)
	count++
	fixUp()
	return
}

// 获取堆顶元素 不弹出
func PeekBuyOrder() *OrderInfo {
	if count <= 0 {
		return nil
	}
	return buyOrderList[0]
}

func PollBuyOrder() *OrderInfo {
	if count <= 0 {
		return nil
	}
	root := buyOrderList[0]
	count--
	buyOrderList[0] = buyOrderList[count]
	fixDown()
	return root
}

// 插入元素: 拆入元素后放在数组最后一个位置上
// 将最后一个元素上浮到它应该在的位置
func fixUp() {
	targetNum := count - 1
	for {
		if targetNum <= 0 {
			break
		}
		parentIndex := (targetNum - 1) / 2

		if buyOrderList[targetNum].Price < buyOrderList[parentIndex].Price {
			break
		}
		buyOrderList[targetNum], buyOrderList[parentIndex] = buyOrderList[parentIndex], buyOrderList[targetNum]
		targetNum = parentIndex
	}
}

// 取出元素 取出元素后 将最后一个元素放在根元素上
// 元素下沉 将元素移动到它合适的地方
func fixDown() {
	rootIndex := 0
	length := count
	for {
		leftChildIndex := 2*rootIndex + 1
		rightChildIndex := 2*rootIndex + 2

		// 孩子节点中较大的节点
		maxChildIndex := leftChildIndex
		// 防止越界
		if leftChildIndex >= length {
			// 至少存在左孩子
			break
		}
		if rightChildIndex < length {
			// 存在 右孩子
			if buyOrderList[rightChildIndex].Price > buyOrderList[leftChildIndex].Price {
				maxChildIndex = rightChildIndex
			}
		}
		if buyOrderList[rootIndex].Price >= buyOrderList[maxChildIndex].Price {
			break
		}
		buyOrderList[rootIndex], buyOrderList[maxChildIndex] = buyOrderList[maxChildIndex], buyOrderList[rootIndex]
		rootIndex = maxChildIndex
	}
}
