package simple

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Limiter struct {
	mx         sync.Mutex
	rate       int
	bucketSize int
	nTokens    int
	lastToken  time.Time
}

func (l *Limiter) Wait() {
	l.mx.Lock()
	defer l.mx.Unlock()

	if l.nTokens > 0 {
		l.nTokens--
		return
	}

	tElapsed := time.Since(l.lastToken)
	period := time.Second / time.Duration(l.rate)
	nTokens := tElapsed.Nanoseconds() / period.Nanoseconds()

	l.nTokens = int(nTokens)
	if l.nTokens > l.bucketSize {
		l.nTokens = l.bucketSize
	}

	l.lastToken = l.lastToken.Add(time.Duration(nTokens) * period)
	if l.nTokens > 0 {
		l.nTokens--
		return
	}

	next := l.lastToken.Add(period)
	wait := next.Sub(time.Now())
	if wait >= 0 {
		time.Sleep(wait)
	}
	l.lastToken = next
}

func NewLimiter(rate, limit int) *Limiter {
	return &Limiter{
		rate:       rate,
		bucketSize: limit,
		nTokens:    limit,
		lastToken:  time.Now(),
	}
}

func Localmain() {
	limiter := NewLimiter(5, 10)

	for i := 0; i < 100; i++ {
		limiter.Wait()
		fmt.Printf("Request: %v %+v\n", time.Now(), limiter)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(400)))
	}

	time.Sleep(time.Second * 2)

	for i := 0; i < 100; i++ {
		limiter.Wait()
		fmt.Printf("Request: %v %+v\n", time.Now(), limiter)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(400)))
	}
}
