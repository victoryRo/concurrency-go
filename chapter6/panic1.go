package chapter6

// En un programa de servidor, normalmente, una rutina separada maneja cada solicitud.
// La mayoría de los marcos de servidores (incluido el paquete net/http de la biblioteca estándar)
// manejan estos pánicos sin fallar al imprimir una pila y fallar la solicitud.
// Si está escribiendo un servidor sin utilizar dicha biblioteca o si desea proporcionar más información cuando le entra el pánico,
// debe encargarse de ello usted mismo:

func PanicHandler(next func(http.ResponseWriter,*http.Request)) func(http.ResponseWriter,*http.Request) {
  return func(wr http.ResponseWriter, req *http.Request) {
    defer func() {
        if err:=recover(); err!=nil {
           // print panic info
        }
    }()
    next(wr,req)
   }
}
func main() {
     http.Handle("/path",PanicHandler(pathHandler))
}

// El fragmento de código está escrito en lenguaje Go 
// y sirve para manejar situaciones de pánico que pueden ocurrir en una aplicación de servidor HTTP. 
// En Go, las situaciones de pánico son excepciones. 
// Cuando esto ocurre, detienen el flujo normal de control y comienzan a entrar en pánico. 
// Si una función encuentra un pánico, su ejecución se detendrá, 
// se ejecutarán todas las funciones diferidas y luego el pánico continúa en la pila de llamadas. 
// Si el pánico no se recupera, el programa falla. 
// Este fragmento incluye dos funciones, `PanicHandler` y `main`. `PanicHandler` 
// es en realidad una función que toma una función como argumento y devuelve otra función. 
// Este es un patrón común en Go llamado "middleware". 
// La función `siguiente` que espera tiene que ser una función de controlador que se ajuste al tipo `func(http.ResponseWriter,*http.Request)`. 
// Esta función devuelta es otro controlador HTTP. 
// Dentro de la función devuelta, 
// utiliza "diferir" para garantizar que la función anónima que la sigue se ejecute sin importar cómo continúe el flujo del programa. 
// Esta función anónima utiliza `recover()` para detectar un pánico si ocurre durante la ejecución de la función del controlador. 
// Si `recover()` detecta un pánico, devuelve el valor que se pasó a `panic()`. 
// Si no se produjo ningún pánico, devuelve nil. 
// En la función `main`, el controlador HTTP para "/path" se establece llamando a `http.Handle()`, 
// pasando el patrón de URL y `PanicHandler` como argumentos. 
// La función `pathHandler` no está definida en su fragmento, 
// pero se supone que maneja las solicitudes HTTP en "/path". `PanicHandler(pathHandler)` 
// se asegura de que la función `pathHandler` 
// esté protegida o envuelta con el mecanismo de recuperación de pánico definido en `PanicHandler`. 
// Por lo tanto, se puede recuperar cualquier pánico que pueda ocurrir dentro de `pathHandler`, 
// evitando que toda la aplicación falle. 
// En pocas palabras, este código se utiliza para crear un controlador HTTP de recuperación de pánico para "/ruta".
