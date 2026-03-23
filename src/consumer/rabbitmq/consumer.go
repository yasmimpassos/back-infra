package rabbitmq

import (
	"log"
	"encoding/json"
	"consumer/models"
	"consumer/db"
)

func ProcessTelemetryMessage(body []byte, save func(models.TelemetryMessage) error) error {
	var data models.TelemetryMessage

	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	return save(data)
}

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

		err := ProcessTelemetryMessage(msg.Body, db.InsertTelemetry)
		if err != nil {
			log.Println("Erro ao processar mensagem:", err)
			continue
		}

		log.Println("Dados salvos no banco")
	}

	return nil
}