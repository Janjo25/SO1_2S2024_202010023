# **Manual Técnico: Proyecto Olimpiadas USAC**

## **Introducción**

Este manual detalla la configuración e implementación de una arquitectura en Google Cloud Platform (GCP) utilizando
Google Kubernetes Engine (GKE). El proyecto busca monitorizar en tiempo real los resultados de las Olimpiadas de la
Universidad de San Carlos de Guatemala, donde los participantes de las facultades de Ingeniería y Agronomía competirán
en varias disciplinas. Utilizaremos Golang para el procesamiento concurrente y Kafka para el manejo de flujos de datos.
Grafana y Prometheus serán claves para la visualización de datos y monitoreo en tiempo real.
---

## **1. Creación del Cluster en GKE**

Para lograr una infraestructura robusta y escalable en la nube, creamos un clúster en Google Kubernetes Engine (GKE). El
cluster servirá como la base de nuestra arquitectura, permitiéndonos gestionar de manera eficiente los recursos, cargas
de trabajo, y el escalado automático de las aplicaciones según la demanda.

Ejecutamos el siguiente comando para crear el cluster:

```bash
gcloud container clusters create olympics-us-central1-a --disk-size=25GB --disk-type=pd-standard --num-nodes=5 --zone=us-central1-a
```

Luego, obtenemos las credenciales para poder interactuar con el cluster:

```bash
gcloud container clusters get-credentials <CLUSTER-NAME> --zone <ZONE> --project <PROJECT-ID>
```

---

## **2. Configuración de Herramientas y Componentes**

### **Locust (Generador de tráfico)**

Locust es una herramienta open-source basada en Python para realizar pruebas de carga y generar tráfico simulado. En
este proyecto, Locust nos permite replicar grandes volúmenes de peticiones que simulan la actividad de los participantes
en la plataforma, permitiéndonos probar la resistencia de nuestra arquitectura ante situaciones de alto tráfico.

Para lanzar Locust, ejecutamos:

```bash
python -m locust -f main.py
```

Este comando abre la interfaz web de Locust, desde donde podemos controlar y observar la generación de tráfico y su
efecto en nuestro sistema.

### **Strimzi & Kafka**

Apache Kafka es una plataforma de streaming distribuida ideal para manejar grandes volúmenes de datos en tiempo real. En
este proyecto, Kafka se utiliza para gestionar y enrutar los resultados de los participantes, clasificándolos en
ganadores y perdedores. Strimzi es una implementación de Kafka en Kubernetes que facilita su despliegue y gestión en
nuestro cluster.

Para configurar Kafka mediante Strimzi, primero creamos un namespace específico:

```bash
kubectl create namespace kafka
```

Luego, instalamos Strimzi usando Helm, lo cual simplifica y automatiza la instalación de Kafka:

```bash
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
```

### **Instalación de Helm**

Helm es un gestor de paquetes para Kubernetes, que facilita la implementación de aplicaciones y servicios al permitirnos
manejar configuraciones complejas de forma simplificada. Usaremos Helm para instalar Redis, Grafana y Prometheus en el
cluster, manteniendo la configuración estandarizada y repetible.

Para instalar Helm:

```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

Luego añadimos el repositorio de Bitnami, que contiene los charts oficiales de varias aplicaciones populares:

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

### **Redis**

Redis es una base de datos en memoria, diseñada para manejar almacenamiento temporal de datos a alta velocidad. En este
proyecto, Redis actúa como un almacenamiento intermedio donde se guardan temporalmente los resultados de los
participantes antes de que sean visualizados en Grafana.

Para desplegar Redis:

```bash
helm install redis-release bitnami/redis
```

Si necesitamos acceder a la contraseña de Redis, podemos obtenerla con:

```bash
kubectl get secret --namespace default redis-release -o jsonpath="{.data.redis-password}" | base64 -d
```

### **Grafana**

Grafana es una plataforma de visualización de datos que permite monitorear métricas en tiempo real. En nuestro proyecto,
Grafana utiliza los datos de Redis y las métricas de Prometheus para visualizar el conteo de participantes, las
disciplinas en las que compiten, y otros detalles de la infraestructura en un formato gráfico y dinámico.

Para implementar Grafana:

```bash
helm install grafana bitnami/grafana --set plugins="redis-datasource" --set service.type=LoadBalancer
```

Para acceder a Grafana y ver los paneles, podemos obtener la contraseña de administrador con:

```bash
echo "Password: $(kubectl get secret grafana-admin --namespace default -o jsonpath="{.data.GF_SECURITY_ADMIN_PASSWORD}" | base64 -d)"
```

### **Prometheus**

Prometheus es una solución de monitoreo de código abierto diseñada para recolectar métricas de los servicios en
Kubernetes. Nos ayuda a monitorear en tiempo real el rendimiento de los deployments, recursos del sistema, y el tráfico
de red, asegurando que nuestra infraestructura se mantenga estable y optimizada.

Para desplegar Prometheus en el cluster:

```bash
helm install prometheus bitnami/prometheus
```

Instalamos además exportadores que permiten la recolección de métricas de los recursos de Kubernetes:

```bash
helm install kube-state-metrics bitnami/kube-state-metrics
helm install node-exporter bitnami/node-exporter
```

Finalmente, configuramos Prometheus para rastrear servicios específicos editando el configmap:

```bash
kubectl edit configmap prometheus-server
```

Agregamos configuraciones para monitorizar servicios importantes:

```yaml
- job_name: 'kube-state-metrics'
  static_configs:
    - targets: [ '<CLUSTER-IP>:<PORT>' ]
- job_name: 'node-exporter'
  static_configs:
    - targets: [ '<CLUSTER-IP>:<PORT>' ]
- job_name: 'agronomy-metrics'
  static_configs:
    - targets: [ '<CLUSTER-IP>:<PORT>' ]
- job_name: 'engineering-metrics'
  static_configs:
    - targets: [ '<CLUSTER-IP>:<PORT>' ]
```

Luego reiniciamos Prometheus:

```bash
kubectl rollout restart deployment prometheus-server
```

---

## **3. Implementación de Deployments y Servicios**

### **Deployments de Facultades**

Cada facultad, Ingeniería y Agronomía, tiene su propio deployment con contenedores en Golang o Rust. Estos deployments
son responsables de recibir y procesar las solicitudes de los participantes, dirigiendo el tráfico a las disciplinas
correspondientes mediante gRPC, aprovechando la concurrencia de Golang (go routines) y el manejo de threads en Rust.

### **Deployments de Disciplinas**

Los contenedores para las disciplinas (Natación, Boxeo, Atletismo) están diseñados en Go. Cada servidor implementa un
algoritmo simple de probabilidad (como lanzar una moneda) para decidir si el alumno es ganador o perdedor. Kafka se usa
para enviar estos resultados a tópicos específicos, lo cual nos permite organizar el flujo de datos de una forma
eficiente y escalable.

### **Kafka**

Configuramos Kafka para gestionar los mensajes de ganadores y perdedores. Kafka organiza los datos en tópicos,
permitiendo que los consumidores accedan y procesen esta información de manera paralela sin duplicación de mensajes.

### **Redis, Grafana, y Prometheus**

Redis guarda los datos temporales para ser visualizados en Grafana, que muestra las métricas y resultados en tiempo
real. Prometheus monitorea el rendimiento y la estabilidad de los servicios, generando métricas importantes sobre el
estado del cluster.
