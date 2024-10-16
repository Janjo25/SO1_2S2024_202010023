FROM golang:1.23 AS builder
LABEL authors="Luis Lizama"

WORKDIR /app

# Copia los archivos 'go.mod' y 'go.sum' y descarga las dependencias.
COPY ./proto/go.mod ./proto/
COPY ./src/disciplines/go.mod ./src/disciplines/go.sum ./src/disciplines/
RUN cd ./proto && go mod download
RUN cd ./src/disciplines && go mod download

# Copia el código fuente y los Protocol Buffers al contenedor.
COPY ./proto ./proto
COPY ./src/disciplines ./src/disciplines

# Compila la aplicación.
RUN cd ./src/disciplines && go build -o disciplines .

# Crea una imagen mínima para ejecutar la aplicación.
FROM debian:bookworm-slim

# Instala las versiones necesarias de GLIBC.
RUN apt-get update && apt-get install -y libc6

# Copia el binario y los Protocol Buffers al contenedor.
COPY --from=builder /app/src/disciplines/disciplines /app/src/disciplines/disciplines
COPY --from=builder /app/proto /app/proto

CMD ["./app/src/disciplines/disciplines"]
