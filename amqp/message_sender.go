package amqp

import (
	"time"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/logging"
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
		logging.Printf("Adding a delay before sending message: %d\n", m.Notify.Amqp.Delay)
		time.Sleep(time.Duration(m.Notify.Amqp.Delay) * time.Second)
	}

	return sendMessage(m)
}

func sendMessage(m *definition.Mock) bool {
	conn, err := amqp.Dial(m.Notify.Amqp.URL)

	if err != nil {
		logging.Printf("Failed to connect to server: %s\n", err)
		return false
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		logging.Printf("Failed to open a channel: %s\n", err)
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
		logging.Printf("Failed to publish a message: %s\n", err)
		return false
	}

	logging.Println("Notified message by AMQP")
	return true
}

func failOnError(err error, msg string) {
	if err != nil {
		logging.Fatalf("%s: %s", msg, err)
	}
}
