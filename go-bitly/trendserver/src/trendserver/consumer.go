package main

import (
	"encoding/json"
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
var redirectlinkqueue = "redirectlinkqueue_ts"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Receive Long URL and Short URL Code from Queue to store in database
func queue_receive() {
	conn, err := amqp.Dial("amqp://" + rabbitmq_user + ":" + rabbitmq_pass + "@" + rabbitmq_server + ":" + rabbitmq_port + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		redirectlinkqueue, // name
		false,             // durable
		false,             // delete when usused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue")

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

	err = ch.QueueBind(
		q.Name,                      // queue name
		"redirect-link-routing-key", // routing key
		"redirectlinkexchange",      // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name,        // queue
		"trendserver", // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	failOnError(err, "Failed to register a consumer")

	queueData := map[string]string{}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
			if jsonErr := json.Unmarshal(d.Body, &queueData); jsonErr != nil {
				failOnError(jsonErr, "Failed to unmarshal json")
			}
			addToListOfLinks(queueData["shortlink_code"])
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
