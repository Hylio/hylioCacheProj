package main

import (
	"fmt"
	"time"
)

var set = make(map[int]bool)
var total int

func printOnce(num int) {
	if _, exist := set[num]; !exist {
		total++
	}
	set[num] = true
}

func main() {
	for i := 0; i < 100; i++ {
		go printOnce(100)
	}
	time.Sleep(time.Second)
	fmt.Print(total)
}
