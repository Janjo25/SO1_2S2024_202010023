FROM rust:1.80 AS builder
LABEL authors="Luis Lizama"

WORKDIR /app

# Copia el archivo 'Cargo.toml' y 'Cargo.lock' y descarga dependencias.
COPY Cargo.toml Cargo.lock ./
RUN cargo fetch

# Copia el código fuente al contenedor.
COPY . .

# Compila la aplicación.
RUN cargo build --release

# Crea una imagen mínima para ejecutar la aplicación.
FROM debian:bookworm-slim

# Instala las versiones necesarias de GLIBC.
RUN apt-get update && apt-get install -y libc6

# Copia el binario compilado al contenedor.
# La ruta que se ha especificado es el lugar estándar que se usa para los programas que deseamos ejecutar globalmente.
COPY --from=builder /app/target/release/engineering /usr/local/bin/engineering

CMD ["engineering"]
