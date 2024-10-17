/*
En un archivo 'build.rs', se puede escribir instrucciones que se ejecuten antes de la compilación del proyecto. Por
ejemplo, se puede generar archivos de configuración, procesar archivos que se necesitan en el proyecto, o establecer
variables de entorno. En general, 'build.rs' permite realizar cualquier operación que necesite ser ejecutada antes de la
compilación para garantizar que todo esté listo y configurado adecuadamente. Aunque 'build.rs' facilita la preparación
de todos los recursos necesarios, no siempre garantiza que el binario final sea completamente autónomo. Dependiendo del
proyecto, el binario podría seguir necesitando archivos externos, como configuraciones o recursos específicos, para su
correcta ejecución.
*/
fn main() {
    tonic_build::compile_protos("../../proto/discipline.proto").unwrap();
}
