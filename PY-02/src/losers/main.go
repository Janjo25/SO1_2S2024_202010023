package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis"
	"log"
	"os"
	"regexp"
	"strconv"
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
		log.Fatalf("Ocurrió un error al conectar con Redis: %s", err)
	}

	log.Printf("Conexión exitosa con Redis en '%s:%s'", redisHost, redisPort)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Ocurrió un error al crear el consumidor de Kafka: %s", err)
	}

	defer func(consumer sarama.Consumer) {
		err = consumer.Close()
		if err != nil {
			log.Printf("Ocurrió un error al cerrar el consumidor de Kafka: %s", err)
		}
	}(consumer)

	/*
		Crea un consumidor de partición para leer mensajes de la partición 0 del topic especificado.
		Se comenzará a leer desde los mensajes más recientes, ya que se especifica 'sarama.OffsetNewest'.
	*/
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Ocurrió un error al consumir del topic '%s': %s", topic, err)
	}

	defer func(partitionConsumer sarama.PartitionConsumer) {
		err = partitionConsumer.Close()
		if err != nil {
			log.Printf("Ocurrió un error al cerrar el consumidor de partición: %s", err)
		}
	}(partitionConsumer)

	log.Println("Consumidor escuchando mensajes de los perdedores...")

	for message := range partitionConsumer.Messages() {
		//fmt.Printf("Mensaje recibido (perdedor): %s\n", string(message.Value))

		regex := regexp.MustCompile(`Nombre: ([^,]+), Edad: (\d+), Facultad: ([^,]+), Disciplina: (\d+)`)
		matches := regex.FindStringSubmatch(string(message.Value))

		if len(matches) < 5 {
			log.Printf("Ocurrió un error al leer el mensaje: %s", string(message.Value))

			continue
		}

		name := matches[1]
		age, _ := strconv.Atoi(matches[2])
		faculty := matches[3]
		discipline, _ := strconv.Atoi(matches[4])

		key := fmt.Sprintf("loser-%d", message.Offset)
		value := map[string]interface{}{
			"nombre":     name,
			"edad":       age,
			"facultad":   faculty,
			"disciplina": discipline,
		}

		err = redisClient.HMSet(key, value).Err()
		if err != nil {
			log.Printf("Ocurrió un error al guardar el mensaje en Redis: %s", err)
		}

		// Contador de participantes por facultad.
		if faculty == "engineering" {
			err = redisClient.Incr("engineering-count").Err()
			if err != nil {
				log.Printf("Ocurrió un error al incrementar el contador de ingeniería: %s", err)
			}
		} else if faculty == "agronomy" {
			err = redisClient.Incr("agronomy-count").Err()
			if err != nil {
				log.Printf("Ocurrió un error al incrementar el contador de agronomía: %s", err)
			}
		}

		// Contador de perdedores por facultad.
		if faculty == "Ingeniería" {
			err = redisClient.Incr("engineering-loser-count").Err()
			if err != nil {
				log.Printf("Ocurrió un error al incrementar el contador de perdedores de ingeniería: %s", err)
			}
		} else if faculty == "Agronomía" {
			err = redisClient.Incr("agronomy-loser-count").Err()
			if err != nil {
				log.Printf("Ocurrió un error al incrementar el contador de perdedores de agronomía: %s", err)
			}
		}

		log.Printf("Mensaje guardado en Redis con clave '%s'", key)
	}
}
