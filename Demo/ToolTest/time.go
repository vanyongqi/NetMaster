package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	for {
		cur := time.Now()
		elapsed := cur.Sub(start)
		if elapsed.Seconds() > 5 {
			fmt.Println("Time is up!")
			break
		}
	}
}
