package signal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Localmain This program prints "Tick..." every second until the program
// terminates by a signal. When the signal is caught, it prints
// "Terminating".
func Localmain() {
	term := make(chan struct{})
	done := make(chan struct{})
	sig := make(chan os.Signal, 1)

	tick := time.NewTicker(time.Second / 2)

	go func() {
		<-sig
		close(term)
	}()

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() {
			tick.Stop()
			close(done)
		}()

		for {
			select {
			case <-term:
				fmt.Println("Terminating...")
				return
			case <-tick.C:
				fmt.Println("Tick...")
			}
		}
	}()
	<-done
}
