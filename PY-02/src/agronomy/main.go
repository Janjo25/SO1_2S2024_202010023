package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/janjo25/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

type PrometheusMetrics struct {
	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
}

type grpcClient struct {
	Connection *grpc.ClientConn
	Client     pb.DisciplineServiceClient
}

type FacultyRequest struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Faculty    string `json:"faculty"`
	Discipline int    `json:"discipline"`
}

func initializeMetrics() *PrometheusMetrics {
	metrics := &PrometheusMetrics{
		HttpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Número total de peticiones HTTP",
			}, []string{"method", "endpoint"},
		),
		HttpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duración de las peticiones HTTP en segundos",
				Buckets: prometheus.DefBuckets,
			}, []string{"method", "endpoint"},
		),
	}

	prometheus.MustRegister(metrics.HttpRequestsTotal)
	prometheus.MustRegister(metrics.HttpRequestDuration)

	return metrics
}

func newClient(serverAddress string) (*grpcClient, error) {
	connection, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("ocurrió un error al crear el cliente: %v", err)
	}

	client := pb.NewDisciplineServiceClient(connection)

	return &grpcClient{Connection: connection, Client: client}, nil
}

func assignStudent(metrics *PrometheusMetrics, client *grpcClient, request FacultyRequest) {
	// Iniciar el temporizador para medir la latencia.
	start := time.Now()
	metrics.HttpRequestsTotal.WithLabelValues("POST", "/agronomy").Inc()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcRequest := &pb.DisciplineRequest{
		Name:       request.Name,
		Age:        int32(request.Age),
		Faculty:    request.Faculty,
		Discipline: int32(request.Discipline),
	}

	grpcResponse, err := client.Client.Assign(ctx, grpcRequest)
	if err != nil {
		log.Printf("Ocurrió un error al llamar al servicio: %v", err)

		return
	}

	if grpcResponse.Success {
		log.Printf("El estudiante '%s' ha sido asignado a la disciplina '%d'", request.Name, request.Discipline)
	} else {
		log.Printf("No se ha podido asignar al estudiante '%s' a la disciplina '%d'", request.Name, request.Discipline)
	}

	// Registrar la duración de la petición.
	metrics.HttpRequestDuration.WithLabelValues("POST", "/agronomy").Observe(time.Since(start).Seconds())
}

func requestHandler(writer http.ResponseWriter, request *http.Request, metrics *PrometheusMetrics, client *grpcClient) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Solo se permiten peticiones POST", http.StatusMethodNotAllowed)

		return
	}

	var facultyRequest FacultyRequest

	err := json.NewDecoder(request.Body).Decode(&facultyRequest)
	if err != nil {
		http.Error(writer, "No se pudo decodificar el cuerpo de la petición", http.StatusBadRequest)

		return
	}

	assignStudent(metrics, client, facultyRequest)

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(fmt.Sprintf("Estudiante '%s' asignado a la disciplina '%d'", facultyRequest.Name, facultyRequest.Discipline)))
	if err != nil {
		log.Printf("Ocurrió un error al escribir la respuesta: %v", err)
	}
}

func healthCheckHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		writer.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(writer, "OK")
		if err != nil {
			log.Printf("Ocurrió un error al escribir la respuesta: %v", err)
		}

		return
	}

	writer.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	metrics := initializeMetrics()

	serverAddress := "discipline-service:8080"
	client, err := newClient(serverAddress)
	if err != nil {
		log.Fatalf("Ocurrió un error al crear el cliente: %v", err)
	}
	defer func(Connection *grpc.ClientConn) {
		err = Connection.Close()
		if err != nil {
			log.Printf("Ocurrió un error al cerrar la conexión: %v", err)
		}
	}(client.Connection)

	http.HandleFunc("/agronomy", func(writer http.ResponseWriter, request *http.Request) {
		requestHandler(writer, request, metrics, client)
	})
	http.HandleFunc("/agronomy/healthz", healthCheckHandler)

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Servidor de Agronomía escuchando en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
