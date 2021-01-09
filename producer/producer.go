package main

import (
	"encoding/json"
	"github.com/muchlist/go-rabbit/shared"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err, "tidak dapat terhubung ke AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "tidak dapat membuat AMQP channel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("add", false, false, false, false, nil)
	handleError(err, "tidak dapat deklarasi `add` queue")

	rand.Seed(time.Now().UnixNano())

	addTask := shared.AddTask{Number1: rand.Intn(999), Number2: rand.Intn(999)}
	body, err := json.Marshal(addTask)
	if err != nil {
		handleError(err, "Error encoding JSON")
	}

	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	log.Printf("AddTask: %d+%d", addTask.Number1, addTask.Number2)

}
