package rabbitmq

import (
	"log"
	"encoding/json"
	"consumer/models"
	"consumer/db"
)

func StartConsumer(queueName string) error {
	ch := GetChannel()

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("Consumer rodando... aguardando mensagens")

	for msg := range msgs {
		log.Printf("Mensagem recebida: %s\n", msg.Body)

		var data models.TelemetryMessage
		err := json.Unmarshal(msg.Body, &data)
		if err != nil {
			log.Println("Erro ao converter JSON:", err)
			continue
		}

		err = db.InsertTelemetry(data)
		if err != nil {
			log.Println("Erro ao salvar no banco:", err)
			continue
		}

		log.Println("Dados salvos no banco")
	}

	return nil
}