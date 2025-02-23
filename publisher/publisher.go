package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Koneksi ke RabbitMQ
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	failOnError(err, "Gagal terhubung ke RabbitMQ")
	defer conn.Close()

	// Membuat channel
	ch, err := conn.Channel()
	failOnError(err, "Gagal membuka channel")
	defer ch.Close()

	// Membuat queue
	q, err := ch.QueueDeclare(
		"test_queue", // Nama queue
		false,        // Durable
		false,        // Auto-delete
		false,        // Exclusive
		false,        // No-wait
		nil,          // Arguments
	)
	failOnError(err, "Gagal mendeklarasikan queue")

	// Mengirim pesan
	for i := 1; i <= 10; i++ {
		body := fmt.Sprintf("Pesan ke-%d", i)
		err = ch.Publish(
			"",     // Exchange
			q.Name, // Routing key (queue name)
			false,  // Mandatory
			false,  // Immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)
		failOnError(err, "Gagal mengirim pesan")
		fmt.Println(" [x] Dikirim:", body)
		time.Sleep(1 * time.Second) // Delay untuk simulasi pengiriman bertahap
	}

	fmt.Println("Semua pesan telah dikirim!")
}
