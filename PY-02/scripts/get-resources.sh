#!/bin/zsh

# Obtener los pods actuales.
echo "Verificando el estado de los pods..."
kubectl get pods

# Obtener los servicios actuales.
echo "Verificando el estado de los servicios..."
kubectl get services

# Obtener el Ingress.
echo "Verificando el estado del Ingress..."
kubectl get ingress

# Obtener los HPA actuales.
echo "Verificando el estado de los HPA..."
kubectl get hpa

echo "¡Los recursos actuales han sido obtenidos!"
