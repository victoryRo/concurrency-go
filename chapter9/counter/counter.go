package counter

import (
	"fmt"
	"sync/atomic"
)

var count int64

func LocalCounter() {
	for i := 0; i < 10000; i++ {
		go func() {
			atomic.AddInt64(&count, 1)
		}()
	}

	for {
		v := atomic.LoadInt64(&count)
		fmt.Println(v)
		if v == 10000 {
			break
		}
	}
}

// Contadores
// Los elementos atómicos se pueden utilizar como contadores eficientes y seguros para la concurrencia.
// El siguiente programa crea muchas goroutines,
// cada una de las cuales sumará 1 al contador compartido.
// Otra goroutine se repite hasta que el contador llega a 10000.
// Debido al uso de elementos atómicos aquí,
// este programa no tiene restricciones de carrera y siempre terminará imprimiendo 10000:
