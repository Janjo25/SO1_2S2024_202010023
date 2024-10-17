FROM rust:1.81.0 AS builder
LABEL authors="Luis Lizama"

WORKDIR /app

# Copia el archivo 'Cargo.toml' y 'Cargo.lock' y descarga dependencias.
COPY ./src/engineering/Cargo.toml ./src/engineering/Cargo.lock ./src/engineering/
RUN cd ./src/engineering && cargo fetch

# Copia el código fuente y los Protocol Buffers al contenedor.
COPY ./proto ./proto
COPY ./src/engineering ./src/engineering

# Instala el compilador de Protocol Buffers.
RUN apt-get update && apt-get install -y protobuf-compiler

# Compila la aplicación.
RUN cd ./src/engineering && cargo build --release

# Crea una imagen mínima para ejecutar la aplicación.
FROM debian:bookworm-slim

# Copia el binario al contenedor.
COPY --from=builder /app/src/engineering/target/release/engineering /app/src/engineering/engineering

CMD ["./app/src/engineering/engineering"]
