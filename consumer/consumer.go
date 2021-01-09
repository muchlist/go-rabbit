package main

import (
	"encoding/json"
	"github.com/muchlist/go-rabbit/shared"
	"github.com/muchlist/go-rabbit/utils"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.HandleError(err, "tidak dapat terhubung ke AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	utils.HandleError(err, "tidak dapat membuat AMQP channel")
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("add", false, false, false, false, nil)
	utils.HandleError(err, "tidak dapat deklarasi `add` queue")

	err = amqpChannel.Qos(1, 0, false)
	utils.HandleError(err, "tidak dapat mengkonfigurasi qos")

	messageChannel, err := amqpChannel.Consume(
		queue.Name, "", false, false, false, false, nil)
	utils.HandleError(err, "tidak dapat meegistrasi konsumer")

	stopChan := make(chan bool)

	go func() {
		log.Printf("konsumer siap, PID : %d", os.Getpid())

		for d := range messageChannel {
			log.Printf("menerima message, %s", d.Body)

			addTask := &shared.AddTask{}

			err := json.Unmarshal(d.Body, addTask)
			if err != nil {
				log.Printf("error decoding JSON : %s", err)
			}

			log.Printf("hasil dari %d + %d adalah %d", addTask.Number1, addTask.Number2, addTask.Number1+addTask.Number2)
			if err := d.Ack(false); err != nil {
				log.Printf("error acknowledging message : %s", err)
			} else {
				log.Print("acknowledging message")
			}

		}
	}()

	// stop untuk terminasi program
	<-stopChan
}
