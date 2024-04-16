package chapter6

// cuando hablo sobre servidores, Hablo principalmente de su naturaleza orientada a solicitudes
// y no de sus características de comunicación. 
// Las solicitudes pueden provenir de una red a través de HTTP o gRPC o pueden provenir de la línea de comando. 
// Por lo general, cada solicitud se maneja en una rutina separada. 
// Por lo tanto, depende de la pila de manejo de solicitudes propagar errores significativos 
// que puedan usarse para generar una respuesta para el usuario. 
// Si ese usuario es otro programa (es decir, si hablamos de un servicio web, por ejemplo), 
// tiene sentido incluir un código de error y algún mensaje de diagnóstico. 
// Los errores estructurados son tu mejor amigo:

// Embed this error to all other structured errors that can be returned from the API
type Error struct {
    Code int
    HTTPStatus int
    DiagMsg string
}
// HTTPError extracts HTTP information from an error
type HTTPError interface {
   GetHTTPStatus() int
   GetHTTPMessage() string
}
func (e Error) GetHTTPStatus() int {
  return e.HTTPStatus
}
func (e Error) GetHTTPMessage() string {
   return fmt.Sprintf("%d: %s",e.Code,e.DiagMsg)
}
// Treat HTTPErrors and other unrecognized errors
//separately
func WriteError(w http.ResponseWriter, err error) {
   if e, ok:=err.(HTTPError); ok {
      http.Error(w,e.HTTPStatus(),e.HTTPMessage())
   } else {
      http.Error(w,http.InternalServerError,err.Error())
   }
}


// Básicamente, este enfoque proporciona una forma de manejar tanto errores específicos de HTTP (con información detallada)
// como un método general para otros tipos de errores desconocidos o no especificados.

// Las implementaciones de errores como la anterior le ayudarán a devolver errores significativos a sus usuarios, 
// para que puedan solucionar problemas comunes
