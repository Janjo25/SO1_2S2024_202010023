use std::io::Write;
use std::net::TcpListener;

fn main() {
    let listener = TcpListener::bind("0.0.0.0:8080").unwrap();

    println!("Servidor corriendo en el puerto 8080");

    for stream in listener.incoming() {
        let mut stream = stream.unwrap();
        let response = "HTTP/1.1 200 OK\r\n\r\n¡Hola, Olimpiadas de la USAC! Esta es la facultad de Ingeniería 👷‍♂️⚙️";
        stream.write(response.as_bytes()).unwrap();
        stream.flush().unwrap();
    }
}
