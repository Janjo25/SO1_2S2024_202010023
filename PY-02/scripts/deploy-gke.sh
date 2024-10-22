#!/bin/zsh

# Autenticar Docker con GCP.
echo "Autenticando Docker con Google Cloud..."
gcloud auth configure-docker

# Construir y subir la imagen para Agronomía.
echo "Construyendo y subiendo la imagen para Agronomía..."
docker build --no-cache -t gcr.io/<PROJECT-ID>/agronomy:latest -f ../src/agronomy/Dockerfile ..
docker push gcr.io/<PROJECT-ID>/agronomy:latest

# Construir y subir la imagen para Ingeniería.
echo "Construyendo y subiendo la imagen para Ingeniería..."
docker build --no-cache -t gcr.io/<PROJECT-ID>/engineering:latest -f ../src/engineering/Dockerfile ..
docker push gcr.io/<PROJECT-ID>/engineering:latest

# Construir y subir la imagen para las disciplinas.
echo "Construyendo y subiendo la imagen para las disciplinas..."
docker build --no-cache -t gcr.io/<PROJECT-ID>/disciplines:latest -f ../src/disciplines/Dockerfile ..
docker push gcr.io/<PROJECT-ID>/disciplines:latest

# Borrar los deployments actuales.
echo "Borrando deployments actuales..."
kubectl delete deployment agronomy-faculty-deployment engineering-faculty-deployment disciplines-deployment

# Borrar los servicios actuales.
echo "Borrando servicios actuales..."
kubectl delete service agronomy-faculty-service engineering-faculty-service disciplines-service

# Borrar el Ingress actual.
echo "Borrando el Ingress actual..."
kubectl delete ingress faculties-ingress

# Borrar los HPA actuales.
echo "Borrando los HPA actuales..."
kubectl delete hpa agronomy-faculty-hpa engineering-faculty-hpa

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

# Aplicar los deployments para las disciplinas.
echo "Aplicando el deployment para las disciplinas..."
kubectl apply -f ../deployments/disciplines-deployment.yaml

# Aplicar los servicios para las disciplinas.
echo "Aplicando el servicio para las disciplinas..."
kubectl apply -f ../services/disciplines-service.yaml

# Aplicar el Ingress de las facultades.
echo "Aplicando el Ingress para las facultades..."
kubectl apply -f ../ingresses/faculties-ingress.yaml

# Aplicar el HPA de Agronomía.
echo "Aplicando el HPA para Agronomía..."
kubectl apply -f ../hpa/agronomy-faculty-hpa.yaml

# Aplicar el HPA de Ingeniería.
echo "Aplicando el HPA para Ingeniería..."
kubectl apply -f ../hpa/engineering-faculty-hpa.yaml

sleep 5

# Verificar los pods.
echo "Verificando el estado de los pods..."
kubectl get pods

# Verificar los servicios.
echo "Verificando el estado de los servicios..."
kubectl get services

# Verificar el Ingress.
echo "Verificando el estado del Ingress..."
kubectl get ingress

# Verificar los HPA.
echo "Verificando el estado del HPA..."
kubectl get hpa

echo "¡El despliegue ha finalizado!"
