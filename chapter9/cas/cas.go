package cas

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type SomeStruct struct {
	v int
}

var sharedValue atomic.Pointer[SomeStruct]

func computedNewCopy(in SomeStruct) SomeStruct {
	return SomeStruct{v: in.v + 1}
}

func updateSharedValue(index int) {
	myCopy := sharedValue.Load()
	newCopy := computedNewCopy(*myCopy)

	if sharedValue.CompareAndSwap(myCopy, &newCopy) {
		fmt.Printf("Set value %d\n", index)
	} else {
		fmt.Printf("Cannot set value %d\n", index)
	}
}

func LocalCas() {
	sharedValue.Store(&SomeStruct{})
	wg := sync.WaitGroup{}

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			updateSharedValue(num)
		}(i)
	}
	wg.Wait()
}
