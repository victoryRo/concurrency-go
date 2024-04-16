package chapter6


// Sólo puedes recuperar el pánico en la rutina si se inicia.
// Eso significa que si inicias una rutina que puede iniciar un pánico y no quieres que ese pánico termine el programa,
// tienes que recupera

go func(errCh chan<- error, resultCh chan<- result) {
     defer func() {
          if err:=recover(); err!=nil {
               // panic recovered, return error instead
               errCh <- err
          close(resultCh)
          }
     }()
     // Do the work
}()

//  Esta rutina toma dos canales como argumentos:
// 1. `errCh chan<- error`: 
// este es un canal al cual la función puede enviar valores de tipo `error`, 
// típicamente usado para transmitir errores ocurridos durante alguna operación.
// 2. `resultCh chan<- resultado`: 
// este es un canal al cual la función puede enviar valores de tipo `resultado`, 
// generalmente usado para transferir el resultado de una tarea. 
// Dentro de la rutina, hay una función diferida que está configurada para ejecutarse al final de la función envolvente. 
// Esta función maneja la recuperación de pánico, es decir, 
// detecta cualquier excepción no controlada que podría haber causado que la función entre en pánico y detenga la ejecución. 
// La primera línea de esta función diferida intenta `recuperar` cualquier pánico: 
// - Si tiene éxito (es decir, `err! = nil`), envía el error a través del canal `errCh`, 
// indicando efectivamente que se ha producido un pánico (y por lo tanto un error). 
// Luego cierra el canal `resultCh`. 
// Esta es una forma común de indicar que no se enviarán más datos en este canal y que es seguro dejar de leerlo. 
// El comentario `// Do the work` es un marcador de posición que indica dónde debe estar la lógica central de esta función. . 
// Esta es la sección que podría entrar en pánico y, por lo tanto, 
// está protegida por la función de recuperación diferida en la parte superior. 
// Este patrón se usa generalmente en Go para manejar tareas asincrónicas, 
// que se realizan en paralelo con otras y sus resultados y errores se administran por separado a través de canales. 
// El uso de "recuperar" dentro de un "diferir" es un modismo común para manejar los pánicos en tiempo de ejecución de forma controlada, 
// sin bloquear todo el programa.
