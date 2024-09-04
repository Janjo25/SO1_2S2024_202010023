use serde::Deserialize;
use std::fs::File;
use std::io::{self, Read};
use std::process::Command;
use std::thread::sleep;
use std::time::Duration;

#[derive(Clone, Debug, Deserialize)]
struct Process {
    pid: u64,
    name: String,
    cmd_line: String,
    vsz: u64,
    rss: u64,
    memory_usage: f64,
    cpu_usage: f64,
}

#[derive(Debug, Deserialize)]
struct SysInfo {
    total_ram: u64,
    free_ram: u64,
    used_ram: u64,
    processes: Vec<Process>,
}

fn execute_exporter() -> io::Result<()> {
    let working_directory = "../kernel-module";
    let exporter_path = "../kernel-module/exporter";
    let output = Command::new(exporter_path).current_dir(working_directory).output()?;

    if !output.status.success() {
        return Err(io::Error::new(io::ErrorKind::Other, "El exportador no se ejecutó correctamente."));
    }

    if !output.stderr.is_empty() {
        let error_message = String::from_utf8_lossy(&output.stderr);

        return Err(io::Error::new(io::ErrorKind::Other, format!("Error en el exportador: {}", error_message)));
    }

    let output_message = String::from_utf8_lossy(&output.stdout);
    println!("{}", output_message);

    sleep(Duration::from_secs(5));

    Ok(())
}

fn read_file(path: &str) -> io::Result<SysInfo> {
    let mut file = File::open(path)?; // Se abre el archivo.
    let mut contents = String::new(); // Se crea un buffer para almacenar el contenido del archivo.

    file.read_to_string(&mut contents)?; // Se lee el contenido del archivo y se almacena en el buffer.

    let system_information: SysInfo = serde_json::from_str(&contents)?; // Se deserializa el contenido JSON.

    println!("{:#?}", system_information);

    Ok(system_information)
}

fn sort_and_select_processes(processes: &[Process]) -> (Vec<Process>, Vec<Process>) {
    /* Ordenar por uso de CPU descendente. Si dos procesos tienen el mismo uso de CPU, se procede a ordenarlos por su memoria. */
    let mut sorted_processes = processes.to_vec(); // Se copian los procesos de un slice a un vector, ya que un slice no puede ser mutado.

    /* El orden en el que se colocan las letras 'a' y 'b' determinará el orden ascendente o descendente.
     * Si se coloca 'b' antes de 'a', se ordenará de forma descendente. Para ordenar de forma ascendente se hace lo contrario.
     * Se usa 'partial_cmp', ya que en el vector hay flotantes y existe la posibilidad de valores que no estén definidos.
     * El ordenamiento en Rust se realiza en dos partes.
       * Primero se realiza la fase de comparación, donde se determina si un elemento es mayor, menor o igual a otro. Rust toma nota de esto.
       * Luego, se realiza la fase de ordenamiento, donde se ordenan los elementos según la comparación realizada anteriormente.
     * Se usa 'then_with' para realizar una segunda comparación en caso de que la primera sea igual. */
    sorted_processes.sort_by(|a, b| {
        b.cpu_usage.partial_cmp(&a.cpu_usage).unwrap().then_with(|| b.memory_usage.partial_cmp(&a.memory_usage).unwrap())
    });

    /* Se toman los dos primeros procesos que más recursos consumen.
     * Se usa 'iter' para crear un iterador sobre el vector de procesos.
     * Se usa 'take' para tomar los dos primeros elementos del iterador (solamente se han seleccionado, pero no se ha hecho nada más).
     * Se usa 'cloned' para clonar los elementos seleccionados a un nuevo iterador.
     * Se usa 'collect' para consumir el iterador y almacenar los elementos en un nuevo vector. */
    let high_resource_processes: Vec<Process> = sorted_processes.iter().take(2).cloned().collect();

    sorted_processes.sort_by(|a, b| {
        a.cpu_usage.partial_cmp(&b.cpu_usage).unwrap().then_with(|| a.memory_usage.partial_cmp(&b.memory_usage).unwrap())
    });

    let low_resource_processes: Vec<Process> = sorted_processes.iter().take(3).cloned().collect();

    (high_resource_processes, low_resource_processes)
}

fn main() -> io::Result<()> {
    execute_exporter()?;

    let path = "./sysinfo.json";
    let sysinfo = read_file(path)?;

    let (high, low) = sort_and_select_processes(&sysinfo.processes);

    println!("Procesos de alto consumo:");
    for process in high {
        println!("{:?}", process);
    }

    println!("Procesos de bajo consumo:");
    for process in low {
        println!("{:?}", process);
    }

    Ok(())
}
