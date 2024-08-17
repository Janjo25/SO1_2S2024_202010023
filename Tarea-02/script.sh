#!/bin/bash

NUMBER_CONTAINERS=10 # Cantidad de contenedores que serán creados.

for i in $(seq 1 $NUMBER_CONTAINERS); do
    CONTAINER_NAME=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 8) # Genera un nombre aleatorio de 8 caracteres alfanuméricos para el contenedor.

    docker run -d --name "$CONTAINER_NAME" alpine sleep 3600 # Se crea un contenedor utilizando 'alpine' y se mantendrá en ejecución por 1 hora.

    echo "Se ha creado el contenedor $i: $CONTAINER_NAME"
done
