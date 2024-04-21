package chapter7

import (
	"fmt"
	"time"
)

func TickerSlowReceiver() {
	dur := 100 * time.Millisecond
	start := time.Now()
	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	n := 0
	go func() {
		for {
			<-ticker.C
			n++

			var sleep time.Duration
			if n >= 5 {
				sleep = 180 * time.Millisecond
				n = 0
			} else {
				sleep = 10 * time.Millisecond
			}
			fmt.Printf("Tick at %d, delaying for %d msecs\n", time.Since(start).Milliseconds(), sleep.Milliseconds())
			time.Sleep(sleep)
		}
	}()
	select {}
}
