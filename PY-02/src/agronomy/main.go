package main

import (
	"context"
	pb "github.com/janjo25/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func assignStudent(studentID int, discipline string, serverAddress string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		log.Printf("Ocurrió un error al llamar al servicio: %v", err)
	}

	if response.Success {
		log.Printf("El estudiante con ID '%d' ha sido asignado a la disciplina '%s'", studentID, discipline)
	} else {
		log.Printf("No se ha podido asignar al estudiante con ID '%d' a la disciplina '%s'", studentID, discipline)
	}
}

func main() {
	studentID := 1
	discipline := "athletics"
	serverAddress := "disciplines-service:80"

	for {
		assignStudent(studentID, discipline, serverAddress)
		time.Sleep(5 * time.Second)
	}
}
