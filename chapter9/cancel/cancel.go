package cancel

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func CancelSupport() (cancel func(), isCancelled func() bool) {
	v := atomic.Bool{}

	cancel = func() {
		v.Store(true)
	}
	isCancelled = func() bool {
		return v.Load()
	}

	return
}

// La CancelSupport función devuelve dos cierres:
// cancel() se puede llamar a la función para señalar la cancelación y
// isCancelled() se puede utilizar para verificar si se ha registrado una solicitud de cancelación.
// Ambos cierres comparten un bool valor atómico. Esto se puede utilizar de la siguiente manera:

func LocalCancel() {
	cancel, isCanceled := CancelSupport()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			time.Sleep(100 * time.Millisecond)
			if isCanceled() {
				fmt.Println("Cancelled")
				return
			}
		}
	}()

	time.AfterFunc(5*time.Second, cancel)
	wg.Wait()
}

// Cancelaciones
// Ya hemos visto cómo utilizar el cierre de un canal para señalar cancelaciones.
// Las implementaciones de contexto utilizan este paradigma para señalar cancelaciones y tiempos de espera.
// También se puede implementar un esquema de cancelación simple usando átomos:
