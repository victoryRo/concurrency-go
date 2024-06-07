package progress

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type ProgressMeter struct {
	progress int64
}

func (pm *ProgressMeter) Progress() {
	atomic.AddInt64(&pm.progress, 1)
}

func (pm *ProgressMeter) Get() int64 {
	return atomic.LoadInt64(&pm.progress)
}

// -------------------------------------------------

func longGoroutine(ctx context.Context, pm *ProgressMeter) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Canceled")
			return
		default:
		}
		time.Sleep(time.Duration(rand.Intn(120)) * time.Millisecond)
		pm.Progress()
	}
}

func observer(ctx context.Context, cancel func(), progress *ProgressMeter) {
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	var lastProgress int64
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			p := progress.Get()
			if p == lastProgress {
				fmt.Println("No progress since last time, canceling")
				cancel()
				return
			}
			fmt.Printf("Progress %d\n", p)
			lastProgress = p
		}
	}
}

func LocalProgress() {
	var progresss ProgressMeter

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		longGoroutine(ctx, &progresss)
	}()
	go observer(ctx, cancel, &progresss)

	wg.Wait()
}

// Medidor de latidos y progreso
// A veces, una gorutina puede dejar de responder o no progresar tan rápido como sea necesario.
// un latido del corazón puede utilizar un medidor de progreso para observar tales gorutinas.
// Hay varias formas de hacer esto; por ejemplo,
// la rutina observada puede usar envíos sin bloqueo para anunciar el progreso,
// o puede anunciar su progreso incrementando una variable compartida protegida por un mutex.
// Los atómicos nos permiten implementar el esquema de variable compartida sin un mutex.
// Esto también tiene la ventaja de ser observable por múltiples gorutinas sin sincronización adicional.
