package prosumer

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func producer(index int, wg *sync.WaitGroup, done <-chan struct{}, output chan<- int) {
	defer wg.Done()

	for {
		value := rand.Int()
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

		select {
		case output <- value:
		case <-done:
			return
		}
		fmt.Printf("Producer %d sent %d\n", index, value)
	}
}

func consumer(index int, wg *sync.WaitGroup, input <-chan int) {
	defer wg.Done()

	for value := range input {
		fmt.Printf("Consumer %d received %d\n", index, value)
	}
}

func LocalMain() {
	doneCh := make(chan struct{})
	dataCh := make(chan int)

	producers := sync.WaitGroup{}
	consumers := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		producers.Add(1)
		go producer(i, &producers, doneCh, dataCh)
	}

	for i := 0; i < 10; i++ {
		consumers.Add(1)
		go consumer(i, &consumers, dataCh)
	}

	// select {}
	time.Sleep(time.Second * 2)
	close(doneCh)
	producers.Wait()
	close(dataCh)
	consumers.Wait()
}
