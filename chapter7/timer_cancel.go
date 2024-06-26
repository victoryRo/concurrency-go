package chapter7

import (
	"fmt"
	"time"
)

func MainTimerCancel() {
	// timer := time.NewTimer(100 * time.Millisecond)
	timeout := make(chan struct{})

	go func() {
		// <-timer.C
		// close(timeout)
		// fmt.Println("Timeout")
		// ----------------------
		time.AfterFunc(200*time.Millisecond, func() {
			close(timeout)
			fmt.Println("TIMEOUT")
		})
	}()

	x := 0
	done := false
	for !done {
		select {
		case <-timeout:
			done = true
		default:
		}
		time.Sleep(time.Millisecond)
		x++
	}
	fmt.Println(x)
}
