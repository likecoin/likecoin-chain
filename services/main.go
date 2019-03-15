package main

import (
	"math/rand"
	"time"

	"github.com/likecoin/likechain/services/cmd"
)

func main() {
	now := time.Now()
	seed := int64(now.Second())*1000000000 + int64(now.Nanosecond())
	rand.Seed(seed)
	cmd.Execute()
}
