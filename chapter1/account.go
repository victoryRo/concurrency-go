package chapter1

import (
	"errors"
	"fmt"
	"sync"
)

type AccountOne struct {
	Name    string
	Balance int
}

type AccountTwo struct {
	Name    string
	Balance int
}

var ErrInsufficient = errors.New("error: insufficient funds")

func TransferOne(from, to *AccountOne, amt int) {
	if from.Balance < amt {
		fmt.Printf("⛔: %s\n", fmt.Sprintf("%v %v", from, to))
		return
	}

	from.Balance -= amt
	to.Balance += amt
	fmt.Printf("✅: %s\n", fmt.Sprintf("%v %v", from, to))
}

func TransferTwo(from, to *AccountTwo, amt int) {
	if from.Balance < amt {
		fmt.Printf("⛔: %s\n", fmt.Sprintf("%v %v", from, to))
		return
	}

	from.Balance -= amt
	to.Balance += amt
	fmt.Printf("✅: %s\n", fmt.Sprintf("%v %v", from, to))
}

// ----------------------------------------------------------------------

func ExecuteOne() {
	wg := sync.WaitGroup{}
	mx := sync.Mutex{}
	wg.Add(2)

	money := []int{5, 6}
	maria := AccountOne{Name: "Maria", Balance: 10}
	josefa := AccountOne{Name: "Josefa", Balance: 11}

	for _, v := range money {
		go func(value int) {
			mx.Lock()
			TransferOne(&maria, &josefa, value)
			mx.Unlock()
			wg.Done()
		}(v)
	}

	wg.Wait()
}

func ExecuteTwo() {
	wg := sync.WaitGroup{}
	mx := sync.Mutex{}
	wg.Add(2)

	maria := AccountTwo{Name: "Maria", Balance: 10}
	josefa := AccountTwo{Name: "Josefa", Balance: 15}

	go func() {
		mx.Lock()
		TransferTwo(&maria, &josefa, 5)
		mx.Unlock()
		wg.Done()
	}()
	go func() {
		mx.Lock()
		TransferTwo(&josefa, &maria, 10)
		mx.Unlock()
		wg.Done()
	}()

	wg.Wait()
}
