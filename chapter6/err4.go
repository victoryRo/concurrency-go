package chapter6

// El enfoque de manejo de errores que veo a menudo en el campo es el uso de variables de error dedicadas 
// en el alcance adjunto para cada rutina. Este enfoque necesita un WaitGroup, 
// y no hay forma de cancelar el trabajo cuando falla una de las rutinas. 
// Sin embargo, puede resultar útil si ninguna de las rutinas realiza operaciones cancelables. 
// Si termina implementando este patrón, 
// asegúrese de que los errores se lean después de la Wait() llamada del grupo de espera porque, 
// según el modelo de memoria Go, 
// la configuración de las variables de error ocurre antes del retorno de esa Wait() llamada, 
// pero son concurrentes hasta entonces:

wg := sync.WaitGroup{}
wg.Add(2)
var err1 error
go func() {
     defer wg.Done()
     if err := doSomething1(); err != nil {
           err1 = err
           return
     }
}()
var err2 error
go func() {
     defer wg.Done()
     if err := doSomething2(); err != nil {
           err2 = err
           return
     }
}()
wg.Wait()
// Collect results and deal with errors here
if err1 != nil {
     // handle err1
}
if err2 != nil {
     // handle err2
}
