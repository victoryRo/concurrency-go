package chapter7

import (
	"fmt"
	"time"
)

type TimerMockup struct {
	C chan<- time.Time
}

func NewTimerMockup(d time.Duration) *TimerMockup {
	t := &TimerMockup{
		C: make(chan time.Time, 1),
	}

	go func() {
		time.Sleep(d)
		t.C <- time.Now()
	}()

	return t
}

func ExecuteNewTimer() {
	t := NewTimerMockup(1 * time.Second)

	fmt.Println(t)
}
