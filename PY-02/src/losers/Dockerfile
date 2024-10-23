FROM golang:1.23 AS builder
LABEL authors="Luis Lizama"

WORKDIR /app

# Copia los archivos 'go.mod' y 'go.sum' y descarga las dependencias.
COPY ./src/losers/go.mod ./src/losers/go.sum ./src/losers/
RUN cd ./src/losers && go mod download

# Copia el código fuente.
COPY ./src/losers ./src/losers

# Compila la aplicación.
RUN cd ./src/losers && go build -o losers .

# Crea una imagen mínima para ejecutar la aplicación.
FROM debian:bookworm-slim

# Instala las versiones necesarias de GLIBC.
RUN apt-get update && apt-get install -y libc6

# Copia el binario.
COPY --from=builder /app/src/losers/losers /app/src/losers/losers

CMD ["./app/src/losers/losers"]
