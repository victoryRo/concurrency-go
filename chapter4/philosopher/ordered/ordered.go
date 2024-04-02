package ordered

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func philosopher(index int, firstFork, secondFork *sync.Mutex) {
	for {
		fmt.Printf("philosopher %d thinking\n", index)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

		firstFork.Lock()
		fmt.Printf("philosopher %d got left fork\n", index)
		secondFork.Lock()
		fmt.Printf("philosopher %d got right fork\n", index)

		fmt.Printf("Philosopher %d is eating\n", index) // Eat
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

		secondFork.Unlock()
		firstFork.Unlock()
	}
}

func Localmain() {
	forks := [5]sync.Mutex{}
	go philosopher(0, &forks[0], &forks[4])
	go philosopher(1, &forks[0], &forks[1])
	go philosopher(2, &forks[1], &forks[2])
	go philosopher(3, &forks[2], &forks[3])
	go philosopher(4, &forks[3], &forks[4])
	select {}
}
