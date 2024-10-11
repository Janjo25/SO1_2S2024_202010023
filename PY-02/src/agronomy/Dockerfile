FROM golang:1.20 AS builder
LABEL authors="Luis Lizama"

WORKDIR /app

# Copia el archivo 'go.mod' y 'go.sum' y descarga dependencias.
COPY go.mod go.sum ./
RUN go mod download

# Copia el código fuente al contenedor.
COPY . .

# Compila la aplicación.
RUN go build -o agronomy .

# Crea una imagen mínima para ejecutar la aplicación.
FROM debian:bookworm-slim

# Instala las versiones necesarias de GLIBC.
RUN apt-get update && apt-get install -y libc6

# Copia el binario compilado al contenedor.
# La ruta que se ha especificado es el lugar estándar que se usa para los programas que deseamos ejecutar globalmente.
COPY --from=builder /app/agronomy /usr/local/bin/agronomy

CMD ["agronomy"]
