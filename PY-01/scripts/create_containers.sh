#!/bin/bash

IMAGES=("high-ram-image" "high-cpu-image" "low-ram-image" "low-cpu-image")

NUMBER_CONTAINERS=10 # Cantidad de contenedores que serán creados.

for i in $(seq 1 $NUMBER_CONTAINERS); do
    # 'IMAGES' es el arreglo que contiene los nombres de las imágenes que serán utilizadas para crear los contenedores.
    # '$RANDOM' genera un número aleatorio.
    # '${#IMAGES[@]}' obtiene la cantidad de elementos en el arreglo 'IMAGES'.
    # El operador '%' asegura que el número aleatorio generado esté dentro del rango de los índices del arreglo 'IMAGES'.
    IMAGE=${IMAGES[$RANDOM % ${#IMAGES[@]}]} # Selecciona una imagen aleatoria del arreglo 'IMAGES'.

    CONTAINER_NAME=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 8)

    CONTAINER_ID=$(docker run -d --name "$CONTAINER_NAME" "$IMAGE")
    echo "Se ha creado el contenedor $i: $CONTAINER_NAME utilizando la imagen '$IMAGE'"

    CONTAINER_PID=$(docker inspect --format '{{.State.Pid}}' "$CONTAINER_ID")

    echo "$CONTAINER_NAME-$CONTAINER_PID" >> ../kernel-module/containers_pid.txt
done
