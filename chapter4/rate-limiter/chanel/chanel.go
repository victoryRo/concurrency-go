package chanel

import (
	"fmt"
	"math/rand"
	"time"
)

type ChannelRate struct {
	bucket chan struct{}
	ticker *time.Ticker
	done   chan struct{}
}

func (c *ChannelRate) Wait() {
	<-c.bucket
}

func (c *ChannelRate) Close() {
	close(c.done)
	c.ticker.Stop()
}

func NewChannelRate(rate float64, limit int) *ChannelRate {
	ret := &ChannelRate{
		bucket: make(chan struct{}),
		ticker: time.NewTicker(time.Duration(1 / rate * 1000)),
		done:   make(chan struct{}),
	}

	for i := 0; i < limit; i++ {
		ret.bucket <- struct{}{}
	}

	go func() {
		for {
			select {
			case <-ret.done:
				return
			case <-ret.ticker.C:
				select {
				case ret.bucket <- struct{}{}:
				default:
				}
			}
		}
	}()

	return ret
}

func LocalMain() {
	limiter := NewChannelRate(5, 10)

	for i := 0; i < 100; i++ {
		limiter.Wait()
		fmt.Printf("Request: %v\n", time.Now())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(400)))
	}

	time.Sleep(time.Second * 2)

	for i := 0; i < 100; i++ {
		limiter.Wait()
		fmt.Printf("Request: %v\n", time.Now())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(400)))
	}
}
