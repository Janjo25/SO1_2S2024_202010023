package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func assignStudent(studentID int, discipline string, serverAddress string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ocurrió un error al crear el cliente: %v", err)
	}
	defer func(client *grpc.ClientConn) {
		err = client.Close()
		if err != nil {
			log.Fatalf("Ocurrió un error al cerrar la conexión: %v", err)
		}
	}(client)

	disciplineClient := pb.NewDisciplineServiceClient(client)

	request := &pb.DisciplineRequest{
		StudentId:  int32(studentID),
		Discipline: discipline,
	}

	response, err := disciplineClient.Assign(ctx, request)
	if err != nil {
		log.Fatalf("Ocurrió un error al llamar al servicio: %v", err)
	}

	if response.Success {
		fmt.Printf("El estudiante con ID %d ha sido asignado a la disciplina %s", studentID, discipline)
	} else {
		fmt.Printf("No se ha podido asignar al estudiante con ID %d a la disciplina %s", studentID, discipline)
	}
}

func main() {
	studentID := 1
	discipline := "athletics"

	serverAddresses := map[string]string{
		"athletics": "athletics-competition-service:8080",
		"boxing":    "boxing-competition-service:8080",
		"swimming":  "swimming-competition-service:8080",
	}

	serverAddress, exists := serverAddresses[discipline]

	if !exists {
		log.Fatalf("La disciplina %s no existe", discipline)
	}

	assignStudent(studentID, discipline, serverAddress)
}
