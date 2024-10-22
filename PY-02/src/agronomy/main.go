package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/janjo25/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

type FacultyRequest struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Faculty    string `json:"faculty"`
	Discipline int    `json:"discipline"`
}

func assignStudent(request FacultyRequest, serverAddress string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Ocurrió un error al crear el cliente: %v", err)
	}
	defer func(client *grpc.ClientConn) {
		err = client.Close()
		if err != nil {
			log.Printf("Ocurrió un error al cerrar la conexión: %v", err)
		}
	}(client)

	disciplineClient := pb.NewDisciplineServiceClient(client)

	grpcRequest := &pb.DisciplineRequest{
		Name:       request.Name,
		Age:        int32(request.Age),
		Faculty:    request.Faculty,
		Discipline: int32(request.Discipline),
	}

	response, err := disciplineClient.Assign(ctx, grpcRequest)
	if err != nil {
		log.Printf("Ocurrió un error al llamar al servicio: %v", err)
	}

	if response.Success {
		log.Printf("El estudiante '%s' ha sido asignado a la disciplina '%d'", request.Name, request.Discipline)
	} else {
		log.Printf("No se ha podido asignar al estudiante '%s' a la disciplina '%d'", request.Name, request.Discipline)
	}
}

func requestHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Solo se permiten peticiones POST", http.StatusMethodNotAllowed)

		return
	}

	var facultyRequest FacultyRequest
	serverAddress := "disciplines-service:80"

	err := json.NewDecoder(request.Body).Decode(&facultyRequest)
	if err != nil {
		http.Error(writer, "No se pudo decodificar el cuerpo de la petición", http.StatusBadRequest)

		return
	}

	assignStudent(facultyRequest, serverAddress)

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(fmt.Sprintf("Estudiante '%s' asignado a la disciplina '%d'", facultyRequest.Name, facultyRequest.Discipline)))
	if err != nil {
		log.Printf("Ocurrió un error al escribir la respuesta: %v", err)

		return
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
	http.HandleFunc("/agronomy", requestHandler)
	http.HandleFunc("/agronomy/healthz", healthCheckHandler)

	log.Println("Servidor de Agronomía escuchando en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
