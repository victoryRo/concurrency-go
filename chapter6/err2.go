package chapter6


// Separate result channels for goroutines
resultCh1 := make(chan Result1)
resultCh2 := make(chan Result2)
// canceled channel is closed once when a goroutine
// sends to cancelCh
canceled := make(chan struct{})
// cancelCh can receive many cancellation requests,
// but closes canceled channel once
cancelCh := make(chan struct{})
// Make sure cancelCh is closed, otherwise the
// goroutine that reads from it leaks
defer close(cancelCh)
go func() {
     // close canceled channel once when received 
     // from cancelCh
     once := sync.Once{}
     for range cancelCh {
           once.Do(func() {
                close(canceled)
           })
     }
}()
// Goroutine 1 computes Result1
go func() {
     result, err := computeResult1()
     if err != nil {
          // cancel other goroutines
           cancelCh <- struct{}{}
           // Send error back. Do not close channel
           resultCh1 <- Result1{Error: err}
           return
     }
     // If other goroutines failed, stop computation
     select {
          case <-canceled:
               // close resultCh1, so the listener does 
               // not block
               close(resultCh1)
               return
     default:
     }
     // Do more computations
}()
// Goroutine 2 computes Result2
go func() {
   ...
}()
// Receive results. The channel will be closed if
// the goroutine was canceled (ok will be false)
result1, ok1 := <-resultCh1
result2, ok2 := <-resultCh2
