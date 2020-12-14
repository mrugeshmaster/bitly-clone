package main

import (
	"log"
	"os"
	"github.com/streadway/amqp"
)

// GCP
var rabbitmq_server = os.Getenv("rabbitmq_server")
var rabbitmq_port = os.Getenv("rabbitmq_port")
var rabbitmq_user = os.Getenv("rabbitmq_user")
var rabbitmq_pass = os.Getenv("rabbitmq_pass")

// AWS
// var rabbitmq_server = "10.0.2.217"
// var rabbitmq_port = "5672"
// var rabbitmq_user = "bitly"
// var rabbitmq_pass = "bitly"

// RabbitMQ Config
// var rabbitmq_server = "localhost"
// var rabbitmq_port = "5672"
// var rabbitmq_user = "guest"
// var rabbitmq_pass = "guest"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Send Order to Queue for Processing
func redirectLinkqueue_send(message string) {
	conn, err := amqp.Dial("amqp://" + rabbitmq_user + ":" + rabbitmq_pass + "@" + rabbitmq_server + ":" + rabbitmq_port + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"redirectlinkexchange", // name
		"fanout",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	body := message
	err = ch.Publish(
		"redirectlinkexchange",      // exchange
		"redirect-link-routing-key", // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
