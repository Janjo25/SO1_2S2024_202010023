#!/bin/zsh

# Conectar Docker a Minikube.
echo "Conectando Docker con Minikube..."
eval $(minikube docker-env)

# Construir la imagen para Agronomía.
echo "Construyendo la imagen para Agronomía..."
docker build --no-cache -t agronomy:latest -f ../src/agronomy/Dockerfile ..

# Construir la imagen para Ingeniería.
echo "Construyendo la imagen para Ingeniería..."
docker build --no-cache -t engineering:latest -f ../src/engineering/Dockerfile ..

# Construir la imagen para las disciplinas.
echo "Construyendo la imagen para las disciplinas..."
docker build --no-cache -t disciplines:latest -f ../src/disciplines/Dockerfile ..

# Construir la imagen para el consumidor de los ganadores.
echo "Construyendo la imagen para el consumidor de los ganadores..."
docker build --no-cache -t winners:latest -f ../src/winners/Dockerfile ..

# Construir la imagen para el consumidor de los perdedores.
echo "Construyendo la imagen para el consumidor de los perdedores..."
docker build --no-cache -t losers:latest -f ../src/losers/Dockerfile ..

# Borrar los deployments actuales.
echo "Borrando deployments actuales..."
kubectl delete deployment agronomy-faculty-deployment engineering-faculty-deployment disciplines-deployment winners-deployment losers-deployment

# Borrar los servicios actuales.
echo "Borrando servicios actuales..."
kubectl delete service agronomy-faculty-service engineering-faculty-service disciplines-service

# Borrar el Ingress actual.
echo "Borrando el Ingress actual..."
kubectl delete ingress faculties-ingress

# Borrar el clúster de Kafka actual.
echo "Borrando el clúster de Kafka actual..."
kubectl delete -f ../kafka/kafka-cluster.yaml -n kafka

# Borrar el topic de Kafka actual.
echo "Borrando el topic de Kafka actual..."
kubectl delete -f ../kafka/kafka-topic.yaml -n kafka

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

# Aplicar el deployment para las disciplinas.
echo "Aplicando el deployment para las disciplinas..."
kubectl apply -f ../deployments/disciplines-deployment.yaml

# Aplicar los servicios para las disciplinas.
echo "Aplicando el servicio para las disciplinas..."
kubectl apply -f ../services/disciplines-service.yaml

# Aplicar el Ingress de las facultades.
echo "Aplicando el Ingress para las facultades..."
kubectl apply -f ../ingresses/faculties-ingress.yaml

# Aplicar el clúster de Kafka.
echo "Aplicando el clúster de Kafka..."
kubectl apply -f ../kafka/kafka-cluster.yaml -n kafka

# Aplicar el topic de los ganadores.
echo "Aplicando el topic de los ganadores..."
kubectl apply -f ../kafka/kafka-winners-topic.yaml -n kafka

# Aplicar el topic de los perdedores.
echo "Aplicando el topic de los perdedores..."
kubectl apply -f ../kafka/kafka-losers-topic.yaml -n kafka

# Aplicar el deployment del consumidor de los ganadores.
echo "Aplicando el deployment para el consumidor de los ganadores..."
kubectl apply -f ../deployments/winners-deployment.yaml

# Aplicar el deployment del consumidor de los perdedores.
echo "Aplicando el deployment para el consumidor de los perdedores..."
kubectl apply -f ../deployments/losers-deployment.yaml

sleep 15

# Verificar los pods.
echo "Verificando el estado de los pods..."
kubectl get pods

# Verificar los servicios.
echo "Verificando el estado de los servicios..."
kubectl get services

# Verificar el Ingress.
echo "Verificando el estado del Ingress..."
kubectl get ingress

# Verificar el clúster de Kafka.
echo "Verificando el estado del Kafka..."
kubectl get pods -n kafka
kubectl get services -n kafka

echo "¡El despliegue ha finalizado!"
