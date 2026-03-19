package main

import (
	"log"
	"backend/routes"
	"backend/rabbitmq"
)

func main() {
	err := rabbitmq.InitRabbitMQ()
	if err != nil {
		log.Fatal("Erro ao conectar no RabbitMQ:", err)
	}

	router := routes.SetupRouter()
	router.Run(":8080")
}