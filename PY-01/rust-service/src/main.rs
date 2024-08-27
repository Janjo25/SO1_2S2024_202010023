use serde_json::Value;
use std::fs::File;
use std::io::{self, Read};

fn main() -> io::Result<()> {
    let path = "./sysinfo.json";
    let mut file = File::open(path)?; // Se abre el archivo.

    let mut contents = String::new(); // Se crea un buffer para almacenar el contenido del archivo.

    file.read_to_string(&mut contents)?; // Se lee el contenido del archivo y se almacena en el buffer.

    let data: Value = serde_json::from_str(&contents)?; // Se deserializa el contenido JSON.

    println!("{:#?}", data);

    Ok(())
}
