package amqp

import (
	"log"
	"time"

	"github.com/jmartin82/mmock/definition"
	"github.com/streadway/amqp"
)

//MessageSender sends message
type MessageSender struct {
}

//Send message if needed
func (msender MessageSender) Send(m *definition.Mock) bool {
	if m.Notify.Amqp.URL == "" {
		return true
	}
	if m.Notify.Amqp.Delay > 0 {
		log.Printf("Adding a delay before sending message")
		time.Sleep(time.Duration(m.Notify.Amqp.Delay) * time.Second)
	}

	return sendMessage(m)
}

func sendMessage(m *definition.Mock) bool {
	conn, err := amqp.Dial(m.Notify.Amqp.URL)

	if err != nil {
		log.Println("Failed to connect to RabbitMQ")
		return false
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Println("Failed to open a channel")
		return false
	}

	defer ch.Close()

	err = ch.Publish(
		m.Notify.Amqp.Exchange,   // exchange
		m.Notify.Amqp.RoutingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Body:            []byte(m.Notify.Amqp.Body),
			ContentType:     m.Notify.Amqp.ContentType,
			ContentEncoding: m.Notify.Amqp.ContentEncoding,
			Priority:        m.Notify.Amqp.Priority,
			CorrelationId:   m.Notify.Amqp.CorrelationID,
			ReplyTo:         m.Notify.Amqp.ReplyTo,
			Expiration:      m.Notify.Amqp.Expiration,
			MessageId:       m.Notify.Amqp.Expiration,
			Timestamp:       m.Notify.Amqp.Timestamp,
			Type:            m.Notify.Amqp.Type,
			UserId:          m.Notify.Amqp.UserID,
			AppId:           m.Notify.Amqp.AppID,
			DeliveryMode:    2,
		})

	if err != nil {
		log.Println("Failed to publish a message")
		return false
	}

	log.Println("Notified message by AMQP")
	return true
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
