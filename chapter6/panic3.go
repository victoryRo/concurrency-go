package chapter6

// un pánico Por lo general, indica un error en el programa,
// pero finalizar una canalización de larga duración después de horas de procesamiento no es la mejor solución.
// Por lo general, querrás tener un registro de todos los pánicos y errores una vez que se complete el procesamiento.
// Por lo tanto, debe asegurarse de que la recuperación del pánico se realice en el lugar correcto.
// Por ejemplo, en el siguiente fragmento de código,
// la recuperación de pánico se produce en torno a la función de procesamiento de la etapa de canalización real,
// por lo que se registra un pánico, pero el for bucle continúa procesándose:

func pipelineStage[IN any, OUT WithError](input <-chan IN, output chan<- OUT, errCh chan<-error, process func(IN) OUT) {
     defer close(output)
     for data := range input {
           // Process the next input
           result, err := func() (res OUT,err error) {
                defer func() {
                      // Convert panics to errors
                      if err = recover(); err != nil {
                           return
                      }
                }()
                return process(data),nil
           }()
           if err!=nil {
                // Report error and continue
                errCh<-err
                continue
           }
           output<-result
     }
}
