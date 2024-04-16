package chapter6

// Veamos algunos patrones comunes. 
// Si envía trabajo a una rutina y espera recibir un resultado más adelante, 
// asegúrese de que ese resultado incluya información de error. 
// El patrón ilustrado aquí es útil si tiene varias tareas que pueden realizarse simultáneamente. 
// Comienza cada tarea en su propia rutina y luego recopila los resultados o errores según sea necesario. 
// Esta también es una buena forma de solucionar errores cuando tienes un grupo de trabajadores:

// Result type keeps the expected result, and the
// error information.
type Result1 struct {
     Result ResultType1
     Error err
}
type Result2 struct {
     Result ResultType2
     Error err
}
...
result1Ch:=make(chan Result1)    
go func() {
     result, err := handleRequest1()
     result1Ch <- Result1{ Result: result, Error: err }
}()
result2Ch:=make(chan Result2)    
go func() {
     result, err := handleRequest2()
     result2Ch <- Result2{ Result: result, Error: err }
}() 
// Do other work
...

// Collect the results from the goroutines
result1:=<-result1Ch
if result1.Error!=nil {
// result2Ch is never read. Goroutine leaks!
       return result1.Error
}
result2:=<-result2Ch
...
