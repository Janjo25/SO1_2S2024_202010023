FROM alpine:latest
LABEL authors="luis-lizama"

# Se instala el paquete 'stress-ng' que permite realizar pruebas de estrés en la memoria RAM.
RUN apk add --no-cache stress-ng

# Se ejecuta el comando 'stress-ng' para realizar pruebas de estrés en la memoria RAM.
CMD ["stress-ng", "--vm", "1", "--vm-bytes", "75%", "--vm-hang", "1", "--vm-keep", "-t", "600s"]
