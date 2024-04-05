package chapter1

import (
	"fmt"
	"sync"
)

type Account struct {
	Name    string
	Balance int
}

func Transfer(from, to *Account, amt int) {
	if from.Balance < amt {
		fmt.Printf("⛔: %s\n", fmt.Sprintf("%v %v", from, to))
		return
	}

	from.Balance -= amt
	to.Balance += amt
	fmt.Printf("✅: %s\n", fmt.Sprintf("%v %v", from, to))
}

// ----------------------------------------------------------------------

func Execute() {
	wg := sync.WaitGroup{}
	mx := sync.Mutex{}
	wg.Add(2)

	money := []int{5, 6}
	maria := Account{Name: "Maria", Balance: 10}
	josefa := Account{Name: "Josefa", Balance: 11}

	for _, v := range money {
		go func(value int) {
			defer mx.Unlock()
			mx.Lock()
			Transfer(&maria, &josefa, value)
			wg.Done()
		}(v)
	}

	wg.Wait()
}
