package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// ExampleOne()
	// ExampleTwo()
	// ExampleThree()
	RaceFree()
}

// Race-free use of atomic as a synchronization tool. The number of
// times the program will print 1 will be different at each run. It
// will never print 0.
func RaceFree() {
	for i := 0; i < 1000000; i++ {
		var done atomic.Bool
		var a int
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			a = 1
			done.Store(true)
		}()
		if done.Load() {
			fmt.Println(a)
		}
		wg.Wait()
	}
}

// ----------------------------

// ExampleThree Si el efecto de una escritura en memoria atómica
// se observa mediante una lectura atómica,
// entonces la escritura atómica ocurrió antes de la lectura atómica.
func ExampleThree() {
	var done atomic.Bool
	var a int

	go func() {
		a = 1
		done.Store(true)
	}()

	if done.Load() {
		fmt.Println(a)
	}
}

// Tenga en cuenta que aquí todavía hay una condición de carrera,
// pero no una carrera de memoria.
// Dependiendo del orden de ejecución de las declaraciones,
// la rutina principal puede ver o no done como true.
// Sin embargo, si la rutina principal ve done como true, entonces está garantizado que .a=1

// ----------------------------

// ExampleTwo Aquí es donde atómic marca la diferencia.
// El siguiente programa no tiene carreras race-free
func ExampleTwo() {
	var str atomic.Value
	var done atomic.Bool
	str.Store("")

	go func() {
		str.Store("Done !")
		done.Store(true)
	}()

	for !done.Load() {
	}
	fmt.Println(str.Load())
}

// ----------------------------

// ExampleOne Hay una carrera de memoria,
// porque las variables str y done se escriben en una rutina y
// se leen en otra sin sincronización explícita.
// Hay varias formas en que este programa puede comportarse
func ExampleOne() {
	var str string
	var done bool

	go func() {
		str = "Done!"
		done = true
	}()

	for !done {
	}
	fmt.Println(str)
}

// ----------------------------
