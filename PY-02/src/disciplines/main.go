package main

import (
	"context"
	"fmt"
	pb "github.com/janjo25/proto"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"time"
)

/*
La estructura 'athleticsServer' representa el servidor para el servicio de disciplinas en gRPC.

- Este tipo (struct) se define para manejar las solicitudes que lleguen a este servicio.
- Al heredar de 'pb.UnimplementedDisciplineServiceServer', se asegura que nuestro servidor tenga todas las
  implementaciones necesarias para funcionar correctamente.
- Sin este registro, aunque tengamos el servidor, no sabríamos cómo responder a las solicitudes de los clientes, ya que
  el servidor no tendría la lógica de negocio asociada. Esta plantilla nos ayuda a evitar errores al implementar los
  métodos requeridos.

En resumen, 'athleticsServer' es el modelo que se encargará de las solicitudes, y
'pb.UnimplementedDisciplineServiceServer' nos proporciona una base para asegurarnos de que todo esté en orden.
*/
type athleticsServer struct {
	pb.UnimplementedDisciplineServiceServer
}

func coinToss() bool {
	randomSource := rand.NewSource(time.Now().UnixNano()) // Se inicializa la semilla del generador con el tiempo.
	random := rand.New(randomSource)                      // Se crea un generador de números aleatorios con la semilla.

	return random.Intn(2) == 1
}

func (server *athleticsServer) Assign(ctx context.Context, request *pb.DisciplineRequest) (*pb.DisciplineResponse, error) {
	fmt.Printf("Estudiante con ID '%d' compitiendo...\n", request.StudentId)

	won := coinToss()

	log.Printf("¿Estudiante con ID '%d' ganó la competencia? %t\n", request.StudentId, won)

	return &pb.DisciplineResponse{
		Success: won,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("No se pudo establecer la conexión en el puerto 8080: %v", err)
	}

	/*
		Crear una instancia del servidor gRPC no es suficiente para permitir la comunicación con los clientes.
		Para que los clientes puedan interactuar con el servidor, es necesario registrar el servicio en el servidor gRPC.
		Esto vincula la lógica de negocio del servicio con el servidor, permitiendo que el servidor maneje las solicitudes.
	*/
	grpcServer := grpc.NewServer()
	pb.RegisterDisciplineServiceServer(grpcServer, &athleticsServer{})

	fmt.Println("Servidor de competencias iniciado en el puerto 8080...")

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Ocurrió un error al servir el servidor: %v", err)
	}
}
