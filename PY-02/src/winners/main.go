package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis"
	"log"
	"os"
)

func main() {
	brokers := []string{os.Getenv("KAFKA_BROKER")} // Recupera la dirección del broker de Kafka usando la variable de entorno 'KAFKA_BROKER'.
	topic := os.Getenv("KAFKA_TOPIC")              // Recupera el nombre del topic de Kafka usando la variable de entorno 'KAFKA_TOPIC'.
	redisHost := os.Getenv("REDIS_HOST")           // Recupera la dirección del host de Redis usando la variable de entorno 'REDIS_HOST'.
	redisPort := os.Getenv("REDIS_PORT")           // Recupera el puerto de Redis usando la variable de entorno 'REDIS_PORT'.
	redisPassword := os.Getenv("REDIS_PASSWORD")   // Recupera la contraseña de Redis usando la variable de entorno 'REDIS_PASSWORD'.

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Printf("Ocurrió un error al conectar con Redis: %v", err)
	}

	log.Printf("Conexión exitosa con Redis en '%s:%s'", redisHost, redisPort)

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

		key := fmt.Sprintf("winner-%d", message.Offset)
		value := string(message.Value)

		err = redisClient.Set(key, value, 0).Err()
		if err != nil {
			log.Printf("Ocurrió un error al guardar el mensaje en Redis: %v", err)
		}

		log.Printf("Mensaje guardado en Redis con clave '%s'", key)
	}
}
