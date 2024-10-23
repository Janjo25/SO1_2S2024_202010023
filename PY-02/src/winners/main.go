package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
)

func main() {
	brokers := []string{os.Getenv("KAFKA_BROKER")} // Recupera la dirección del broker de Kafka usando la variable de entorno 'KAFKA_BROKER'.
	topic := os.Getenv("KAFKA_TOPIC")              // Recupera el nombre del topic de Kafka usando la variable de entorno 'KAFKA_TOPIC'.

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Printf("Ocurrió un error al crear el consumidor de Kafka: %v", err)
	}

	defer func(consumer sarama.Consumer) {
		err = consumer.Close()
		if err != nil {
			log.Printf("Ocurrió un error al cerrar el consumidor de Kafka: %v", err)
		}
	}(consumer)

	/*
		Crea un consumidor de partición para leer mensajes de la partición 0 del topic especificado.
		Se comenzará a leer desde los mensajes más recientes, ya que se especifica 'sarama.OffsetNewest'.
	*/
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Ocurrió un error al consumir del topic '%s': %v", topic, err)
	}

	defer func(partitionConsumer sarama.PartitionConsumer) {
		err = partitionConsumer.Close()
		if err != nil {
			log.Printf("Ocurrió un error al cerrar el consumidor de partición: %v", err)
		}
	}(partitionConsumer)

	log.Println("Consumidor escuchando mensajes de los ganadores...")

	for message := range partitionConsumer.Messages() {
		fmt.Printf("Mensaje recibido (ganador): %s\n", string(message.Value))
	}
}
