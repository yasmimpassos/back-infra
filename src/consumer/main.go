package main

import (
	"log"
	"consumer/rabbitmq"
	"consumer/db"
)

func main() {
	err := rabbitmq.InitRabbitMQ()
	if err != nil {
		log.Fatal("Erro ao conectar no RabbitMQ:", err)
	}
	defer rabbitmq.CloseRabbitMQ()

	err = db.InitDB()
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}

	err = rabbitmq.StartConsumer("telemetry_queue")
	if err != nil {
		log.Fatal("Erro ao iniciar consumer:", err)
	}
}