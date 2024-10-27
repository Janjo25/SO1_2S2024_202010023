package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	pb "github.com/janjo25/proto"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"time"
)

/*
La estructura 'disciplinesServer' representa el servidor para el servicio de disciplinas en gRPC.

- Este tipo (struct) se define para manejar las solicitudes que lleguen a este servicio.
- Al heredar de 'pb.UnimplementedDisciplineServiceServer', se asegura que nuestro servidor tenga todas las
  implementaciones necesarias para funcionar correctamente.
- Sin este registro, aunque tengamos el servidor, no sabríamos cómo responder a las solicitudes de los clientes, ya que
  el servidor no tendría la lógica de negocio asociada. Esta plantilla nos ayuda a evitar errores al implementar los
  métodos requeridos.

En resumen, 'disciplinesServer' es el modelo que se encargará de las solicitudes, y
'pb.UnimplementedDisciplineServiceServer' nos proporciona una base para asegurarnos de que todo esté en orden.
*/
type disciplinesServer struct {
	kafkaProducer sarama.SyncProducer
	pb.UnimplementedDisciplineServiceServer
}

func coinToss() bool {
	randomSource := rand.NewSource(time.Now().UnixNano()) // Se inicializa la semilla del generador con el tiempo.
	random := rand.New(randomSource)                      // Se crea un generador de números aleatorios con la semilla.

	return random.Intn(2) == 1
}

func createKafkaProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll // Espera a que todos los brokers confirmen.
	config.Producer.Retry.Max = 5                    // Número máximo de reintentos en caso de error.
	config.Producer.Return.Successes = true          // Habilita el retorno de mensajes exitosos.

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear el productor de Kafka: %v", err)
	}

	return producer, nil
}

func (server *disciplinesServer) Assign(_ context.Context, request *pb.DisciplineRequest) (*pb.DisciplineResponse, error) {
	fmt.Printf("Estudiante con nombre '%s' compitiendo en la disciplina '%d'...\n", request.Name, request.Discipline)

	won := coinToss()

	log.Printf("¿Estudiante con nombre '%s' ganó la competencia en la disciplina '%d'? %v\n", request.Name, request.Discipline, won)

	// Crea un mensaje con el resultado de la competencia para enviar a Kafka.
	message := fmt.Sprintf("Nombre: %s, Edad: %d, Facultad: %s, Disciplina: %d, Ganó: %v", request.Name, request.Age, request.Faculty, request.Discipline)

	topic := "olympics-losers"

	if won {
		topic = "olympics-winners"
	}

	// Publica el mensaje en el topic 'olympics-results' de Kafka.
	kafkaMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := server.kafkaProducer.SendMessage(kafkaMessage)
	if err != nil {
		log.Printf("Ocurrió un error al enviar mensaje a Kafka: %v", err)

		return nil, err
	} else {
		log.Printf("Mensaje enviado a Kafka: %s", message)
	}

	return &pb.DisciplineResponse{Success: true}, nil
}

func main() {
	// Crea un productor de Kafka. Un productor es lo que se utiliza para enviar mensajes a un topic de Kafka.
	brokers := []string{"kafka-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092"}
	producer, err := createKafkaProducer(brokers)
	if err != nil {
		log.Printf("Ocurrió un error al crear el productor de Kafka: %v", err)
	}

	defer func(producer sarama.SyncProducer) {
		err = producer.Close()
		if err != nil {
			log.Printf("Ocurrió un error al cerrar el productor de Kafka: %v", err)
		}
	}(producer)

	server := &disciplinesServer{
		kafkaProducer: producer,
	}

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
	pb.RegisterDisciplineServiceServer(grpcServer, server)

	fmt.Println("Servidor de competencias iniciado en el puerto 8080...")

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Ocurrió un error al servir el servidor: %v", err)
	}
}
