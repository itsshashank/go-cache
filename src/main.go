package main

import (
	"fmt"
	gc "itsshashank/go-cache"
	"time"
)

func main() {
	my_cache := gc.New()
	my_cache.Set(1, 2, time.Minute*1)
	my_cache.Set(1.2, "a", time.Minute*2)
	fmt.Println("Sleeping for 1 min")
	time.Sleep(time.Minute * 1)
	getvalueof(1, my_cache)
	getvalueof(1.2, my_cache)
	fmt.Println("Sleeping for 1 min")
	time.Sleep(time.Minute * 1)
	getvalueof(1, my_cache)
	getvalueof(1.2, my_cache)
}

func getvalueof(key any, my_c *gc.Cache) {
	value, ok := my_c.Get(key)
	if !ok {
		fmt.Printf("key %v Not found\n", key)
	} else {
		fmt.Printf("Found %v has %v\n", key, value)
	}
}
