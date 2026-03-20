package rabbitmq

import (
	"log"
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
	}

	return nil
}