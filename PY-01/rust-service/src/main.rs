use chrono::Utc;
use reqwest::blocking::Client;
use reqwest::Error;
use serde::{Deserialize, Serialize};
use serde_json::json;
use std::collections::HashMap;
use std::fs::{File, OpenOptions};
use std::io::{self, BufWriter, Read, Write};
use std::process::{Command, Stdio};
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use std::thread;
use std::time::Duration;

#[derive(Clone, Debug, Deserialize, Serialize)]
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

fn create_cron_job() -> io::Result<()> {
    let new_job = "* * * * * /home/luis-lizama/CLionProjects/SO1_2S2024_202010023/PY-01/scripts/create_containers.sh\n";

    let mut process = Command::new("crontab").stdin(Stdio::piped()).spawn()?;

    let stdin = process.stdin.as_mut().ok_or(io::Error::new(io::ErrorKind::Other, "Error al abrir la entrada estándar del proceso."))?;
    stdin.write_all(new_job.as_bytes())?;

    process.wait()?;

    Ok(())
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

    Ok(())
}

fn read_file(path: &str) -> io::Result<SysInfo> {
    let mut file = File::open(path)?; // Se abre el archivo.
    let mut contents = String::new(); // Se crea un buffer para almacenar el contenido del archivo.

    file.read_to_string(&mut contents)?; // Se lee el contenido del archivo y se almacena en el buffer.

    let system_information: SysInfo = serde_json::from_str(&contents)?; // Se deserializa el contenido JSON.

    // println!("{:#?}", system_information);

    Ok(system_information)
}

fn sort_and_select_processes(processes: &[Process]) -> (HashMap<String, u64>, HashMap<String, u64>) {
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

    let mut high_resource_processes: HashMap<String, u64> = HashMap::new();

    for process in sorted_processes.iter().take(2) {
        high_resource_processes.insert(process.name.clone(), process.pid);
    }

    sorted_processes.sort_by(|a, b| {
        a.cpu_usage.partial_cmp(&b.cpu_usage).unwrap().then_with(|| a.memory_usage.partial_cmp(&b.memory_usage).unwrap())
    });

    let mut low_resource_processes: HashMap<String, u64> = HashMap::new();

    for process in sorted_processes.iter().take(3) {
        low_resource_processes.insert(process.name.clone(), process.pid);
    }

    (high_resource_processes, low_resource_processes)
}

fn kill_container(name: &str) -> io::Result<()> {
    let output = Command::new("docker").args(&["rm", "-f", &name]).output()?;

    if !output.status.success() {
        let error_message = String::from_utf8_lossy(&output.stderr);

        return Err(io::Error::new(io::ErrorKind::Other, format!("Error al eliminar el contenedor: {}", error_message)));
    }

    Ok(())
}

fn send_log(process: &Process) -> Result<(), Error> {
    let client = Client::new();
    let timestamp = Utc::now().to_rfc3339();

    let log_message = json!({
        "timestamp": timestamp,
        "process": process
    });

    let response = client.post("http://localhost:80/log").json(&log_message).send()?;

    if !response.status().is_success() {
        eprintln!("La respuesta del servidor no fue exitosa: {}", response.status());
    }

    Ok(())
}

fn overwrite_file(high: HashMap<String, u64>, low: HashMap<String, u64>) -> io::Result<()> {
    let path = "../kernel-module/containers_pid.txt";
    let file = OpenOptions::new().write(true).truncate(true).open(path)?;
    let mut writer = BufWriter::new(file);

    for (name, pid) in high.iter() {
        writeln!(writer, "{}-{}", name, pid)?;
    }

    for (name, pid) in low.iter() {
        writeln!(writer, "{}-{}", name, pid)?;
    }

    writer.flush()?;

    Ok(())
}

fn delete_cron_job() -> io::Result<()> {
    let mut process = Command::new("crontab").stdin(Stdio::piped()).spawn()?;

    let stdin = process.stdin.as_mut().ok_or(io::Error::new(io::ErrorKind::Other, "Error al abrir la entrada estándar del proceso."))?;
    stdin.write_all(b"")?;

    process.wait()?;

    Ok(())
}

fn generate_graphs() -> Result<(), Error> {
    let client = Client::new();
    let response = client.get("http://localhost:80/generate-processes-graph").send()?;

    if !response.status().is_success() {
        eprintln!("La respuesta del servidor no fue exitosa: {}", response.status());
    }

    Ok(())
}

fn main() -> io::Result<()> {
    let running = Arc::new(AtomicBool::new(true));
    let running_clone = Arc::clone(&running);

    ctrlc::set_handler(move || {
        running_clone.store(false, Ordering::SeqCst);
    }).expect("Ocurrió un error al configurar el controlador 'Ctrl+C'.");

    if let Err(error) = create_cron_job() {
        eprintln!("Ocurrió un error al crear el trabajo programado: {:?}", error);
    }

    while running.load(Ordering::SeqCst) {
        execute_exporter()?;

        let path = "./sysinfo.json";
        let sysinfo = read_file(path)?;

        let (high, low) = sort_and_select_processes(&sysinfo.processes);

        println!("Total de RAM: {} KB", sysinfo.total_ram);
        println!("RAM libre: {} KB", sysinfo.free_ram);
        println!("RAM usada: {} KB", sysinfo.used_ram);

        println!("Procesos de alto consumo:");
        for process in &high {
            println!("{:?}", process);
        }

        println!("Procesos de bajo consumo:");
        for process in &low {
            println!("{:?}", process);
        }

        println!();

        let mut handles = vec![]; // En este vector se guardarán todos los hilos que se creen.

        for process in sysinfo.processes {
            if !high.contains_key(&process.name) && !low.contains_key(&process.name) {
                let process_name = process.name.clone(); // Se clona para que el código principal no pierda el 'ownership'.

                let handle = thread::spawn(move || {
                    println!("Eliminando contenedor con nombre '{}'", process_name);

                    if let Err(error) = kill_container(&process_name) {
                        eprintln!("Ocurrió un error al eliminar el contenedor '{}': {}", process_name, error);
                    } else {
                        if let Err(error) = send_log(&process) {
                            eprintln!("Ocurrió un error al enviar el registro del contenedor '{}': {}", process_name, error);
                        }
                    }
                });

                handles.push(handle);
            }
        }

        /* Este bucle espera a que todos los hilos terminen y maneja posibles errores. */
        for handle in handles {
            if let Err(error) = handle.join() {
                eprintln!("Ocurrió un error al esperar a que el hilo termine: {:?}", error);
            }
        }

        overwrite_file(high, low)?;

        println!("Esperando 10 segundos...");

        thread::sleep(Duration::from_secs(10));
    }

    if let Err(error) = delete_cron_job() {
        eprintln!("Ocurrió un error al eliminar el trabajo programado: {:?}", error);
    }

    if let Err(error) = generate_graphs() {
        eprintln!("Ocurrió un error al generar las gráficas: {:?}", error);
    }

    println!("Saliendo del programa...");

    Ok(())
}
