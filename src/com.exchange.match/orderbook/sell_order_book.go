package match

import (
	. "exchange-match/src/com.exchange.match/domain"
)

// 成交前提是 买方价 >= 卖方价
// 出高价买入的订单  优先成交, 因此按价格 倒序 (大顶堆)
// 出低价卖的订单 优先成交  因此按价格 升序    (小顶堆)

// 卖方订单薄
// 按价格从小到大排序
var sellOrderBook = make([]*OrderInfo, 0)

// 订单薄大小
var sellOrderCount = 0

// 订单列表大小
func SellOrderBookSize() int {
	return sellOrderCount
}

// 获取订单列表
func SellOrderBook() []*OrderInfo {
	return sellOrderBook
}

// 添加订单
func PushSellOrder(order *OrderInfo) {
	sellOrderBook = append(sellOrderBook, order)
	sellOrderCount++
	sellOrderFixUp()
	return
}

// 获取堆顶元素 不弹出
func PeekSellOrder() *OrderInfo {
	if sellOrderCount <= 0 {
		return nil
	}
	return sellOrderBook[0]
}

func PollSellOrder() *OrderInfo {
	if sellOrderCount <= 0 {
		return nil
	}
	root := sellOrderBook[0]
	sellOrderCount--
	sellOrderBook[0] = sellOrderBook[sellOrderCount]
	sellOrderFixDown()
	return root
}

// 插入元素: 拆入元素后放在数组最后一个位置上
// 将最后一个元素上浮到它应该在的位置
func sellOrderFixUp() {
	targetNum := sellOrderCount - 1
	for {
		if targetNum <= 0 {
			break
		}
		parentIndex := (targetNum - 1) / 2

		if sellOrderBook[targetNum].Price > sellOrderBook[parentIndex].Price {
			break
		}
		sellOrderBook[targetNum], sellOrderBook[parentIndex] = sellOrderBook[parentIndex], sellOrderBook[targetNum]
		targetNum = parentIndex
	}
}

// 取出元素 取出元素后 将最后一个元素放在根元素上
// 元素下沉 将元素移动到它合适的地方
func sellOrderFixDown() {
	rootIndex := 0
	length := sellOrderCount
	for {
		leftChildIndex := 2*rootIndex + 1
		rightChildIndex := 2*rootIndex + 2

		// 孩子节点中较小的节点
		childIndex := leftChildIndex
		// 防止越界
		if leftChildIndex >= length {
			// 至少存在左孩子
			break
		}
		if rightChildIndex < length {
			// 存在 右孩子
			if sellOrderBook[rightChildIndex].Price < sellOrderBook[leftChildIndex].Price {
				childIndex = rightChildIndex
			}
		}
		if sellOrderBook[rootIndex].Price <= sellOrderBook[childIndex].Price {
			break
		}
		sellOrderBook[rootIndex], sellOrderBook[childIndex] = sellOrderBook[childIndex], sellOrderBook[rootIndex]
		rootIndex = childIndex
	}
}
