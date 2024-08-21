FROM alpine:latest
LABEL authors="luis-lizama"

# Se instala el paquete 'stress-ng' que permite realizar pruebas de estrés en la CPU.
RUN apk add --no-cache stress-ng

# Se ejecuta el comando 'stress-ng' para realizar pruebas de estrés en la CPU.
CMD ["stress-ng", "--cpu", "4", "--timeout", "600s"]
