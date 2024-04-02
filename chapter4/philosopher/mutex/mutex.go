package mutex

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func philosopher(index int, leftFork, rightFork *sync.Mutex) {
	for {
		fmt.Printf("Philosopher %d is thinking\n", index)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

		leftFork.Lock()
		fmt.Printf("Philosopher %d got left fork\n", index)

		if rightFork.TryLock() {
			fmt.Printf("Philosopher %d got right fork\n", index)
			// eat
			fmt.Printf("Philosopher %d is eating\n", index)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			rightFork.Unlock()
		}

		leftFork.Unlock()
	}
}

func LocalMain() {
	forks := [5]sync.Mutex{}
	go philosopher(0, &forks[4], &forks[0])
	go philosopher(1, &forks[0], &forks[1])
	go philosopher(2, &forks[1], &forks[2])
	go philosopher(3, &forks[2], &forks[3])
	go philosopher(4, &forks[3], &forks[4])
	select {}
}
