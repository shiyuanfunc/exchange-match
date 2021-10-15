package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func OrderId() int64 {
	return time.Millisecond.Milliseconds()
}
