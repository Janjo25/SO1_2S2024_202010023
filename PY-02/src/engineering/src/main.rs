use proto::discipline_service_client::DisciplineServiceClient;
use proto::DisciplineRequest;
use tokio::time::{timeout, Duration};
use tonic::transport::Channel;
use tonic::Request;

/*
Este módulo 'proto' incluye todas las definiciones generadas automáticamente a partir del archivo '.proto'. Contiene los
servicios y mensajes necesarios para la comunicación gRPC entre el cliente y el servidor, lo cual permite utilizar las
funcionalidades definidas en el archivo '.proto' dentro del código de Rust.
*/
pub mod proto {
    tonic::include_proto!("discipline");
}

async fn assign_student(client: &mut DisciplineServiceClient<Channel>, student_id: i32, discipline: &str) -> Result<(), Box<dyn std::error::Error>> {
    let duration = Duration::from_secs(10);
    let request = Request::new(DisciplineRequest {
        student_id,
        discipline: discipline.to_string(),
    });

    let response = timeout(duration, client.assign(request)).await??;

    if response.get_ref().success {
        println!("El estudiante con ID '{}' ha sido asignado a la disciplina '{}'", student_id, discipline);
    } else {
        println!("No se ha podido asignar al estudiante con ID '{}' a la disciplina '{}'", student_id, discipline);
    }

    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let student_id = 2;
    let discipline = "boxing";
    let server_address = "http://disciplines-service:80";

    let mut client = loop {
        match DisciplineServiceClient::connect(server_address).await {
            Ok(client) => break client,
            Err(_) => tokio::time::sleep(Duration::from_secs(5)).await,
        }
    };

    loop {
        if let Err(error) = assign_student(&mut client, student_id, discipline).await {
            eprintln!("Ocurrió un error al asignar al estudiante: {:?}", error);
        }

        tokio::time::sleep(Duration::from_secs(5)).await;
    }

    Ok(())
}
