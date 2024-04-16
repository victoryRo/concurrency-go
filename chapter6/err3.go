package chapter6

// una variación de esto está usando un canal de error en lugar de un canal de cancelación. 
// Una gorutina separada escucha el canal de error y captura los errores de las gorutinas:

// errCh will communicate errors
errCh := make(chan error)
// Any error will close canceled channel
canceled := make(chan struct{})
// Ensure error listener terminates
defer close(errCh)
// collect all errors.
errs := make([]error, 0)
go func() {
     once := sync.Once{}
     for err := range errCh {
           errs = append(errs, err)
           // cancel all goroutines when error received
           once.Do(func() { close(canceled) })
     }
}()
resultCh1 := make(chan Result1)
go func() {
    defer close(resultCh1)
     result, err := computeResult()
     if err != nil {
           errCh <- err
           // Make sure listener does not block
           return
     }
     // If canceled, stop
     select {
     case <-canceled:
           return
     default:
     }
     resultCh1 <- result
}()
result, ok := <-resultCh1


// Este fragmento de código crea coordinación entre múltiples gorutinas utilizando canales en Golang. 
// Aquí, se pueden ejecutar múltiples tareas en paralelo, 
// y si alguna tarea enfrenta un error, todas las demás tareas se detendrán. 
// Aquí hay un desglose del código: 
// - `errCh := make(chan error)`: 
// esta línea está inicializando un canal llamado `errCh` que se usará para propagar cualquier error que ocurra. 
// - `cancelado: = make(chan struct{})`: 
// Esto está inicializando un canal de estructura vacío llamado `cancelado`, 
// que se usará para señalar gorutinas si hay algún error. ocurren y necesitan detener la ejecución.
// - `errs := make([]error, 0)`: 
// Aquí estamos creando una porción vacía de errores que se usará para recopilar cualquier error que ocurra.
// - `defer close(errCh )`: 
// Esta línea cerrará el canal de error una vez que finalice la función circundante. 
// Esto se utiliza para garantizar que no se envíen más datos a este canal. 
// - `go func() {... }`: 
// esto crea una rutina (hilo concurrente en Go) que detecta errores en `errCh`. 
// Si se recibe un error, agrega este error al segmento `errs` y luego cancela todas las demás rutinas cerrando el canal `cancelado`. 
// - `resultCh1 := make(chan Result1)`: 
// aquí, el canal `resultCh1` es hecho para comunicar el resultado de algún cálculo (`Resultado1` es un tipo asumido). 
// - La siguiente gorutina maneja el cálculo del resultado y envía el resultado en `resultCh1`. 
// Si ocurre un error durante el cálculo, esto se envía a través de `errCh` y la función regresa para garantizar que no haya bloqueo. 
// - La instrucción `select` verifica si el canal `cancelado` está cerrado, 
// si es así, esta rutina se detiene y si no, envía el resultado calculado al canal `resultCh1`. 
// - Finalmente, `resultado, ok := <-resultCh1` 
// Lee el resultado de `resultCh1`. La variable `ok` será `false` si `resultCh1` se cierra antes de que recibamos un resultado. 
// Si "ok" es "verdadero", tenemos el "resultado" de la ejecución de la función.
