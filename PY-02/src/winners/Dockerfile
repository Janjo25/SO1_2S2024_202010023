FROM golang:1.23 AS builder
LABEL authors="Luis Lizama"

WORKDIR /app

# Copia los archivos 'go.mod' y 'go.sum' y descarga las dependencias.
COPY ./src/winners/go.mod ./src/winners/go.sum ./src/winners/
RUN cd ./src/winners && go mod download

# Copia el código fuente.
COPY ./src/winners ./src/winners

# Compila la aplicación.
RUN cd ./src/winners && go build -o winners .

# Crea una imagen mínima para ejecutar la aplicación.
FROM debian:bookworm-slim

# Instala las versiones necesarias de GLIBC.
RUN apt-get update && apt-get install -y libc6

# Copia el binario.
COPY --from=builder /app/src/winners/winners /app/src/winners/winners

CMD ["./app/src/winners/winners"]
