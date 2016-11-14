package amqp

import (
	"log"
	"time"

	"fmt"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/parse"
	"github.com/streadway/amqp"
)

//RabbitMQSender sends message to RabbitMQ
type RabbitMQSender struct {
	Parser parse.ResponseParser
}

//Send message to rabbitMQ if needed
func (rmqs RabbitMQSender) Send(per *definition.Persist, req *definition.Request, res *definition.Response) bool {
	if per.AMQP.URL == "" {
		return true
	}

	per.AMQP.Body = rmqs.Parser.ParseBody(req, res, per.AMQP.Body, per.AMQP.BodyAppend)

	if per.AMQP.Delay > 0 {
		log.Printf("Adding a delay before sending message")
		time.Sleep(time.Duration(per.AMQP.Delay) * time.Second)
	}

	return sendMessage(per.AMQP, res)
}

//NewRabbitMQSender creates a new RabbitMQSender
func NewRabbitMQSender(parser parse.ResponseParser) *RabbitMQSender {
	result := RabbitMQSender{Parser: parser}
	return &result
}

func sendMessage(publishInfo definition.AMQPPublishing, res *definition.Response) bool {
	conn, err := amqp.Dial(publishInfo.URL)

	if hasError(err, "Failed to connect to RabbitMQ", res) {
		return false
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if hasError(err, "Failed to open a channel", res) {
		return false
	}
	defer ch.Close()

	err = ch.Publish(
		publishInfo.Exchange,   // exchange
		publishInfo.RoutingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Body:            []byte(publishInfo.Body),
			ContentType:     publishInfo.ContentType,
			ContentEncoding: publishInfo.ContentEncoding,
			Priority:        publishInfo.Priority,
			CorrelationId:   publishInfo.CorrelationID,
			ReplyTo:         publishInfo.ReplyTo,
			Expiration:      publishInfo.Expiration,
			MessageId:       publishInfo.Expiration,
			Timestamp:       publishInfo.Timestamp,
			Type:            publishInfo.Type,
			UserId:          publishInfo.UserID,
			AppId:           publishInfo.AppID,
			DeliveryMode:    2,
		})

	if hasError(err, "Failed to publish a message", res) {
		return false
	}
	log.Printf(" [x] Sent %s", publishInfo.Body)
	return true
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func hasError(err error, msg string, res *definition.Response) bool {
	if err != nil {
		log.Print(err)
		res.Body = fmt.Errorf("%s: %s", msg, err).Error()
		res.StatusCode = 500
	}
	return err != nil
}
