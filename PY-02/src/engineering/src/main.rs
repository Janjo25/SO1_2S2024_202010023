use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::routing::{get, post};
use axum::{Json, Router};
use serde::Deserialize;
use std::net::SocketAddr;
use tonic::Request;
use tracing::{error, info};

/*
Este módulo 'proto' incluye todas las definiciones generadas automáticamente a partir del archivo '.proto'. Contiene los
servicios y mensajes necesarios para la comunicación gRPC entre el cliente y el servidor, lo cual permite utilizar las
funcionalidades definidas en el archivo '.proto' dentro del código de Rust.
*/
pub mod proto {
    tonic::include_proto!("discipline");
}

use crate::proto::discipline_service_client::DisciplineServiceClient;
use crate::proto::DisciplineRequest;

#[derive(Clone, Deserialize)]
struct FacultyRequest {
    name: String,
    age: i32,
    faculty: String,
    discipline: i32,
}

async fn assign_student(
    request: FacultyRequest, server_address: String
) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    // Crea un cliente gRPC para conectarse al servicio de disciplinas, con 10 segundos de espera.
    let mut client = tokio::time::timeout(
        std::time::Duration::from_secs(10),
        DisciplineServiceClient::connect(server_address),
    ).await??;

    // Crea la solicitud gRPC con los datos del estudiante.
    let grpc_request = DisciplineRequest {
        name: request.name.clone(),
        age: request.age,
        faculty: request.faculty.clone(),
        discipline: request.discipline,
    };

    // Envía la solicitud gRPC al servicio de disciplinas.
    let response = client.assign(Request::new(grpc_request)).await?;

    // Verifica si la asignación fue exitosa.
    if response.into_inner().success {
        info!(
            "El estudiante '{}' ha sido asignado a la disciplina '{}'",
            request.name, request.discipline
        );
    } else {
        info!(
            "No se ha podido asignar al estudiante '{}' a la disciplina '{}'",
            request.name, request.discipline
        );
    }

    Ok(())
}

async fn request_handler(Json(faculty_request): Json<FacultyRequest>) -> impl IntoResponse {
    let server_address = "http://disciplines-service:80".to_string();

    // Ejecuta la tarea de asignar al estudiante a la disciplina de manera asíncrona.
    match tokio::spawn(assign_student(faculty_request.clone(), server_address)).await {
        Ok(Ok(_)) => {
            let message = format!(
                "Estudiante '{}' asignado a la disciplina '{}'",
                faculty_request.name, faculty_request.discipline
            );
            info!("{}", message);
            (StatusCode::OK, message)
        }
        Ok(Err(error)) => {
            error!("Ocurrió un error al asignar al estudiante: {}", error);
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                "Ocurrió un error al asignar al estudiante".to_string(),
            )
        }
        Err(error) => {
            error!("Ocurrió un error al ejecutar la tarea: {}", error);
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                "Ocurrió un error al procesar la solicitud".to_string(),
            )
        }
    }
}

async fn health_check_handler() -> impl IntoResponse {
    (StatusCode::OK, "OK")
}

#[tokio::main]
async fn main() {
    // Inicializa el logger.
    tracing_subscriber::fmt::init();

    // Define las rutas y los manejadores de las peticiones.
    let app = Router::new()
        .route("/engineering", post(request_handler))
        .route("/engineering/healthz", get(health_check_handler));

    // Define la dirección en la que el servidor escuchará las peticiones.
    let address = SocketAddr::from(([0, 0, 0, 0], 8080));
    info!("Servidor de Ingeniería escuchando en el puerto {}...", address);

    // Crea el servidor y lo pone en funcionamiento.
    axum_server::Server::bind(address).serve(app.into_make_service()).await.unwrap();
}
