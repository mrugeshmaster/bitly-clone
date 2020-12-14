package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Fatalf("%s: %s", msg, err)
// 	}
// }

// AWS
// var rabbitmq_server = "10.0.2.29"
// var rabbitmq_port = "5672"
// var rabbitmq_user = "bitly"
// var rabbitmq_pass = "bitly"

// // RabbitMQ Config
// var rabbitmq_server = "localhost"
// var rabbitmq_port = "5672"
// var rabbitmq_user = "guest"
// var rabbitmq_pass = "guest"
var shortlinkqueue = "shortlinkcreatequeue_lr"

// Receive Long URL and Short URL Code from Queue to store in database
func queue_receive() {
	conn, err := amqp.Dial("amqp://" + rabbitmq_user + ":" + rabbitmq_pass + "@" + rabbitmq_server + ":" + rabbitmq_port + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		shortlinkqueue, // name
		false,          // durable
		false,          // delete when usused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.ExchangeDeclare(
		"shortlinkexchange", // name
		"fanout",            // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.QueueBind(
		q.Name,                   // queue name
		"short-link-routing-key", // routing key
		"shortlinkexchange",      // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name,         // queue
		"linkredirect", // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
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
			writeInCache(queueData["shortlink_code"], queueData["longurl"])
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
