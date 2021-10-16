package main

import (
	"exchange-match/src/com.exchange.match/util"
	"fmt"
)

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println(util.NextId())
	}
}
