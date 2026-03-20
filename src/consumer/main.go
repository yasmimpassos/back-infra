package main

import (
	"log"
	"consumer/rabbitmq"
)

func main() {
	err := rabbitmq.InitRabbitMQ()
	if err != nil {
		log.Fatal("Erro ao conectar no RabbitMQ:", err)
	}
	defer rabbitmq.CloseRabbitMQ()

	err = rabbitmq.StartConsumer("telemetry_queue")
	if err != nil {
		log.Fatal("Erro ao iniciar consumer:", err)
	}
}