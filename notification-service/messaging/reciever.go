package messaging

import (
	"encoding/json"
	"log"

	"github.com/BernardN38/ebuy-server/notification-service/service"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageReceiver interface {
	ListenForMessages() error
}

type RabbitmqReceiver struct {
	notificationService service.Service
	conn                *amqp.Connection
	channel             *amqp.Channel
	queue               amqp.Queue
}

func New(connection *amqp.Connection, ns service.Service) (MessageReceiver, error) {
	r := &RabbitmqReceiver{
		conn: connection,
	}
	err := r.init()
	if err != nil {
		return nil, err
	}
	r.notificationService = ns
	return r, nil
}

func (r *RabbitmqReceiver) init() error {
	channel, err := r.conn.Channel()
	if err != nil {
		return err
	}

	// Declare exchange
	err = channel.ExchangeDeclare(
		"media_events", // exchange name
		"topic",        // exchange type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	// Declare queue
	queue, err := channel.QueueDeclare(
		"",    // queue name (empty to let RabbitMQ generate one)
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// Bind queue to exchange with routing key
	err = channel.QueueBind(
		queue.Name,         // queue name
		"media.compressed", // routing key
		"media_events",     // exchange name
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}
	err = channel.QueueBind(
		queue.Name,       // queue name
		"media.uploaded", // routing key
		"media_events",   // exchange name
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	r.channel = channel
	r.queue = queue

	return nil
}

func (r *RabbitmqReceiver) ListenForMessages() error {
	// Consume messages from the queue
	msgs, err := r.channel.Consume(
		r.queue.Name, // queue name
		"",           // consumer tag (empty to let RabbitMQ generate one)
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	// Process incoming messages
	go func() {
		for msg := range msgs {
			switch msg.RoutingKey {
			case "media.compressed":
				// var mediaUploadedPayload mediaUploadedPayload
				// json.Unmarshal(msg.Body, &mediaUploadedPayload)
				log.Printf("media compressed message recieved: %s", msg.Body)

			case "media.uploaded":
				var mediaUploadedPayload mediaUploadedPayload
				json.Unmarshal(msg.Body, &mediaUploadedPayload)
				log.Printf("media uploaded message recieved: %+v", mediaUploadedPayload)
				json, _ := json.Marshal(mediaUploadedPayload)
				r.notificationService.SendMessageToUser(mediaUploadedPayload.UserId, json)
			default:
				log.Printf("Received message: %s", msg.Body)
			}

		}
	}()

	// Block indefinitely to keep the consumer running
	select {}
}

type mediaUploadedPayload struct {
	MediaId              int32     `json:"mediaId"`
	ExternalIdCompressed uuid.UUID `json:"externalIdCompressed"`
	UserId               int32     `json:"userId"`
}