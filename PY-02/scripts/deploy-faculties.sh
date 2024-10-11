#!/bin/zsh

# Conectar Docker a Minikube.
echo "Conectando Docker con Minikube..."
eval $(minikube docker-env)

# Construir la imagen para Agronomía.
echo "Construyendo la imagen para Agronomía..."
docker build --no-cache -t agronomy:latest ../src/agronomy

# Construir la imagen para Ingeniería.
echo "Construyendo la imagen para Ingeniería..."
docker build --no-cache -t engineering:latest ../src/engineering

# Borrar los deployments actuales.
echo "Borrando deployments actuales..."
kubectl delete deployment agronomy-faculty-deployment engineering-faculty-deployment

# Aplicar el deployment de Agronomía.
echo "Aplicando el deployment para Agronomía..."
kubectl apply -f ../deployments/agronomy-faculty-deployment.yaml

# Aplicar el servicio de Agronomía.
echo "Aplicando el servicio para Agronomía..."
kubectl apply -f ../services/agronomy-faculty-service.yaml

# Aplicar el deployment de Ingeniería.
echo "Aplicando el deployment para Ingeniería..."
kubectl apply -f ../deployments/engineering-faculty-deployment.yaml

# Aplicar el servicio de Ingeniería.
echo "Aplicando el servicio para Ingeniería..."
kubectl apply -f ../services/engineering-faculty-service.yaml

sleep 5

# Verificar los pods.
echo "Verificando el estado de los pods..."
kubectl get pods

# Verificar los servicios.
echo "Verificando el estado de los servicios..."
kubectl get services

echo "¡El despliegue ha finalizado!"
